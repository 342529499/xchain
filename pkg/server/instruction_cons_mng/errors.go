package instruction_cons_mng

import (
	"errors"
)

var (
	ErrConnectionNotExist     = errors.New("connection not exists")
	ErrConnectionAlreadyExist = errors.New("connection already exists")
	ErrConnectionsOutOfLimit  = errors.New("connection out of limit")
)
