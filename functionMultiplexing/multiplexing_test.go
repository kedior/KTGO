package functionMultiplexing

import (
	"sync"
	"testing"
)

func TestMakeMultiplex(t *testing.T) {
	var (
		callCount = 1000
		called    = 0
	)

	from := func(id int) int {
		called++
		return 1
	}
	var wrapper func(int) int
	err := MakeMultiplex(&from, &wrapper)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	for i := 0; i < callCount; i++ {
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
