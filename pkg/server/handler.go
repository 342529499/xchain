package server

import (
	"errors"
	pb "github.com/1851616111/xchain/pkg/protos"
	"google.golang.org/grpc/peer"
	"log"
	"os"
)

var (
	handlerLog = log.New(os.Stderr, "[handler]", log.LstdFlags)
)

func (n *Node) handshakeHandler(in *pb.HandShake, out chan *pb.Message, stream pb.Net_ConnectServer) {
	p, _ := peer.FromContext(stream.Context())

	validateFn := func(msg *pb.HandShake) error {
		if validateSource(msg, p.Addr.String()) {
			return nil
		}

		return InvalidateHandShakeAddressErr
	}

	//验证客户端连接的id， address， type 的正确性
	if err := validateFirstHandShake(in, validateFn); err != nil {
		handlerLog.Printf("[handshake] validate client error:%v: invalidated endpoint{ID:%s}\n", err, in.EndPoint.Id)
		out <- makeFirstHandShakeRspMsg(err)
		return
	}

	if n.Exist(p.Addr.String()) {
		handlerLog.Printf("[handshake] client{ID:%s} connection already exists\n", in.EndPoint.Id)
		out <- makeFirstHandShakeRspMsg(errors.New("connection already exists"))
	} else {
		out <- makeSecondHandShakeReqMsg(node.GetLocalEndPoint())
	}

	if err := n.Accept(*in.EndPoint, stream); err != nil {
		handlerLog.Printf("[handshake] accept client{ID:%s} connection err: %v\n", in.EndPoint.Id, err)
		out <- makeFirstHandShakeRspMsg(errors.New("accept connection falied"))
	}

	handlerLog.Printf("[handshake] {ID:%s} connect success", in.EndPoint.Id)
	return
}

func (n *Node) pingHandler(in *pb.Message, out chan *pb.Message) {
	if in == nil {
		return
	}
	switch in.Action {
	case pb.Action_Request:
		out <- makePingRspMsg(n.epManager.list())
		return
	case pb.Action_Response:

		pbList, err := parsePingRspMsg(in)
		if err != nil {
			out <- makeErrRspMsg(err)
			return
		}
		pbList = ListWithOutLocalEP(pbList)

		n.epManager.findNewEndPointHandler(pbList, func(ep *pb.EndPoint) {
			//todo 这里需不需要处理err？
			err := n.ConnectEntryPoint(ep.Address + ":10690")
			handlerLog.Printf("[ping] handle endpoint %s err: %v\n", ep, err)
		})

	default:
		handlerLog.Printf("[pingHandler] unsupport message %v\n", in)
		return
	}
}
