package server

import (
	"errors"
	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"strings"
	"sync"
	"time"
)

type handshakeManager struct {
	sync.RWMutex
	from    pb.EndPoint
	to      pb.EndPoint
	success bool
}

func newHandshakeAgreement(localEndPoint pb.EndPoint, targetAddress string) *handshakeManager {
	return &handshakeManager{
		from:    localEndPoint,
		to:      pb.EndPoint{Address: targetAddress},
		success: false,
	}
}

//first handshake:
//cli(192.168.1.1)---->server(192.168.1.2):
//server你好，我是ID=michael，Address=192.168.1.1的客户端。请先确认，请问，server你的名字是多少？
func makeFirstHandShakeReqMsg(localEPInfo *pb.EndPoint) *pb.Message {
	handshakeInfo := &pb.HandShake{
		Type:     pb.HandShake_Net_HANDSHAKE_FIRST,
		EndPoint: localEPInfo,
	}

	payload, _ := proto.Marshal(handshakeInfo)
	timeStamp, _ := ptypes.TimestampProto(time.Now())

	return &pb.Message{
		Action:    pb.Action_Request,
		Type:      pb.Message_Net_HANDSHAKE,
		Payload:   payload,
		Timestamp: timeStamp,
	}
}

//second handshake:
//server(192.168.1.2) ----> cli(192.168.1.1)
//michael你好,经过验证你的ip正确，准许接入。我的名字为Jessie。
func makeSecondHandShakeReqMsg(localEPInfo *pb.EndPoint) *pb.Message {
	handshakeInfo := &pb.HandShake{
		Type:     pb.HandShake_Net_HANDSHAKE_SECOND,
		EndPoint: localEPInfo,
	}

	payload, _ := proto.Marshal(handshakeInfo)
	timeStamp, _ := ptypes.TimestampProto(time.Now())

	return &pb.Message{
		Action:    pb.Action_Request,
		Type:      pb.Message_Net_HANDSHAKE,
		Payload:   payload,
		Timestamp: timeStamp,
	}
}

////third handshake:我已经确认你的信息（ID）。可以建立连接。
//func makeThirdHandShakeMsg() *pb.Message {
//	handshakeInfo := &pb.HandShake{
//		Type: pb.HandShake_Net_HANDSHAKE_THIRD,
//	}
//
//	payload, _ := proto.Marshal(handshakeInfo)
//	timeStamp, _ := ptypes.TimestampProto(time.Now())
//
//	return &pb.Message{
//		Action: pb.Action_Request,
//		Type:   pb.Message_Net_HANDSHAKE,
//		Payload: payload,
//		Timestamp: timeStamp,
//	}
//}

func (h *handshakeManager) handlerJoin(con Connection) (err error) {
	h.Lock()
	defer h.Unlock()

	msg := makeFirstHandShakeReqMsg(&h.from)
	if err = con.Send(msg); err != nil {
		return err
	}

	for {
		msg, err = con.Recv()
		if err != nil {
			return err
		}

		if !isHandShakeMsg(msg) {
			continue
		}

		err = handleSecHandShakeFunc(msg, validateSecondHandShakeMsg)
		if err != nil {
			return err
		} else {
			break
		}
	}

	h.success = true
	return nil
}

func validateSecondHandShakeMsg(msg *pb.HandShake) error {
	if len(strings.TrimSpace(msg.EndPoint.Id)) == 0 {
		return errors.New("second handshake failded")
	}

	return nil
}

func isHandShakeMsg(msg *pb.Message) bool {
	if msg.Type != pb.Message_Net_HANDSHAKE {
		return false
	}

	return true
}

func handleSecHandShakeFunc(msg *pb.Message, callback func(hs *pb.HandShake) error) error {
	handshake := &pb.HandShake{}
	if err := proto.Unmarshal(msg.Payload, handshake); err != nil {
		return err
	}

	if handshake.Type == pb.HandShake_Net_HANDSHAKE_SECOND {
		if err := callback(handshake); err != nil {
			return err
		}

		return nil
	}

	return errors.New("unknown net hankshake type")
}
