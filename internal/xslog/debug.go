//go:build debug

package xslog

import (
	"time"

	"github.com/env25/mpdlrc/internal/slog"
)

type Args []any

func (a *Args) Any(k string, v any)                { *a = append(*a, slog.Any(k, v)) }
func (a *Args) Bool(k string, v bool)              { *a = append(*a, slog.Bool(k, v)) }
func (a *Args) Duration(k string, v time.Duration) { *a = append(*a, slog.Duration(k, v)) }
func (a *Args) Float64(k string, v float64)        { *a = append(*a, slog.Float64(k, v)) }
func (a *Args) Group(k string, v ...slog.Attr)     { *a = append(*a, slog.Group(k, v...)) }
func (a *Args) Int(k string, v int)                { *a = append(*a, slog.Int(k, v)) }
func (a *Args) Int64(k string, v int64)            { *a = append(*a, slog.Int64(k, v)) }
func (a *Args) String(k, v string)                 { *a = append(*a, slog.String(k, v)) }
func (a *Args) Time(k string, v time.Time)         { *a = append(*a, slog.Time(k, v)) }
func (a *Args) Uint64(k string, v uint64)          { *a = append(*a, slog.Uint64(k, v)) }

func (a *Args) Rune(k string, v rune) { *a = append(*a, slog.String(k, string(v))) }
