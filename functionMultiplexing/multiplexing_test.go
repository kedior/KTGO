package functionMultiplexing

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestMakeMultiplex(t *testing.T) {
	var (
		callCount = int64(1000)
		called    = int64(0)
	)

	from := func(id int) int {
		atomic.AddInt64(&called, 1)
		return 1
	}
	var wrapper func(int) int
	err := MakeMultiplex(&from, &wrapper)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	for i := int64(0); i < callCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wrapper(1)
		}()
	}
	wg.Wait()
	if called >= callCount {
		t.Fatal("multiplex error")
	}
}
