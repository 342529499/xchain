package xcode

import "sync"

var (
	identifier int64         = 1
	locker     *sync.RWMutex = new(sync.RWMutex)
)

func getIdentifier() int64 {
	locker.Lock()
	defer locker.Unlock()

	id := identifier
	identifier++

	return id
}
