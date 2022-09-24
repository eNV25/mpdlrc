package panics

import (
	"context"
	"log"
	"runtime"
)

func Handle(ctx context.Context) {
	r := recover()
	if r == nil {
		return
	}
	runHooksFromContext(ctx)
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	log.Printf("\npanic: %v\n%s\n", r, buf)
}
