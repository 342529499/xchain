package connection_manager

type Manager interface {
	Exist(key string) bool
	Add(key string, con Connection) error
	Del(key string)
	Get(key string) (Connection, error)

	Keys() []string
	BroadcastFunc(ignoreError bool, cb func(string, Connection) error) error
}
