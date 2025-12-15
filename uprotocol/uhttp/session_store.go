package uhttp

import (
	"context"
	"encoding/json"
	"github.com/whosafe/uf/uprotocol/unet"
	"sync"
	"time"
)

// MemorySession 内存会话
type MemorySession struct {
	id         string
	data       map[string]any
	lastAccess time.Time
	mu         sync.RWMutex
}

// ID 获取会话 ID
func (s *MemorySession) ID() string {
	return s.id
}

// Get 获取会话数据
func (s *MemorySession) Get(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

// Set 设置会话数据
func (s *MemorySession) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	s.lastAccess = time.Now()
}

// Delete 删除会话数据
func (s *MemorySession) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// Clear 清空会话数据
func (s *MemorySession) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]any)
}

// Save 保存会话 (内存存储不需要)
func (s *MemorySession) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastAccess = time.Now()
	return nil
}

// MemoryStore 内存会话存储
type MemoryStore struct {
	sessions    map[string]*MemorySession
	maxLifetime time.Duration
	mu          sync.RWMutex
}

// NewMemoryStore 创建内存会话存储
func NewMemoryStore() *MemoryStore {
	return NewMemoryStoreWithLifetime(30 * time.Minute)
}

// NewMemoryStoreWithLifetime 创建内存会话存储并指定存活周期
func NewMemoryStoreWithLifetime(maxLifetime time.Duration) *MemoryStore {
	store := &MemoryStore{
		sessions:    make(map[string]*MemorySession),
		maxLifetime: maxLifetime,
	}

	// 启动 GC (每 5 分钟清理一次)
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			store.GC(store.maxLifetime)
		}
	}()

	return store
}

// Get 获取会话
func (s *MemoryStore) Get(id string) (unet.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[id]
	if !ok {
		return nil, nil
	}

	session.lastAccess = time.Now()
	return session, nil
}

// Create 创建会话
func (s *MemoryStore) Create() (unet.Session, error) {
	id, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	session := &MemorySession{
		id:         id,
		data:       make(map[string]any),
		lastAccess: time.Now(),
	}

	s.mu.Lock()
	s.sessions[id] = session
	s.mu.Unlock()

	return session, nil
}

// Destroy 销毁会话
func (s *MemoryStore) Destroy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
	return nil
}

// GC 垃圾回收
func (s *MemoryStore) GC(maxLifetime time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for id, session := range s.sessions {
		session.mu.RLock()
		if now.Sub(session.lastAccess) > maxLifetime {
			delete(s.sessions, id)
		}
		session.mu.RUnlock()
	}
}

// RedisSession Redis 会话
type RedisSession struct {
	id      string
	data    map[string]any
	store   *RedisStore
	changed bool
	mu      sync.RWMutex
}

// ID 获取会话 ID
func (s *RedisSession) ID() string {
	return s.id
}

// Get 获取会话数据
func (s *RedisSession) Get(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

// Set 设置会话数据
func (s *RedisSession) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	s.changed = true
}

// Delete 删除会话数据
func (s *RedisSession) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	s.changed = true
}

// Clear 清空会话数据
func (s *RedisSession) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]any)
	s.changed = true
}

// Save 保存会话到 Redis
func (s *RedisSession) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.changed {
		return nil
	}

	// 序列化数据
	data, err := json.Marshal(s.data)
	if err != nil {
		return err
	}

	// 保存到 Redis
	key := s.store.prefix + s.id
	ctx := context.Background()
	if err := s.store.client.Set(ctx, key, data, s.store.maxLifetime).Err(); err != nil {
		return err
	}

	s.changed = false
	return nil
}

// RedisClient Redis 客户端接口 (兼容 go-redis)
type RedisClient interface {
	Get(ctx context.Context, key string) RedisStringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) RedisStatusCmd
	Del(ctx context.Context, keys ...string) RedisIntCmd
}

// RedisStringCmd String 命令结果接口
type RedisStringCmd interface {
	Bytes() ([]byte, error)
}

// RedisStatusCmd Status 命令结果接口
type RedisStatusCmd interface {
	Err() error
}

// RedisIntCmd Int 命令结果接口
type RedisIntCmd interface {
	Err() error
}

// RedisStore Redis 会话存储
type RedisStore struct {
	client      RedisClient
	prefix      string
	maxLifetime time.Duration
}

// NewRedisStore 创建 Redis 会话存储
func NewRedisStore(client RedisClient, prefix string, maxLifetime time.Duration) *RedisStore {
	return &RedisStore{
		client:      client,
		prefix:      prefix,
		maxLifetime: maxLifetime,
	}
}

// Get 获取会话
func (s *RedisStore) Get(id string) (unet.Session, error) {
	key := s.prefix + id
	ctx := context.Background()

	// 从 Redis 获取数据
	result := s.client.Get(ctx, key)
	data, err := result.Bytes()
	if err != nil {
		// 会话不存在
		return nil, nil
	}

	// 反序列化数据
	sessionData := make(map[string]any)
	if err := json.Unmarshal(data, &sessionData); err != nil {
		return nil, err
	}

	session := &RedisSession{
		id:      id,
		data:    sessionData,
		store:   s,
		changed: false,
	}

	return session, nil
}

// Create 创建会话
func (s *RedisStore) Create() (unet.Session, error) {
	id, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	session := &RedisSession{
		id:      id,
		data:    make(map[string]any),
		store:   s,
		changed: true,
	}

	// 立即保存到 Redis
	if err := session.Save(); err != nil {
		return nil, err
	}

	return session, nil
}

// Destroy 销毁会话
func (s *RedisStore) Destroy(id string) error {
	key := s.prefix + id
	ctx := context.Background()
	return s.client.Del(ctx, key).Err()
}

// GC 垃圾回收 (Redis 自动过期,不需要手动 GC)
func (s *RedisStore) GC(maxLifetime time.Duration) {
	// Redis 会自动过期,无需手动 GC
}

// RedisStoreConfig Redis 存储配置
type RedisStoreConfig struct {
	Client      RedisClient
	Prefix      string
	MaxLifetime time.Duration
}

// NewRedisStoreWithConfig 使用配置创建 Redis 存储
func NewRedisStoreWithConfig(cfg *RedisStoreConfig) *RedisStore {
	if cfg.Prefix == "" {
		cfg.Prefix = "session:"
	}
	if cfg.MaxLifetime == 0 {
		cfg.MaxLifetime = 30 * time.Minute
	}

	return NewRedisStore(cfg.Client, cfg.Prefix, cfg.MaxLifetime)
}

// 示例: 如何使用 go-redis 客户端
//
// import "github.com/redis/go-redis/v9"
//
// redisClient := redis.NewClient(&redis.Options{
//     Addr: "localhost:6379",
// })
//
// store := uhttp.NewRedisStore(redisClient, "session:", 30*time.Minute)
// sessionMgr := uhttp.NewSessionManager(store, "session_id", 3600)
