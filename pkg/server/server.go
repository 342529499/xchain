package server

import (
	"errors"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/1851616111/xchain/pkg/protos"
)

func NewAndStartGrpcServer(option *ServerOptions) error {
	if option == nil {
		option = getDefaultServerOptions()
	} else if err := option.validate(); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", option.Address)
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

	pb.RegisterNetServer(server, &netServer{})

	if err := server.Serve(lis); err != nil {
		return err
	}

	return nil
}

type ServerOptions struct {
	Address      string
	TlsEnabled   bool
	CertFilePath string
	KeyFilePath  string
}

func (opts *ServerOptions) validate() error {
	if len(strings.TrimSpace(opts.Address)) == 0 {
		return errors.New("param listen address must not be nil.")
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

func getDefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		Address: "0.0.0.0:10690",
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
