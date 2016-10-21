package server

import (
	"errors"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/1851616111/xchain/pkg/protos"

	"fmt"
	"github.com/1851616111/xchain/pkg/util"
	"os"
)

func NewAndStartGrpcServer(option *ServerOptions) error {
	if option == nil {
		option = getDefaultServerOptions()
	} else if err := option.validate(); err != nil {
		return err
	}

	option.println()

	lis, err := net.Listen("tcp", option.ListenerAddress)
	if err != nil {
		log.Printf("new grpc server listen address err %v\n", err)
		return err
	}

	var opts []grpc.ServerOption
	if option.TlsEnabled {
		opt, err := newTlsOption(option.CertFilePath, option.KeyFilePath)
		if err != nil {
			log.Printf("new grpc server listen address options err %v\n", err)
			return err
		}
		opts = append(opts, opt)
	}

	server := grpc.NewServer(opts...)

	nodeServer := newNodeServer(option.ID, option.Address, option.IsValidator)
	pb.RegisterNetServer(server, nodeServer)

	go nodeServer.node.RunController()

	go func() {
		//无加入网络的目标，只能等待连接
		if len(option.EntryPointAddress) == 0 {
			return
		}

		//设计上应该为1，防止阻塞
		errCh, doneCh := make(chan error, 1), make(chan struct{}, 1)
		nodeServer.node.lounchConnectCh <- &lounchConnectionMetadata{
			targetAddress: option.EntryPointAddress,
			errCh: errCh,
			doneCh:doneCh,
		}

		select {
		case err := <-errCh:
			fmt.Printf("join entryPointAddress:%s err:%v\n", option.EntryPointAddress, err.Error())
			os.Exit(1)

		case <- doneCh:
			fmt.Printf("join entryPointAddress:%s success\n", option.EntryPointAddress)
		}

	}()

	if err := server.Serve(lis); err != nil {
		return err
	}

	return nil
}

type ServerOptions struct {
	ID string

	//192.168.1.1
	Address string

	//0.0.0.0:10690
	ListenerAddress string
	IsValidator     bool

	//网络的一个节点地址
	EntryPointAddress string

	TlsEnabled   bool
	CertFilePath string
	KeyFilePath  string
}

func (opts *ServerOptions) validate() error {
	if len(strings.TrimSpace(opts.ID)) == 0 {
		return errors.New("param id must not be nil.")
	}

	if len(strings.TrimSpace(opts.Address)) == 0 {
		return errors.New("param netaddress must not be nil.")
	}

	if len(strings.TrimSpace(opts.ListenerAddress)) == 0 {
		return errors.New("param listener address must not be nil.")
	}

	if opts.TlsEnabled {
		if len(strings.TrimSpace(opts.CertFilePath)) == 0 {
			return errors.New("param certFilePath not found.")
		}
		if len(strings.TrimSpace(opts.KeyFilePath)) == 0 {
			return errors.New("param keyFilePath not found.")
		}
	}

	return nil
}

func (o *ServerOptions) println() {
	fmt.Printf("------------- xchain start ---------------\n")
	fmt.Printf("ID: %s\n", o.ID)
	fmt.Printf("Net Address: %s\n", o.Address)
	fmt.Printf("Listener Address %s\n", o.ListenerAddress)
	fmt.Printf("Validator: %v\n", o.IsValidator)

}

func getDefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		ID:              GenerateId(),
		Address:         GenerateAddress(),
		ListenerAddress: "0.0.0.0:10690",
		IsValidator:     true,
	}
}

func newTlsOption(cert, key string) (grpc.ServerOption, error) {
	cred, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
		return grpc.Creds(cred), err
	}

	return grpc.Creds(cred), nil
}

func GenerateId() string {
	hostname, _ := os.Hostname()
	return hostname + "@" + util.GetLocalIP()
}

func GenerateAddress() string {
	return util.GetLocalIP()
}
