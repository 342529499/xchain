package server

import (
	"errors"
	pb "github.com/1851616111/xchain/pkg/protos"
	"google.golang.org/grpc/peer"
	"log"
	"os"
)

var (
	handlerLogger = log.New(os.Stderr, "handler:", log.LstdFlags)
)

func (n *Node) handshakeHandler(in *pb.HandShake, out *pb.Message, stream pb.Net_ConnectServer) {
	p, _ := peer.FromContext(stream.Context())

	validateFn := func(msg *pb.HandShake) error {
		if validateSource(msg, p.Addr.String()) {
			handlerLogger.Printf("[handshake] validate client{%#v} ok\n", *msg.EndPoint)
			return nil
		}

		return InvalidateHandShakeAddressErr
	}

	//验证客户端连接的id， address， type 的正确性
	if err := validateFirstHandShake(in, validateFn); err != nil {
		handlerLogger.Printf("[handshake] validate client error:%v: invalidated endpoint{%#v}\n", err, *in.EndPoint)
		out = makeFirstHandShakeRspMsg(err)
		return
	}

	if n.Exist(p.Addr.String()) {
		handlerLogger.Printf("[handshake] client{%#v} connection already exists\n", *in)
		out = makeFirstHandShakeRspMsg(errors.New("connection already exists"))
	} else {
		out = makeSecondHandShakeReqMsg(node.GetLocalEndPoint())
	}

	if err := n.Accept(*in.EndPoint, stream); err != nil {
		handlerLogger.Printf("[handshake] accept client{%#v} connection err: %v\n", *in.EndPoint, err)
		out = makeFirstHandShakeRspMsg(errors.New("accept connection falied"))
	}

	handlerLogger.Printf("[handshake] handle client{%#v} success", *in)
	return
}
