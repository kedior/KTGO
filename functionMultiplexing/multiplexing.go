package functionMultiplexing

import (
	"errors"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type operator struct {
	input interface{}
	doing int32
	out   chan []reflect.Value
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

func makeInStruct(tp reflect.Type) func([]reflect.Value) interface{} {
	numIn := tp.NumIn()
	inFields := make([]reflect.StructField, numIn)
	for i := 0; i < numIn; i++ {
		inFields[i] = reflect.StructField{
			Name: "I" + strconv.Itoa(i),
			Type: tp.In(i),
		}
	}
	structType := reflect.StructOf(inFields)
	return func(args []reflect.Value) interface{} {
		elem := reflect.New(structType).Elem()
		num := elem.NumField()
		for i := 0; i < num; i++ {
			elem.Field(i).Set(args[i])
		}
		return elem.Interface()
	}
}

func todo(fromFunc *reflect.Value) func([]reflect.Value) []reflect.Value {
	inToStructFn := makeInStruct(fromFunc.Type())
	m := make(map[interface{}]*operator)
	var lock sync.Mutex

	return func(in []reflect.Value) []reflect.Value {
	start:
		lock.Lock()
		inputKey := inToStructFn(in)
		op := m[inputKey]
		if op == nil {
			op = &operator{
				input: inputKey,
				doing: 0,
				out:   make(chan []reflect.Value),
			}
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
