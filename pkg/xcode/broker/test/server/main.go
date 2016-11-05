package main

import (
	"fmt"
	"github.com/1851616111/xchain/pkg/xcode/broker"
)

type example struct {
}

func (e *example) Init(i broker.Instructions, function string, args []string) ([]byte, error) {
	fmt.Printf("-----> function:%s", function)
	fmt.Printf("-----> args:%s", args)

	b, err := i.GetState("123")
	if err != nil {
		fmt.Println("get state -------->err", err)
	}

	fmt.Printf("get state response %v\n", string(b))
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
