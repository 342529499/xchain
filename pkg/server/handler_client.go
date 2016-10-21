package server

import (
	"log"

	pb "github.com/1851616111/xchain/pkg/protos"
	cm "github.com/1851616111/xchain/pkg/server/connection_manager"
	"google.golang.org/grpc"

	"code.google.com/p/go.net/context"
	"os"
)

var (
	logger = log.New(os.Stderr, "handler_client", log.LstdFlags)
)

func (n *Node) startAndJoin(address string, successFn func(target pb.EndPoint, con cm.Connection) error) error {

	//开发阶段，此处临时使用insure option。稍后需要将配置写到node上
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logger.Printf("join net(%s) dial err: %v.", address, err)
		return err
	}

	cli := pb.NewNetClient(conn)

	serverStream, err := cli.Connect(context.Background())
	if err != nil {
		logger.Printf("join net(%s) connect err: %v.", address, err)
		return err
	}

	agreement := newHandshakeAgreement(*n.GetLocalEndPoint(), address)
	if err := agreement.handlerJoin(serverStream); err != nil {
		return err
	}

	return successFn(agreement.aside, serverStream)
}
