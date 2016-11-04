package connection_manager

import "time"

var DefaultKeepaliveConfig KeepaliveConfig = KeepaliveConfig{
	MaxFailTime: 3,
}

type KeepaliveConfig struct {
	MaxFailTime uint8
	MaxTimeOut  time.Duration
}
