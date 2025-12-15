package unet

import "time"

// Session 会话接口
type Session interface {
	ID() string
	Get(key string) (any, bool)
	Set(key string, value any)
	Delete(key string)
	Clear()
	Save() error
}

// SessionStore 会话存储接口
type SessionStore interface {
	Get(id string) (Session, error)
	Create() (Session, error)
	Destroy(id string) error
	GC(maxLifetime time.Duration)
}
