package NetworkTrafficControl

import (
	"fmt"
	"testing"
	"time"
)

func TestBucket_Test(t *testing.T) {
	t.Parallel()
	c := NewController(10, 50, 2*time.Second)
	c.Start()
	time.Sleep(6 * time.Second)
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 100)
		fmt.Println(c.Test())
	}
}
