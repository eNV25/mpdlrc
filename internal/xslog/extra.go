package xslog

import (
	"reflect"
	"runtime"

	"github.com/gdamore/tcell/v2"

	"github.com/env25/mpdlrc/internal/config"
)

func TypeName(v any) string {
	if config.Debug {
		return reflect.TypeOf(v).String()
	}
	return ""
}

func FuncName(fn any) string {
	if config.Debug {
		return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	}
	return ""
}

func KeyName(key tcell.Key) string {
	if config.Debug {
		return tcell.KeyNames[key]
	}
	return ""
}
