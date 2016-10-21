package server

import (
	"errors"
	pb "github.com/1851616111/xchain/pkg/protos"
	"google.golang.org/grpc/peer"
	//"log"
	"fmt"
	"io"
	"log"
)

var (
	UnSupportMsgErr             error = errors.New("unsupport msg type error")
	UnMatchHandShakeAddressErr  error = errors.New("unmatch handshake address")
	InvalidatedHandShakeIDErr   error = errors.New("invalide handshake id")
	InvalidatedHandShakeTypeErr error = errors.New("invalide handshake type")
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

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("read eof")
		}

		if err != nil {
			return err
		}

		switch in.Type {
		case pb.Message_Net_HANDSHAKE:

			p, _ := peer.FromContext(stream.Context())
			var ep *pb.EndPoint
			var retMsg *pb.Message

			if ep, err = validateFirstHandShake(in)(p.Addr.String()); err != nil {
				retMsg = makeFirstHandShakeRspMsg(err)
			} else if s.node.Exist(p.Addr.String()) {
				retMsg = makeFirstHandShakeRspMsg(errors.New("connection established"))
			} else {
				retMsg = makeSecondHandShakeReqMsg(node.GetLocalEndPoint())
			}

			if err = stream.Send(retMsg); err != nil {
				return err
			}

			return s.node.Accept(*ep, stream)

		case pb.Message_Net_PING:
			log.Printf("recv ping msg %s\b.", in.String())

		default:
			log.Printf("recv unsupport ping msg %s\b.", in.String())
			stream.Send(in)
		}
	}
}
