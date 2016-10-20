package connection_manager

import (
	"errors"
)

var (
	ErrConnectionNotExist    = errors.New("connection not exists")
	ErrConnectionsOutOfLimit = errors.New("connection out of limit")
)
