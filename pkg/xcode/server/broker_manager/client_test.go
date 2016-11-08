package server

import "testing"

func Test_connectBroker(t *testing.T) {
	if _, err := connectBroker("10692"); err != nil {
		t.Fatalf("connect to broker err %v\n", err)
	}
}
