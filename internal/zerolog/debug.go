//go:build debug

package zerolog

import (
	"context"
	"io"

	"github.com/rs/zerolog"
)

type Array = zerolog.Array

func Arr() *Array {
	return zerolog.Arr()
}

type Formatter = zerolog.Formatter

type ConsoleWriter = zerolog.ConsoleWriter

func NewConsoleWriter(options ...func(w *ConsoleWriter)) ConsoleWriter {
	return zerolog.NewConsoleWriter(options...)
}

type Context = zerolog.Context

func Ctx(ctx context.Context) *Logger {
	return zerolog.Ctx(ctx)
}

type Event = zerolog.Event

type LogObjectMarshaler = zerolog.LogObjectMarshaler

type LogArrayMarshaler = zerolog.LogArrayMarshaler

func Dict() *Event {
	return zerolog.Dict()
}

const (
	TimeFormatUnix = zerolog.TimeFormatUnix

	TimeFormatUnixMs = zerolog.TimeFormatUnixMs

	TimeFormatUnixMicro = zerolog.TimeFormatUnixMicro

	TimeFormatUnixNano = zerolog.TimeFormatUnixNano
)

func SetGlobalLevel(l Level) {
	zerolog.SetGlobalLevel(l)
}

func GlobalLevel() Level {
	return Disabled
}

func DisableSampling(v bool) {
}

type Hook = zerolog.Hook

type HookFunc = zerolog.HookFunc

type LevelHook = zerolog.LevelHook

func NewLevelHook() LevelHook {
	return zerolog.NewLevelHook()
}

type Level = zerolog.Level

const (
	DebugLevel Level = iota

	InfoLevel

	WarnLevel

	ErrorLevel

	FatalLevel

	PanicLevel

	NoLevel

	Disabled

	TraceLevel Level = -1
)

func ParseLevel(levelStr string) (Level, error) {
	return 0, nil
}

type Logger = zerolog.Logger

func New(w io.Writer) Logger {
	return zerolog.New(w)
}

func Nop() Logger {
	return zerolog.Nop()
}

type Sampler = zerolog.Sampler

type RandomSampler = zerolog.RandomSampler

type BasicSampler = zerolog.BasicSampler

type BurstSampler = zerolog.BurstSampler

type LevelSampler = zerolog.LevelSampler

type SyslogWriter = zerolog.SyslogWriter

func SyslogLevelWriter(w SyslogWriter) LevelWriter {
	return zerolog.SyslogLevelWriter(w)
}

func SyslogCEEWriter(w SyslogWriter) LevelWriter {
	return zerolog.SyslogCEEWriter(w)
}

type LevelWriter = zerolog.LevelWriter

func SyncWriter(w io.Writer) io.Writer {
	return zerolog.SyncWriter(w)
}

func MultiLevelWriter(writers ...io.Writer) LevelWriter {
	return zerolog.MultiLevelWriter(writers...)
}

type TestingLog = zerolog.TestingLog

type TestWriter = zerolog.TestWriter

func NewTestWriter(t TestingLog) TestWriter {
	return zerolog.NewTestWriter(t)
}

func ConsoleTestWriter(t TestingLog) func(w *ConsoleWriter) {
	return zerolog.ConsoleTestWriter(t)
}
