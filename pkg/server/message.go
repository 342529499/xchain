package server

import (
	"fmt"
	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
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

func makePingReqMsg() *pb.Message {
	ping := &pb.Ping{}
	timeStamp, _ := ptypes.TimestampProto(time.Now())
	payLoad, _ := proto.Marshal(ping)
	return &pb.Message{
		Action:    pb.Action_Request,
		Type:      pb.Message_Net_PING,
		Payload:   payLoad,
		Timestamp: timeStamp,
	}
}

func makePingRspMsg(epList []*pb.EndPoint) *pb.Message {
	ping := &pb.Ping{
		EndPoint: epList,
	}

	timeStamp, _ := ptypes.TimestampProto(time.Now())
	payLoad, err := proto.Marshal(ping)
	if err != nil {
		fmt.Printf("make ping response message err %v\n", err)
	}
	return &pb.Message{
		Action:    pb.Action_Response,
		Type:      pb.Message_Net_PING,
		Payload:   payLoad,
		Timestamp: timeStamp,
	}
}

func parsePingRspMsg(in *pb.Message) ([]*pb.EndPoint, error) {
	ping := &pb.Ping{}
	if err := proto.Unmarshal(in.Payload, ping); err != nil {
		return ping.EndPoint, err
	}

	return ping.EndPoint, nil
}

func isMsgRequest(msg *pb.Message) bool {
	return msg.Action == pb.Action_Request
}

func isMsgResponse(msg *pb.Message) bool {
	return msg.Action == pb.Action_Response
}

func makeErrRspMsg(err error) *pb.Message {
	timeStamp, _ := ptypes.TimestampProto(time.Now())

	return &pb.Message{
		Action:    pb.Action_Response,
		Type:      pb.Message_Error,
		Payload:   []byte(err.Error()),
		Timestamp: timeStamp,
	}
}
