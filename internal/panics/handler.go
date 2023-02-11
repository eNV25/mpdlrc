// Package panics implements a panic handler.
package panics

import (
	"context"
	"fmt"
	"os"
	"runtime"
)

// Handle handles panic in the current goroutine. Should be called with defer.
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
	fmt.Fprintf(os.Stderr, "panic: %+#v\n\n%s", r, buf)
}
