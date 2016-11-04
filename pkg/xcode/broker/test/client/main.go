package main

import (
	"code.google.com/p/go.net/context"
	"fmt"
	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/1851616111/xchain/pkg/xcode/broker"
	"google.golang.org/grpc"
	"log"
)

func main() {

	conn, err := grpc.Dial("0.0.0.0:10692", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	cli := pb.NewCodeServiceClient(conn)

	c, err := cli.Execute(context.Background())
	if err != nil {
		log.Panic(err)
	}

	go func() {
		for {

			i, err := c.Recv()
			if err != nil {
				log.Panic("receive instruction err:", err)
				return
			}

			fmt.Printf("receive instruction %v\n", i)

		}
	}()

	if err := c.Send(broker.NewStartInstruction("michael", "127.0.0.1", pb.XCodeSpec_GOLANG.String())); err != nil {
		log.Printf("send start instruction err %v\n", err)
	}

	if err := c.Send(broker.NewInitInstruction("init001", []string{"1", "2", "3"})); err != nil {
		log.Printf("send start instruction err %v\n", err)
	}

	select {}
}
