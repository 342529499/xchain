package server

import (
	cm "github.com/1851616111/xchain/pkg/server/connection_manager"
)

const (
	key_marker_server uint8 = 0
	key_marker_client uint8 = 1
)

type NetManager struct {
	serverConsManager cm.Manager
	clientConsManager cm.Manager

	keyMarker map[string]uint8
}

func newNetManager() *NetManager {
	m := new(NetManager)
	m.serverConsManager = cm.NewConnectionsManager(cm.Default_Suggested_MaxCons)
	m.clientConsManager = cm.NewConnectionsManager(cm.Default_Suggested_MaxCons)
	m.keyMarker = map[string]uint8{}

	return m
}

//Add(key string, con Connection) error
//Del(key string)
//Get(key string) (Connection, error)

func (m *NetManager) serverAdd(key string, con cm.Connection) error {

	m.keyMarker[key] = key_marker_server
	return m.serverConsManager.Add(key, con)
}

func (m *NetManager) clientAdd(key string, con cm.Connection) error {

	m.keyMarker[key] = key_marker_client
	return m.clientConsManager.Add(key, con)

}

func (m *NetManager) delete(key string) {

	marker, exist := m.keyMarker[key]
	if !exist {
		return
	}

	delete(m.keyMarker, key)
	switch marker {
	case key_marker_server:
		m.serverConsManager.Del(key)
	case key_marker_client:
		m.clientConsManager.Del(key)
	}
}

func (m *NetManager) get(key string) (cm.Connection, error) {

	marker, exist := m.keyMarker[key]
	if !exist {
		return nil, cm.ErrConnectionNotExist
	}

	if marker == key_marker_server {
		return m.serverConsManager.Get(key)
	} else {
		return m.clientConsManager.Get(key)
	}
}
