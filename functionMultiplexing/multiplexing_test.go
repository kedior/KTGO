package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestMakeMultiplex(t *testing.T) {
	complex := func(i int, str string) (string, int) {
		for i < 100000 {
			i *= i
			str += str + strconv.Itoa(i)
			time.Sleep(time.Second)
		}
		return str, i
	}
	var superFun func(int, string) (string, int)
	err := MakeMultiplex(&complex, &superFun)
	if err != nil {
		panic(err)
	}
	go fmt.Println(superFun(3, "kedior"))
	go fmt.Println(superFun(3, "kedior"))
	go fmt.Println(superFun(3, "kedior"))
	go fmt.Println(superFun(3, "kedior"))
	go fmt.Println(superFun(3, "kedior"))
	go fmt.Println(superFun(3, "kedior"))
	go fmt.Println(superFun(3, "kedior"))
	go fmt.Println(superFun(3, "kedior"))
	go fmt.Println(superFun(3, "kedior"))
	go fmt.Println(superFun(3, "kedior"))
	go fmt.Println(superFun(3, "kedior"))
	go fmt.Println(superFun(3, "kedior"))
	go fmt.Println(superFun(3, "kedior"))
	go fmt.Println(superFun(3, "kedior"))
	go fmt.Println(superFun(3, "kedior"))

	time.Sleep(time.Second * 15)
}
