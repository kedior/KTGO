package main

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	unverifiedCodeRemaining = errors.New("存在一个未验证的验证码")
)

type VerifyCode int32
type MailAddr string

type pack struct {
	code   VerifyCode
	finish chan struct{}
}

type CodeSender struct {
	m        map[MailAddr]*pack
	exp      time.Duration
	mu       sync.Mutex
	sendFunc func(MailAddr, VerifyCode) error
}

// NewCode 生成一个随机四位验证码
func NewCode() VerifyCode {
	rand.Seed(time.Now().UnixNano() ^ 345674756678)
	i := rand.Int31() % 10000
	if i < 1000 {
		return NewCode()
	}
	return VerifyCode(i)
}

// NewVerifyCoder 创建一个验证码发送器，exp是每个验证码的过期时间
//sendFunc是发送邮件时调用的函数
func NewVerifyCoder(exp time.Duration, sendFunc func(MailAddr, VerifyCode) error) *CodeSender {
	return &CodeSender{
		m:        make(map[MailAddr]*pack),
		exp:      exp,
		sendFunc: sendFunc,
	}
}

// SendNewCode 向指定邮箱发送一个验证码
//如果此邮箱已经发送过验证码并未过期，或者发送失败，则抛出错误
func (m *CodeSender) SendNewCode(email MailAddr) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.m[email] != nil {
		return unverifiedCodeRemaining
	}
	return m.unLockSendNewCode(email)
}

// MustSendNewCode 向指定邮箱发送一个验证码
//如果此邮箱已经发送过验证码并未过期，则强制替换之前的验证码并结束它的生命周期
//如果发送失败会抛出错误
func (m *CodeSender) MustSendNewCode(email MailAddr) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	p := m.m[email]
	if p != nil {
		//发送两次信号
		//第一个信号让map删除
		//最后一个信号用来等待map删除完毕
		p.finish <- struct{}{}
		p.finish <- struct{}{}
	}
	return m.unLockSendNewCode(email)
}

// Verify 验证邮箱对应的验证码是否正确
//验证成功后此验证码生命周期结束
func (m *CodeSender) Verify(email MailAddr, c VerifyCode) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	p := m.m[email]
	if p == nil || p.code != c {
		return false
	}
	//发送两次信号
	//第一个信号让map删除
	//最后一个信号用来等待map删除完毕
	p.finish <- struct{}{}
	p.finish <- struct{}{}
	return true

}

//不带锁的发送验证码，并创建倒计时携程
func (m *CodeSender) unLockSendNewCode(email MailAddr) error {
	newCode := NewCode()
	if err := m.sendFunc(email, newCode); err != nil {
		return err
	}
	m.m[email] = &pack{newCode, make(chan struct{})}
	//用来等待携程对map的操作
	start := make(chan struct{})
	go func(codeMap *CodeSender, em MailAddr) {
		f := codeMap.m[email].finish
		//map操作完毕后让父亲线程结束，来释放锁
		start <- struct{}{}
		select {
		case <-time.Tick(codeMap.exp):
			{
				//锁是因为此处被调用时是完全的未知状态
				m.mu.Lock()
				delete(codeMap.m, em)
				log.Println(em, "的验证码已过期")
				m.mu.Unlock()
			}
		case <-f:
			{
				//不锁是因为运行到此时必在运行有锁的Verify或MustSendNewCode函数
				//如果加上锁会导致此处必在Verify或MustSendNewCode函数执行完再抢占锁
				//如果此时有另一携程也在抢占锁，会导致此处被滞后执行，导致错误
				//只需要让Verify或MustSendNewCode函数等待信号，等此处运行完毕，让解锁
				//避免Verify或MustSendNewCode函数关锁后此处仍在无锁运行
				delete(codeMap.m, em)
				log.Println(em, "的验证码已被验证或被强制替换")
				<-f
			}
		}
	}(m, email)
	<-start
	return nil
}
