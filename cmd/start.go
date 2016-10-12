package cmd

import (
	"io"
	"log"
	"net"

	"github.com/1851616111/xchain/pkg/util/file"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

type StartOptions struct {
	peer *PeerOptions
}

type PeerOptions struct {
	listenAddress string

	//
	tlsEnabled  bool
	tlsCertPath string
	tlsKeyPath  string
}

func newCommandStart(out io.Writer) (*cobra.Command, *StartOptions) {
	options := &StartOptions{
		peer: &PeerOptions{},
	}

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Launch one of xchain peer node",
		Long:  "Launch one of xchain peer node",
		Run: func(c *cobra.Command, args []string) {
			options.Run(c, args)

		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.peer.listenAddress, "listen", "0.0.0.0:10600", "The address to lister for grpc connections")
	flags.BoolVar(&options.peer.tlsEnabled, "tls.enabled", false, "TLS settings for p2p communications")
	flags.StringVar(&options.peer.tlsCertPath, "tls.cert", "", "TLS cert file")
	flags.StringVar(&options.peer.tlsKeyPath, "tls.key", "", "TLS key file")

	return cmd, options
}

func (options *StartOptions) Run(c *cobra.Command, args []string) {
	if err := options.Validate(args); err != nil {
		log.Printf("xchart start validate args err, %v.\n", err)
		return
	}

	if err := options.Start(); err != nil {
		log.Printf("xchart start err, %v.\n", err)
	}
}

func (options *StartOptions) Start() error {

	lis, err := net.Listen("tcp", options.peer.listenAddress)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	if options.peer.tlsEnabled {
		opt := newTlsOption(options.peer.tlsCertPath, options.peer.tlsKeyPath)
		opts = append(opts, opt)
	}

	if err := grpc.NewServer(opts...).Serve(lis); err != nil {
		return err
	}
	return nil
}

func (options *StartOptions) Validate(args []string) error {
	var err error

	if options.peer.tlsEnabled {
		if err = file.IsFileExist(options.peer.tlsCertPath); err != nil {
			return err
		}

		if err = file.IsFileExist(options.peer.tlsKeyPath); err != nil {
			return err
		}
	}

	return nil
}

func newTlsOption(cert, key string) grpc.ServerOption {

	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}

	return grpc.Creds(creds)

}
