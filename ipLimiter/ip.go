package ipLimiter

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	InvalidIpErr = errors.New("invalid ip string")
)

type IP uint32

func NewIP(str string) (error, IP) {
	sp := strings.Split(str, ".")
	if len(sp) != 4 {
		return InvalidIpErr, 0
	}
	ip := uint32(0)
	for i := 3; i >= 0; i-- {
		v, e := strconv.Atoi(sp[i])
		if e != nil || v > 255 {
			return InvalidIpErr, 0
		}
		ip <<= 8
		ip |= uint32(v)
	}
	return nil, IP(ip)
}

func MustIP(str string) IP {
	err, ip := NewIP(str)
	if err != nil {
		panic(err)
	}
	return ip
}

func (ip IP) String() string {
	bf := bytes.NewBufferString("")
	t := uint32(ip)
	for i := 0; i < 4; i++ {
		bf.WriteString(strconv.Itoa(int(t & 255)))
		bf.WriteRune('.')
		t >>= 8
	}
	bf.Truncate(bf.Len() - 1)
	return bf.String()
}

type limiter struct {
	limitTo int64
	lastAsk int64
	cnt     uint8
}

type ipLimiter struct {
	limitTime time.Duration
	banTime   time.Duration
	repeat    uint8
	m         map[IP]*limiter
	mu        sync.Mutex
	super     *ipLimiter
	banAction func(ip IP)
}

// NewIpLimit 创建一个新的limiter。
//如果一个ip在limitTime内访问两次以上，将会进行计数。
//当连续计数超过repeat时，此ip将会被ban掉，banTime为的被ban的时间。
//ip被ban时，limiter会调用一次action，并且会向父级limiter的Put方法传入一次ip。
//super为父级limiter的指针
func NewIpLimit(limitTime, banTime time.Duration,
	repeat uint8, action func(IP),
	super *ipLimiter) *ipLimiter {
	return &ipLimiter{
		limitTime: limitTime,
		repeat:    repeat,
		banTime:   banTime,
		m:         make(map[IP]*limiter),
		super:     super,
		banAction: action,
	}
}

// IsBanning 返回ip在时间戳为time的时刻是否被封禁
func (l *ipLimiter) IsBanning(ip IP, time int64) bool {
	return (l.m[ip] != nil && l.m[ip].limitTo > time) ||
		(l.super != nil && l.super.IsBanning(ip, time))

}

// Put 提醒limiter此ip进行了一次访问
//如果此一个p多次访问，超过limiter的限制，将会ban掉该ip
//如果ip已经被ban了，返回false，否则返回true
func (l *ipLimiter) Put(ip IP) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now().UnixNano()
	il := l.m[ip]
	if il == nil {
		l.m[ip] = &limiter{lastAsk: now}
		return true
	}
	if l.IsBanning(ip, now) {
		return false
	}
	il.lastAsk, now = now, il.lastAsk
	if il.lastAsk-now < l.limitTime.Nanoseconds() {
		if il.cnt >= l.repeat-1 {
			il.cnt = 0
			il.limitTo = now + l.banTime.Nanoseconds()
			l.banAction(ip)
			if l.super != nil {
				l.super.Put(ip)
			}
			return false
		}
		il.cnt++
		return true
	}
	il.cnt = 0
	return true
}
