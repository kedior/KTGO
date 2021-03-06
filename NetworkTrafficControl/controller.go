package NetworkTrafficControl

import (
	"sync"
	"time"
)

type bucket struct {
	once     int64
	max      int64
	c        int64
	last     time.Duration
	interval time.Duration
	mu       sync.Mutex
}

// NewController 实现一个令牌桶,capacity为容积,每隔interval向桶中投入once个令牌
func NewController(once, max int64, interval time.Duration) *bucket {
	return &bucket{
		once:     once,
		max:      max,
		c:        max,
		interval: interval,
		last:     time.Duration(time.Now().UnixNano()),
	}
}

func (b *bucket) update() {
	now := time.Duration(time.Now().UnixNano())
	if now-b.last < b.interval {
		return
	}
	if b.c += int64((now-b.last)/b.interval) * b.once; b.c > b.max {
		b.c = b.max
	}
	b.last = now
}

// Test 取出桶中一个令牌,如果桶中剩余令牌，则返回true
//若无令牌,返回false
func (b *bucket) Test() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.update()
	if b.c > 0 {
		b.c--
		return true
	}
	return false
}

// NewControlFunc 是令牌桶的闭包实现
//调用返回的函数取出桶中一个令牌,如果桶中剩余令牌，则返回true
//若无令牌,返回false
func NewControlFunc(once, max int64, interval time.Duration) func() bool {
	var (
		c    = max
		last = time.Duration(time.Now().UnixNano())
		mu   sync.Mutex
	)
	return func() bool {
		mu.Lock()
		defer mu.Unlock()
		now := time.Duration(time.Now().UnixNano())
		if now-last < interval {
			goto end
		}
		if c += int64((now-last)/interval) * once; c > max {
			c = max
		}
		last = now
	end:
		if c > 0 {
			c--
			return true
		}
		return false

	}
}
