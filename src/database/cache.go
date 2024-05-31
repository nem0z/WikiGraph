package database

type Cache interface {
	Get(key string) (string, error)
	GetInt64(key string) (int64, error)
	Set(key string, value interface{}) error
}
