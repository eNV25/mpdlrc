package timerpool_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/env25/mpdlrc/internal/timerpool"
)

func Example() {
	var wg sync.WaitGroup

	ctx1 := context.TODO()

	wg.Add(2)

	timer1 := timerpool.Get(time.Second)
	go func() {
		defer wg.Done()

		select {
		case <-ctx1.Done():
			timerpool.Put(timer1, false)
			return
		case <-timer1.C:
			timerpool.Put(timer1, true)
		}

		fmt.Println("1")
	}()

	ctx2, cancel2 := context.WithCancel(ctx1)

	timer2 := timerpool.Get(time.Second)
	go func() {
		defer wg.Done()

		select {
		case <-ctx2.Done():
			timerpool.Put(timer2, false)
			return
		case <-timer2.C:
			timerpool.Put(timer2, true)
		}

		fmt.Println("2")
	}()

	cancel2()

	wg.Wait()

	// OUTPUT:
	// 1
}
