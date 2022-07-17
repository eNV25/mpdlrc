package panics

import (
	"context"
)

type key struct{}

func ContextWithHook(ctx context.Context, h func()) context.Context {
	v := ctx.Value(key{}) // *[]func()
	if v == nil {
		v = new([]func())
		ctx = context.WithValue(ctx, key{}, v)
	}
	p := v.(*[]func())
	*p = append(*p, h)
	return ctx
}

func RunHooksFromContext(ctx context.Context) {
	v := ctx.Value(key{})
	if v == nil {
		return
	}
	for _, f := range *v.(*[]func()) {
		f()
	}
}
