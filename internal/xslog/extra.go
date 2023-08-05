package xslog

import (
	"reflect"
	"runtime"

	"github.com/gdamore/tcell/v2"
)

func TypeName[T any](t T) typeName[T] {
	return typeName[T]{}
}

type typeName[T any] struct{}

func (typeName[T]) String() string {
	var t T
	return reflect.TypeOf(t).String()
}

func FuncName(fn any) funcName {
	return funcName(reflect.ValueOf(fn).UnsafePointer())
}

type funcName uintptr

func (fn funcName) String() string {
	return runtime.FuncForPC(uintptr(fn)).Name()
}

type Key rune

func (key Key) String() string {
	if key < 0 {
		return tcell.KeyNames[tcell.Key(-key)]
	}
	return string(key)
}
