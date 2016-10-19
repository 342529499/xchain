package server

import (
	"github.com/golang/protobuf/ptypes"
	pb "github.com/1851616111/xchain/pkg/protos"
	"time"
)

func makeKeepaliveMsg() *pb.Message {
	timeStamp, _ := ptypes.TimestampProto(time.Now())

	return &pb.Message{
		Action:    pb.Action_Request,
		Type:      pb.Message_Net_PING,
		Payload:   []byte("ping"),
		Timestamp: timeStamp,
	}
}
