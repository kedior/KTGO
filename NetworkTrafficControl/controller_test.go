package NetworkTrafficControl

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBucket_Test(t *testing.T) {
	t.Parallel()
	c := NewController(10, 50, 2*time.Second, time.Second)
	c.Start()
	time.Sleep(6 * time.Second)
	var wg sync.WaitGroup
	wg.Add(400)
	f := func() {
		for i := 0; i < 100; i++ {
			time.Sleep(time.Millisecond * 100)
			fmt.Println(c.Test())
			wg.Done()
		}
	}
	go f()
	go f()
	go f()
	go f()
	wg.Wait()
}
