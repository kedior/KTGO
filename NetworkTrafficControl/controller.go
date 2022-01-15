package NetworkTrafficControl

import (
	"sync"
	"time"
)

type bucket struct {
	tokens   int64
	once     int64
	capacity int64
	interval time.Duration
	mu       sync.Mutex
	start    bool
}

// NewController 实现一个令牌桶,capacity为容积,每隔interval向桶中投入once个令牌
func NewController(once, capacity int64, interval time.Duration) *bucket {
	return &bucket{
		tokens:   0,
		once:     once,
		capacity: capacity,
		interval: interval,
		start:    false,
	}
}

// Start 启动令牌桶
func (b *bucket) Start() {
	if b.start {
		return
	}
	b.start = true
	go func() {
		for {
			time.Sleep(b.interval)
			b.mu.Lock()
			b.tokens += b.once
			if b.tokens > b.capacity {
				b.tokens = b.capacity
			}
			b.mu.Unlock()
		}
	}()
}

// Test 取出桶中一个令牌,如果桶中剩余令牌，则返回true，否则返回false
func (b *bucket) Test() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}
