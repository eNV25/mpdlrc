//go:build !debug

package xslog

import (
	"time"

	"github.com/env25/mpdlrc/internal/slog"
)

type Args []any

func (a *Args) Any(k string, v any)                {}
func (a *Args) Bool(k string, v bool)              {}
func (a *Args) Duration(k string, v time.Duration) {}
func (a *Args) Float64(k string, v float64)        {}
func (a *Args) Group(k string, v ...slog.Attr)     {}
func (a *Args) Int(k string, v int)                {}
func (a *Args) Int64(k string, v int64)            {}
func (a *Args) String(k, v string)                 {}
func (a *Args) Time(k string, v time.Time)         {}
func (a *Args) Uint64(k string, v uint64)          {}
