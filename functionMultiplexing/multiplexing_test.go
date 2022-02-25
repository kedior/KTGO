package functionMultiplexing

import (
	"fmt"
	"sync"
	"testing"
)

type foo struct {
	id    int
	value []string
}
type complexStruct struct {
	arr  []foo
	name string
	p    *int8
}

func TestMakeMultiplex(t *testing.T) {
	var (
		callCount = 1000
		called    = 0
	)

	from := func(input complexStruct, str string) int {
		fmt.Println(input)
		called++
		return 1
	}
	var wrapper func(complexStruct, string) int
	err := MakeMultiplex(&from, &wrapper)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	for i := 0; i < callCount; i++ {
		inter := int8(0)
		wg.Add(1)
		go func() {
			defer wg.Done()
			wrapper(complexStruct{
				arr: []foo{
					{1, []string{"123", "456"}},
					{2, []string{"789", "ok"}},
				},
				name: "kedior",
				p:    &inter,
			}, "test")
		}()
	}
	wg.Wait()
}
