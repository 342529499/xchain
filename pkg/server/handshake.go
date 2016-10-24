package server

import (
	"errors"
	pb "github.com/1851616111/xchain/pkg/protos"
	cm "github.com/1851616111/xchain/pkg/server/connection_manager"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"strings"
	"sync"
	"time"
)

type handshakeManager struct {
	sync.RWMutex
	local   pb.EndPoint
	aside   pb.EndPoint
	success bool
}

func newHandshakeAgreement(localEndPoint pb.EndPoint, targetAddress string) *handshakeManager {
	return &handshakeManager{
		local:   localEndPoint,
		aside:   pb.EndPoint{Address: targetAddress},
		success: false,
	}
}

//first handshake:
//cli(192.168.1.1)---->server(192.168.1.2):
//server你好，我是ID=michael，Address=192.168.1.1的validate／nonvalidate。请先确认，请问，server你的名字是多少？
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

func makeFirstHandShakeRspMsg(err error) *pb.Message {
	handshakeInfo := &pb.HandShakeResponse{
		Type: pb.HandShakeResponse_Net_HANDSHAKE_FIRST_RESPONSE,
		Msg:  []byte(err.Error()),
	}

	payload, _ := proto.Marshal(handshakeInfo)
	timeStamp, _ := ptypes.TimestampProto(time.Now())

	return &pb.Message{
		Action:    pb.Action_Response,
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

func (h *handshakeManager) handlerJoin(con cm.Connection) (err error) {
	h.Lock()
	defer h.Unlock()

	msg := makeFirstHandShakeReqMsg(&h.local)
	if err = con.Send(msg); err != nil {
		return err
	}

	for {
		msg, err := con.Recv()
		if err != nil {
			return err
		}

		if !isHandShakeMsg(msg) {
			continue
		}

		successHandleFn := func(secondHandShake *pb.HandShake) error {
			if len(strings.TrimSpace(secondHandShake.EndPoint.Id)) == 0 {
				return errors.New("second handshake failded, endpoint id nil.")
			}

			if secondHandShake.EndPoint.Type != pb.EndPoint_VALIDATOR && secondHandShake.EndPoint.Type != pb.EndPoint_NON_VALIDATOR {
				return errors.New("second handshake failded, endpoint type nil")
			}

			h.aside.Id = secondHandShake.EndPoint.Id
			h.aside.Type = secondHandShake.EndPoint.Type

			return nil
		}

		err = handleSecHandShakeFunc(msg, successHandleFn)
		if err != nil {
			return err
		} else {
			break
		}
	}

	h.success = true
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

func validateFirstHandShake(msg *pb.HandShake, validateFn func(msg *pb.HandShake) error) (err error) {
	if err = validateFn(msg); err != nil {
		return
	}

	if msg.EndPoint.Type != pb.EndPoint_NON_VALIDATOR && msg.EndPoint.Type != pb.EndPoint_VALIDATOR {
		return InvalidatedHandShakeTypeErr
	}

	if msg.EndPoint.Id == "" {
		return InvalidatedHandShakeIDErr
	}

	return nil

}

//source maybe 192.168.1.230:53694 with client random port
func validateSource(msg *pb.HandShake, source string) bool {
	return strings.Contains(source, msg.EndPoint.Address)
}
