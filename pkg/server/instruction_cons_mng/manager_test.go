package instruction_cons_mng

import (
	"errors"
	"testing"
)

func TestNewConnectionsManager(t *testing.T) {
	manager := NewConnectionsManager(2)
	con, err := manager.Get("123")
	if con != nil && err != ErrConnectionNotExist {
		t.Fatal("new connection manager err")
	}

	var c1, c2 Connection = Connection(nil), Connection(nil)
	if err := manager.Add("c1", c1); err != nil {
		t.Fatalf("add connnection to manager err: %v\n", err)
	}
	if err := manager.Add("c2", c2); err != nil {
		t.Fatalf("add connnection to manager err: %v\n", err)
	}

	newC2, err := manager.Get("c2")
	if err != nil || newC2 != c2 {
		t.Fatal("get connnection from manager err")
	}

	var (
		counter int = 0
		waitAll     = true
	)
	manager.BroadcastFunc(waitAll, func(key string, con Connection) error {
		counter++
		return errors.New("broadcast err")
	})
	if counter != 2 {
		t.Fatalf("connection manager broadcast err")
	}

	counter = 0
	waitAll = false

	manager.BroadcastFunc(waitAll, func(key string, con Connection) error {
		counter++
		return errors.New("broadcast err")
	})

	if counter != 1 {
		t.Fatalf("connection manager broadcast err")
	}

	manager.Del("c1")

	if manager.Exist("c1") {
		t.Fatalf("connection manager del err")
	}
}
