package server

import (
	"io"
	"log"
	"os"

	pb "github.com/1851616111/xchain/pkg/protos"
	"fmt"
)

var logger = log.New(os.Stdout, "[broker manager]", log.LstdFlags)

type ResponseWriter interface {
	Send(*pb.Instruction) error
}

func (m *manager) HandleEvent(e Event) {
	switch e.Kind {
	case EVENT_BROKER_START:
		if e.BrokerPort == "" || e.BrokerName == "" {
			return
		}

		m.recordBroker(e.BrokerName, e.BrokerPort)
		m.handleBroker(e.BrokerName)

	case EVENT_BROKER_STOP:
		m.stopBrokerHandler(e.BrokerName)
	}
}

func (m *manager) recordBroker(name, port string) {
	m.Lock()
	defer m.Unlock()

	if _, exist := m.stopChM[name]; exist {
		return
	}

	if _, exist := m.conM[name]; exist {
		return
	}

	m.stopChM[name] = make(chan struct{}, 1)
	m.nameToPortM[name] = port

	return
}

func (m *manager) handleBroker(broker string) {
	go func() {

		select {
		case <-m.stopChM[broker]:
			return
		default:
			con, err := connectBroker(m.nameToPortM[broker])
			if err != nil {
				//TODO:这么做很不优雅，因为连接broker的event会因为err消失
				logger.Printf("connect to broker{name:%v, address:%v} err:v", broker, m.nameToPortM[broker], err)
				m.stopBrokerHandler(broker)
				return
			}

			m.conM[broker] = con
			for {
				i, err := con.Recv()
				if err == io.EOF {
					logger.Println("broker handler read eof")
				}
				if err != nil {
					logger.Printf("broker handler err %v\n", err)
				}

				switch i.Type {
				case pb.Instruction_STATE:

					state, err := parseStateInstruction(i)
					if err != nil {
						logger.Printf("receive state message %v\n", state)
					}

					go m.handleState(broker, state, con)

					fmt.Printf("------->[debug]------>state: %#v\n", state)
				}
			}

		}
	}()

}

func (m *manager) stopBrokerHandler(broker string) {
	m.Lock()
	defer m.Unlock()

	c, exist := m.stopChM[broker]
	if !exist {
		return
	}

	c <- struct{}{}

	delete(m.conM, broker)
	delete(m.stopChM, broker)
	delete(m.nameToPortM, broker)
	return
}
