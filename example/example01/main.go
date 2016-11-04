package main

import (
	"github.com/1851616111/xchain/pkg/xcode/broker"
)

type example struct {
}

func (e *example) Init(i broker.Instructions, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (e *example) Invoke(i broker.Instructions, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (e *example) Query(i broker.Instructions, function string, args []string) ([]byte, error) {
	return nil, nil
}

func main() {
	broker.StartCoderService(&example{})
}
