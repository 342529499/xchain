package connection_manager

type Manager interface {
	Exist(key string) bool
	Add(key string, con Connection) error
	Del(key string)
	Get(key string) (Connection, error)

	BroadcastFunc(waitAll bool, cb func(string, Connection) error) error
}
