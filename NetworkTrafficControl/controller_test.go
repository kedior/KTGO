package NetworkTrafficControl

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBucket_Test(t *testing.T) {
	c := NewController(10, 50, 2*time.Second)
	var wg sync.WaitGroup
	wg.Add(400)
	f := func() {
		for i := 0; i < 100; i++ {
			time.Sleep(time.Millisecond * 100)
			fmt.Println(c.c, c.Test())
			wg.Done()
		}
	}
	go f()
	go f()
	go f()
	go f()
	wg.Wait()
}

func TestNewControllFunc(t *testing.T) {
	pc := NewControlFunc(10, 50, 2*time.Second).([]interface{})
	c := pc[0].(func() bool)
	var wg sync.WaitGroup
	wg.Add(400)
	f := func() {
		for i := 0; i < 100; i++ {
			time.Sleep(time.Millisecond * 100)
			fmt.Println(c())
			wg.Done()
		}
	}
	go f()
	go f()
	go f()
	go f()
	wg.Wait()
}
