package server

import (
	"errors"
	pb "github.com/1851616111/xchain/pkg/protos"
	"google.golang.org/grpc/peer"
	"log"
	"os"
)

var (
	handlerLog = log.New(os.Stderr, "handler:", log.LstdFlags)
)

func (n *Node) handshakeHandler(in *pb.HandShake, out chan *pb.Message, stream pb.Net_ConnectServer) {
	p, _ := peer.FromContext(stream.Context())

	validateFn := func(msg *pb.HandShake) error {
		if validateSource(msg, p.Addr.String()) {
			handlerLog.Printf("[handshake] validate client{%#v} ok\n", *msg.EndPoint)
			return nil
		}

		return InvalidateHandShakeAddressErr
	}

	//验证客户端连接的id， address， type 的正确性
	if err := validateFirstHandShake(in, validateFn); err != nil {
		handlerLog.Printf("[handshake] validate client error:%v: invalidated endpoint{%#v}\n", err, *in.EndPoint)
		out <- makeFirstHandShakeRspMsg(err)
		return
	}

	if n.Exist(p.Addr.String()) {
		handlerLog.Printf("[handshake] client{%#v} connection already exists\n", *in)
		out <- makeFirstHandShakeRspMsg(errors.New("connection already exists"))
	} else {
		out <- makeSecondHandShakeReqMsg(node.GetLocalEndPoint())
	}

	if err := n.Accept(*in.EndPoint, stream); err != nil {
		handlerLog.Printf("[handshake] accept client{%#v} connection err: %v\n", *in.EndPoint, err)
		out <- makeFirstHandShakeRspMsg(errors.New("accept connection falied"))
	}

	handlerLog.Printf("[handshake] handle client{%#v} success", *in)
	return
}

func (n *Node) pingHandler(in *pb.Message, out chan *pb.Message) {
	if isMsgRequest(in) {

		//发送ping的响应时，不包含自己的节点信息，但会返回请求节点的信息
		//在处理相应节点时需要去掉自己的节点信息
		epList := n.epManager.list()
		printEPList("ping", "处理请求，生成节点列表", epList)

		out <- makePingRspMsg(n.epManager.list())
		return

	} else if isMsgResponse(in) {

		pbList, err := parsePingRspMsg(in)
		if err != nil {
			out <- makeErrRspMsg(err)
			return
		}
		pbList = ListWithOutLocalEP(pbList, n.GetLocalEndPoint())

		printEPList("ping", "处理相应，对比节点列表", pbList)
		n.epManager.findNewEndPointHandler(pbList, func(ep *pb.EndPoint) {
			//todo 这里需不需要处理err？
			err := n.ConnectEntryPoint(ep.Address)
			handlerLog.Printf("[ping] handle endpoint %s err: %v\n", ep, err)
		})

	} else {
		handlerLog.Printf("[ping] unknow ping message %v\n", in)
		return
	}
}
