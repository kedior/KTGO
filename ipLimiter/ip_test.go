package ipLimiter

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"
)

var l = NewIpLimit(
	time.Millisecond*300, time.Second*5, 5,
	func(ip IP) {
		log.Println("ip:", ip, "访问过快,已被临时封禁5秒")
	}, NewIpLimit(
		time.Minute*3, time.Hour, 3,
		func(ip IP) {
			log.Println("ip:", ip, "多次访问过快,已被临时封禁1小时")
		}, nil))

func TestNewIpLimit(t *testing.T) {
	ip := MustIP("127.0.0.1")
	for i := 0; i < 100; i++ {
		fmt.Printf("%v  ", l.Put(ip))
		time.Sleep(time.Millisecond * 290)
	}
}

func TestIpLimiter_Put(t *testing.T) {
	ipLimit := func(fun http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_, ip := NewIP("192.168.0.1")
			if !l.Put(ip) {
				fmt.Fprint(w, "FAST！")
				return
			}
			fun(w, r)
		}
	}
	http.HandleFunc("/hello", ipLimit(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello!")
	}))
	http.ListenAndServe(":8080", nil)
}

func TestNewIP(t *testing.T) {

}
