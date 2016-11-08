package server

import (
	im "github.com/1851616111/xchain/pkg/server/instruction_cons_mng"
	"sync"
)

var (
	managerOnce   *sync.Once = new(sync.Once)
	brokerManager *manager
	nodeID        string
	nodeAddress   string
)

func GetBrokerManager(nodeID, nodeAddress string) *manager {
	nodeID, nodeAddress = nodeID, nodeAddress
	if brokerManager == nil {
		managerOnce.Do(func() {
			brokerManager = &manager{
				notifier: make(chan Event, 50),
				stopChM:  map[string]chan struct{}{},
				conM:     map[string]im.Connection{},
				nameToPortM: map[string]string{},
			}
		})
	}
	return brokerManager
}

type manager struct {
	sync.RWMutex
	started     bool
	notifier    chan Event
	stopChM     map[string]chan struct{}
	conM        map[string]im.Connection
	nameToPortM map[string]string
}

func (m *manager) Notify(e Event) {
	if !m.started {
		m.start()
	}
	m.notifier <- e
}

func (m *manager) start() {
	go func() {
		for {
			select {
			case e := <-m.notifier:
				logger.Printf("new event %v\n", e)
				m.HandleEvent(e)
			}
		}
	}()
}
