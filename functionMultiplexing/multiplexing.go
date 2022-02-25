package functionMultiplexing

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type operator struct {
	input string
	doing int32
	out   chan []reflect.Value
}

var OpPool = sync.Pool{
	New: func() interface{} {
		return &operator{}
	},
}

func MakeMultiplex(fromFnPtr, toFnPtr interface{}) error {
	from := reflect.ValueOf(fromFnPtr).Elem()
	to := reflect.ValueOf(toFnPtr).Elem()
	if from.Type().Kind() != reflect.Func {
		return errors.New("input args type must be func")
	}
	if from.Type().Kind() != to.Type().Kind() {
		return errors.New("missing match type")
	}
	v := reflect.MakeFunc(from.Type(), todo(&from))
	to.Set(v)
	return nil
}

func todo(fromFunc *reflect.Value) func([]reflect.Value) []reflect.Value {
	m := make(map[string]*operator)
	var lock sync.Mutex

	return func(in []reflect.Value) []reflect.Value {
	start:
		lock.Lock()
		inputKey := fmt.Sprint(in)
		op := m[inputKey]
		if op == nil {
			op = OpPool.Get().(*operator)
			op.input = inputKey
			op.doing = 0
			op.out = make(chan []reflect.Value)
			m[inputKey] = op
		}
		lock.Unlock()
		if atomic.CompareAndSwapInt32(&op.doing, 0, 1) {
			go func() {
				result := fromFunc.Call(in)
				after := time.After(time.Second * 5)
				for {
					select {
					case op.out <- result:
					case <-after:
						lock.Lock()
						delete(m, op.input)
						close(op.out)
						OpPool.Put(op)
						lock.Unlock()
						return
					}
				}
			}()
		}
		if v, ok := <-op.out; ok {
			return v
		}
		goto start
	}
}
