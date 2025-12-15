package uhttp

import (
	"sync"
	"time"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/udb/redis"
)

// CSRFTokenStore CSRF Token 存储接口
type CSRFTokenStore interface {
	// Store 存储 Token
	Store(token string, expireSeconds int) error

	// Validate 验证 Token 是否存在且未过期
	Validate(token string) bool

	// Delete 删除 Token
	Delete(token string) error

	// Cleanup 清理过期 Token (可选,某些实现如Redis自动过期则无需实现)
	Cleanup()
}

// MemoryCSRFStore 内存存储实现
type MemoryCSRFStore struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
}

// NewMemoryCSRFStore 创建内存存储
func NewMemoryCSRFStore() *MemoryCSRFStore {
	return &MemoryCSRFStore{
		tokens: make(map[string]time.Time),
	}
}

// Store 存储 Token
func (s *MemoryCSRFStore) Store(token string, expireSeconds int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = time.Now().Add(time.Duration(expireSeconds) * time.Second)
	return nil
}

// Validate 验证 Token 是否存在且未过期
func (s *MemoryCSRFStore) Validate(token string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expireTime, exists := s.tokens[token]
	if !exists {
		return false
	}

	return time.Now().Before(expireTime)
}

// Delete 删除 Token
func (s *MemoryCSRFStore) Delete(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, token)
	return nil
}

// Cleanup 清理过期 Token
func (s *MemoryCSRFStore) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for token, expireTime := range s.tokens {
		if now.After(expireTime) {
			delete(s.tokens, token)
		}
	}
}

// RedisCSRFStore Redis 存储实现
type RedisCSRFStore struct {
	conn   *redis.Connection
	prefix string // Key 前缀,默认 "csrf:"
}

// NewRedisCSRFStore 创建 Redis 存储
// conn: Redis 连接
// prefix: Key 前缀,如果为空则使用默认值 "csrf:"
func NewRedisCSRFStore(conn *redis.Connection, prefix string) *RedisCSRFStore {
	if prefix == "" {
		prefix = "csrf:"
	}
	return &RedisCSRFStore{
		conn:   conn,
		prefix: prefix,
	}
}

// Store 存储 Token
func (s *RedisCSRFStore) Store(token string, expireSeconds int) error {
	key := s.prefix + token
	ctx := ucontext.New()
	return s.conn.Set(ctx, key, "1", time.Duration(expireSeconds)*time.Second)
}

// Validate 验证 Token 是否存在且未过期
func (s *RedisCSRFStore) Validate(token string) bool {
	key := s.prefix + token
	ctx := ucontext.New()
	val, err := s.conn.Get(ctx, key)
	return err == nil && val != ""
}

// Delete 删除 Token
func (s *RedisCSRFStore) Delete(token string) error {
	key := s.prefix + token
	ctx := ucontext.New()
	_, err := s.conn.Del(ctx, key)
	return err
}

// Cleanup Redis 自动过期,无需手动清理
func (s *RedisCSRFStore) Cleanup() {
	// Redis 会自动清理过期的 key,无需手动实现
}
