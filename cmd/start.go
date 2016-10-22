package cmd

import (
	"io"
	"log"

	"github.com/spf13/cobra"

	"github.com/1851616111/xchain/pkg/server"
	"github.com/1851616111/xchain/pkg/util"
	"github.com/1851616111/xchain/pkg/util/file"

)

type StartOptions struct {
	peer *PeerOptions
}

type PeerOptions struct {
	id            string
	netAddress    string
	listenAddress string
	isValidator   bool

	entryPointAddress string

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

	flags.StringVar(&options.peer.id, "id", server.GenerateId(), "The node identity consist of hostname@ip .etc michael@192.168.1.1")
	flags.StringVar(&options.peer.netAddress, "netAddress", util.GetLocalIP(), "The node net address for other node to connect.")
	flags.StringVar(&options.peer.listenAddress, "listenAddress", "0.0.0.0:10690", "The address to lister for grpc connections")
	flags.BoolVar(&options.peer.isValidator, "validator", true, "peer service type, validate peer or non-validate peer")

	flags.StringVar(&options.peer.entryPointAddress, "entryPointAddress", "", "The address to join the network")

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
	options.Complete()

	serverOptions := &server.ServerOptions{
		ID:              options.peer.id,
		Address:         options.peer.netAddress,
		ListenerAddress: options.peer.listenAddress,
		IsValidator:     options.peer.isValidator,

		EntryPointAddress: options.peer.entryPointAddress,

		TlsEnabled:   options.peer.tlsEnabled,
		CertFilePath: options.peer.tlsCertPath,
		KeyFilePath:  options.peer.tlsKeyPath,
	}

	server.NewAndStartGrpcServer(serverOptions)
}

func (options *StartOptions) Complete () error {
	if options.peer.netAddress != "" {
		util.SetLocalIP(options.peer.netAddress)
		options.peer.id = server.GenerateId()
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


