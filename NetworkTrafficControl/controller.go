package NetworkTrafficControl

import (
	"time"
)

type bucket struct {
	once      int64
	container chan struct{}
	interval  time.Duration
	waitTime  time.Duration
	start     bool
}

// NewController 实现一个令牌桶,capacity为容积,每隔interval向桶中投入once个令牌
func NewController(once, contain int64, interval, waitTime time.Duration) *bucket {
	return &bucket{
		once:      once,
		container: make(chan struct{}, contain),
		interval:  interval,
		waitTime:  waitTime,
		start:     false,
	}
}

// Start 启动令牌桶
func (b *bucket) Start() {
	if b.start {
		return
	}
	b.start = true
	go func() {
	loop:
		time.Sleep(b.interval)
	send:
		for i := 0; i < int(b.once); i++ {
			select {
			case b.container <- struct{}{}:
			default:
				break send
			}
		}
		goto loop
	}()
}

// Test 取出桶中一个令牌,如果桶中剩余令牌，则返回true\
//否则会先等待waitTime,若仍然无令牌,返回false
func (b *bucket) Test() bool {
	select {
	case <-b.container:
		return true
	case <-time.Tick(b.waitTime):
		return false
	}
}
