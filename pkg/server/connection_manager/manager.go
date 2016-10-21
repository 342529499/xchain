package connection_manager

import (
	errlist "github.com/1851616111/xchain/pkg/util/errors"
)

const Default_Suggested_MaxCons = 100

type conManager struct {
	conns    map[string]Connection
	maxConns int
	kaConfig KeepaliveConfig
}

func NewConnectionsManager(maxCon int) Manager {
	return &conManager{
		conns:    map[string]Connection{},
		maxConns: maxCon,
		kaConfig: DefaultKeepaliveConfig,
	}
}

func (p *conManager) Exist(key string) bool {

	_, exist := p.conns[key]
	return exist
}

func (p *conManager) Add(key string, con Connection) error {

	if len(p.conns) >= p.maxConns {
		return ErrConnectionsOutOfLimit
	}

	_, exist := p.conns[key]
	if exist {
		return ErrConnectionAlreadyExist
	}

	p.conns[key] = con
	return nil
}

func (p *conManager) Del(key string) {

	delete(p.conns, key)
}

func (p *conManager) Get(key string) (Connection, error) {

	con, exist := p.conns[key]
	if !exist {
		return nil, ErrConnectionNotExist
	}

	return con, nil
}

func (p *conManager) BroadcastFunc(waitAll bool, callback func(key string, con Connection) error) error {

	var l errlist.ErrorList = errlist.NewErrorList()
	for key, con := range p.conns {
		if err := callback(key, con); err != nil {
			if waitAll {
				l.Append(err)
			} else {
				return err
			}
		}
	}

	if l.Len() == 0 {
		return nil
	}

	return l
}
