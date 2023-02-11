package panics

import (
	"context"
)

type key struct{}

// ContextWithHook returns a [context.Context] after adding the panic handler hook h.
func ContextWithHook(ctx context.Context, hs ...func()) context.Context {
	v := ctx.Value(key{}) // *[]func()
	if v == nil {
		v = new([]func())
		ctx = context.WithValue(ctx, key{}, v)
	}
	p := v.(*[]func())
	*p = append(*p, hs...)
	return ctx
}

func runHooksFromContext(ctx context.Context) {
	v := ctx.Value(key{})
	if v == nil {
		return
	}
	for _, f := range *v.(*[]func()) {
		f()
	}
}
