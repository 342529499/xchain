package server

import (
	"errors"
	pb "github.com/1851616111/xchain/pkg/protos"
	"google.golang.org/grpc/peer"
	"log"
	"os"
)

var (
	logger = log.New(os.Stderr, "handler:", log.LstdFlags)
)

func (n *Node) handshakeHandler(in *pb.HandShake, out *pb.Message, stream pb.Net_ConnectServer) {
	p, _ := peer.FromContext(stream.Context())

	validateFn := func(msg *pb.HandShake) error {
		if validateSource(msg, p.Addr.String()) {
			return nil
		}

		logger.Printf("validate handshake err, source(%s) handshake{%v} unmatch\n", p.Addr.String(), *msg)
		return InvalidateHandShakeAddressErr
	}

	//验证客户端连接的id， address， type 的正确性
	if err := validateFirstHandShake(in, validateFn); err != nil {
		out = makeFirstHandShakeRspMsg(err)
		return
	}

	if n.Exist(p.Addr.String()) {
		logger.Printf("handshake fail, source(%s) handshake{%v} connect already exists\n", p.Addr.String(), *in)
		out = makeFirstHandShakeRspMsg(errors.New("connection already exists"))
	} else {

		out = makeSecondHandShakeReqMsg(node.GetLocalEndPoint())
	}

	if err := n.Accept(*in.EndPoint, stream); err != nil {
		logger.Printf("handshake fail, source(%s) handshake{%v} record connect err: %v\n", p.Addr.String(), *in, err)
		out = makeFirstHandShakeRspMsg(errors.New("accept connection falied"))
	}

	logger.Printf("handshake success, source(%s) handshake{%v}", p.Addr.String(), *in)
	return
}
