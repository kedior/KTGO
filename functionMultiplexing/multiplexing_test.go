package functionMultiplexing

import (
	"fmt"
	"testing"
	"time"
)

func TestMakeMultiplex(t *testing.T) {
	complex := func() {
		time.Sleep(time.Second * 2)
		fmt.Println("jksghhjkas")
	}
	var superFun func()
	err := MakeMultiplex(&complex, &superFun)
	if err != nil {
		panic(err)
	}
	go superFun()
	go superFun()
	go superFun()
	go superFun()
	go superFun()
	go superFun()
	time.Sleep(time.Second * 7)
}
