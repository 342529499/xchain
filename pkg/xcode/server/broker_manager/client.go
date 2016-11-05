package server

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/1851616111/xchain/pkg/xcode/broker"
	"io"
	"log"
)

func connectBroker(port string) (pb.CodeService_ExecuteClient, error) {
	conn, err := grpc.Dial("0.0.0.0"+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	cli := pb.NewCodeServiceClient(conn)

	c, err := cli.Execute(context.Background())
	if err != nil {
		return nil, err
	}

	startInst := broker.NewStartInstruction(nodeID, nodeAddress, pb.XCodeSpec_GOLANG.String())
	if err := c.Send(startInst); err != nil {
		return nil, err
	}

	for {
		rsp, err := c.Recv()
		if err == io.EOF {
			return nil, err
		}
		if err != nil {
			log.Println("connect broker err:%v\n", err)
		}

		if rsp.Identifier == startInst.Identifier && broker.IsOKInstruction(rsp) {
			return c, nil
		}
	}

}
