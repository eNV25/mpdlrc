package timerpool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func Example() {
	var wg sync.WaitGroup

	ctx1 := context.TODO()

	wg.Add(2)

	timer1 := Get(time.Second)
	go func() {
		defer wg.Done()

		select {
		case <-ctx1.Done():
			Put(timer1, false)
			return
		case <-timer1.C:
			Put(timer1, true)
		}

		fmt.Println("1")
	}()

	ctx2, cancel2 := context.WithCancel(ctx1)

	timer2 := Get(time.Second)
	go func() {
		defer wg.Done()

		select {
		case <-ctx2.Done():
			Put(timer2, false)
			return
		case <-timer2.C:
			Put(timer2, true)
		}

		fmt.Println("2")
	}()

	cancel2()

	wg.Wait()

	// OUTPUT:
	// 1
}
