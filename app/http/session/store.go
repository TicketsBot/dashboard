package session

type SessionStore interface {
	Get(userId uint64) (SessionData, error)
	Set(userId uint64, data SessionData) error
	Clear(userId uint64) error
}

var Store SessionStore