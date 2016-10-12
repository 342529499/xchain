package main

import (
	"github.com/1851616111/xchain/cmd"
	"os"
)

func main() {
	xChain := cmd.NewCommandXChain()

	if err := xChain.Execute(); err != nil {
		os.Exit(0)
	}
}
