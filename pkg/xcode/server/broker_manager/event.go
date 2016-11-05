package server

var (
	EVENT_BROKER_START eventKind = 0
	EVENT_BROKER_STOP  eventKind = 1
)

type Notifier interface {
	Notify(Event)
}

type eventKind int

type Event struct {
	BrokerName string //名字为xcode的名字
	BrokerPort string //broker的监听地址
	Kind       eventKind
}
