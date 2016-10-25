package server

import (
	"log"

	pb "github.com/1851616111/xchain/pkg/protos"
	cm "github.com/1851616111/xchain/pkg/server/connection_manager"
	"google.golang.org/grpc"

	"code.google.com/p/go.net/context"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	clientLogger            = log.New(os.Stderr, "[Event]", log.LstdFlags)
	No_EntryPoint_Err error = errors.New("no entrypoint")
)

func (n *Node) startAndJoin(address string, successFn func(target pb.EndPoint, con cm.Connection) error) error {

	//开发阶段，此处临时使用insure option。稍后需要将配置写到node上
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		clientLogger.Printf("join net(%s) dial err: %v.", address, err)
		return err
	}

	cli := pb.NewNetClient(conn)

	serverStream, err := cli.Connect(context.Background())
	if err != nil {
		clientLogger.Printf("join net(%s) connect err: %v.", address, err)
		return err
	}

	agreement := newHandshakeAgreement(*n.GetLocalEndPoint(), address)
	if err := agreement.handlerJoin(serverStream); err != nil {
		return err
	}

	return successFn(agreement.aside, serverStream)
}

//node 启动时连接的的网络entrypoint address
//当node 的server aside 连接断开后，可能用到此连接方法。
//收到 ping的Message后，与本地对比，如果发现为新的endpoint,进行连接
func (n *Node) ConnectEntryPoint(entryPoint string) error {

	//无加入网络的目标，只能等待连接
	if len(entryPoint) == 0 {
		return No_EntryPoint_Err
	}

	//设计上应该为1，防止阻塞
	errCh, doneCh := make(chan error, 10), make(chan struct{}, 10)
	n.lounchClientCh <- &lounchConnectionMetadata{
		targetAddress: entryPoint,
		errCh:         errCh,
		doneCh:        doneCh,
	}

	select {
	case err := <-errCh:
		return err
		clientLogger.Printf("connect entryPoint:%s err:%v\n", entryPoint, err.Error())
		return err

	case <-doneCh:
		clientLogger.Printf("connect entryPoint:%s success\n", entryPoint)
		return nil
	}
}

func clientConnectionHandler(ep pb.EndPoint, con cm.Connection, done chan struct{}) error {
	node := getNode()

	rsp := make(chan *pb.Message, 30)
	go func() {
		for {
			select {
			case msg := <-rsp:
				if Is_Develop_Mod {
					clientLogger.Printf("[client id:%s] sending message %s\n", ep.Id, formatMessage(msg))
				}
				con.Send(msg)
			}
		}
	}()

	for {

		select {
		case <-done:
			return nil
		default:
		}

		msg, err := con.Recv()
		if err == io.EOF {
			clientLogger.Println("read eof")
		}

		if err != nil {
			return err
		}

		switch msg.Type {
		case pb.Message_Net_PING:
			node.pingHandler(msg, rsp)
		default:
			log.Printf("[client:%s] recv unsupport ping msg %s\b.", msg.String())
		}
	}
	return nil
}

func formatMessage(msg *pb.Message) string {
	return fmt.Sprintf("message{Action: %s, Type: %s, PayLoad:%s}", msg.Action, msg.Type.String(), string(msg.Payload))
}
