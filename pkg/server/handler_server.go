package server

import (
	"errors"
	"fmt"
	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
)

var (
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

	newNode(local)

	return &nodeServer{}
}

type nodeServer struct {
}

func (s *nodeServer) Connect(stream pb.Net_ConnectServer) error {

	return serverConnectionHandler(stream)
}

func serverConnectionHandler(stream pb.Net_ConnectServer) error {
	node := getNode()

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

			node.handshakeHandler(req, rsp, stream)

		case pb.Message_Net_PING:

			node.pingHandler(msg, rsp)

			log.Printf("recv ping msg %s\b.", msg.String())

		default:
			log.Printf("recv unsupport ping msg %s\b.", msg.String())
		}

		if rsp != nil {
			if Is_Develop_Mod {
				fmt.Printf("sending message %v\n", *rsp)
			}

			stream.Send(rsp)
		}

		rsp = new(pb.Message)
	}
	return nil
}
