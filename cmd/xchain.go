package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

const (
	XChainLong = `
xchain

The xchain help you deploy, invoke and query SmartContract,
at the same time, xchain supply you with manager and devops command.
To start an xchain node, run:

	xchain start&

`
)

func NewCommandXChain() *cobra.Command {
	out := os.Stdout

	root := &cobra.Command{
		Use:   "xchain",
		Short: "manager your peer and xcontract",
		Long:  XChainLong,
		Run: func(c *cobra.Command, args []string) {
			c.SetOutput(out)
			c.Help()
		},
	}

	startXChain, _ := newCommandStart(out)
	root.AddCommand(startXChain)

	return root
}
