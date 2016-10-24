package server

import (
	"errors"
	pb "github.com/1851616111/xchain/pkg/protos"

	//"log"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
)

var (
	UnSupportMsgErr               error = errors.New("unsupport msg type error")
	InvalidateHandShakeAddressErr error = errors.New("invalide handshake address")
	InvalidatedHandShakeIDErr     error = errors.New("invalide handshake id")
	InvalidatedHandShakeTypeErr   error = errors.New("invalide handshake type")
)

func newNodeServer(id, address string, isValidator bool) *nodeServer {
	local := pb.EndPoint{
		Id:      id,
		Address: address,
	}

	if isValidator {
		local.Type = pb.EndPoint_VALIDATOR
	} else {
		local.Type = pb.EndPoint_NON_VALIDATOR
	}

	return &nodeServer{
		node: newNode(local),
	}
}

type nodeServer struct {
	node *Node
}

func (s *nodeServer) Connect(stream pb.Net_ConnectServer) error {

	rsp := &pb.Message{}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("read eof")
		}

		if err != nil {
			return err
		}

		switch msg.Type {
		case pb.Message_Net_HANDSHAKE:
			req := &pb.HandShake{}
			proto.Unmarshal(msg.Payload, req)

			s.node.handshakeHandler(req, rsp, stream)

		case pb.Message_Net_PING:

			s.node.pingHandler(msg, rsp)

			log.Printf("recv ping msg %s\b.", msg.String())

		default:
			log.Printf("recv unsupport ping msg %s\b.", msg.String())
		}

		if rsp != nil {
			stream.Send(rsp)
		}
	}
}
