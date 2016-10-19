package server

import (
	"errors"
	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/peer"
	"io"
	"log"
)

var (
	UnSupportMsgErr            error = errors.New("unsupport msg type error")
	UnMatchHandShakeAddressErr error = errors.New("unmatch hankshake address")
)

type netServer struct {}

func (s *netServer) Connect(stream pb.Net_ConnectServer) error {
	return handle(stream)
}

func handle(stream pb.Net_ConnectServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err

		}

		switch in.Type {
		case pb.Message_Net_HANDSHAKE:
			var fromClientAddress string

			p, _ := peer.FromContext(stream.Context())
			handleSecHandShakeFunc(in, func(hs *pb.HandShake) error {
				if p.Addr.String() == hs.EndPoint.Address {
					fromClientAddress = hs.EndPoint.Address
					return nil
				}
				return UnMatchHandShakeAddressErr
			})

			if err = stream.Send(makeSecondHandShakeReqMsg(GetLocalEndPoint())); err != nil {
				return err
			}

			ep := &pb.EndPoint{}
			if err = proto.Unmarshal(in.Payload, ep); err != nil {
				log.Printf("handle net hand shake parse paylaod err %v\n", err)
				continue
			}

			key := pair{
				littler: *GetLocalEndPoint(),
				bigger: *ep,
			}

			ServerAddConnection(key, stream)

			PrintConnectionManager()

		case pb.Message_Net_PING:
			log.Printf("recv ping msg %s\b.", in.String())

		default:
			log.Printf("recv unsupport msg %s\b.", in.String())
			return UnSupportMsgErr
		}

	}
}
