package session

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) error
	Delete(key interface{}) error
	SessionID() string
}
