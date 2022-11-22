//go:build !debug

package zerolog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type Array struct{}

func Arr() *Array {
	return nil
}

func (*Array) MarshalZerologArray(*Array) {
}

func (a *Array) Object(obj LogObjectMarshaler) *Array {
	return a
}

func (a *Array) Str(val string) *Array {
	return a
}

func (a *Array) Bytes(val []byte) *Array {
	return a
}

func (a *Array) Hex(val []byte) *Array {
	return a
}

func (a *Array) RawJSON(val []byte) *Array {
	return a
}

func (a *Array) Err(err error) *Array {
	return a
}

func (a *Array) Bool(b bool) *Array {
	return a
}

func (a *Array) Int(i int) *Array {
	return a
}

func (a *Array) Int8(i int8) *Array {
	return a
}

func (a *Array) Int16(i int16) *Array {
	return a
}

func (a *Array) Int32(i int32) *Array {
	return a
}

func (a *Array) Int64(i int64) *Array {
	return a
}

func (a *Array) Uint(i uint) *Array {
	return a
}

func (a *Array) Uint8(i uint8) *Array {
	return a
}

func (a *Array) Uint16(i uint16) *Array {
	return a
}

func (a *Array) Uint32(i uint32) *Array {
	return a
}

func (a *Array) Uint64(i uint64) *Array {
	return a
}

func (a *Array) Float32(f float32) *Array {
	return a
}

func (a *Array) Float64(f float64) *Array {
	return a
}

func (a *Array) Time(t time.Time) *Array {
	return a
}

func (a *Array) Dur(d time.Duration) *Array {
	return a
}

func (a *Array) Interface(i interface{}) *Array {
	return a
}

func (a *Array) IPAddr(ip net.IP) *Array {
	return a
}

func (a *Array) IPPrefix(pfx net.IPNet) *Array {
	return a
}

func (a *Array) MACAddr(ha net.HardwareAddr) *Array {
	return a
}

func (a *Array) Dict(dict *Event) *Array {
	return a
}

type Formatter func(interface{}) string

type ConsoleWriter struct {
	Out                 io.Writer
	FormatMessage       Formatter
	FormatFieldValue    Formatter
	FormatExtra         func(map[string]interface{}, *bytes.Buffer) error
	FormatErrFieldValue Formatter
	FormatErrFieldName  Formatter
	FormatTimestamp     Formatter
	FormatLevel         Formatter
	FormatCaller        Formatter
	FormatFieldName     Formatter
	TimeFormat          string
	FieldsExclude       []string
	PartsExclude        []string
	PartsOrder          []string
	NoColor             bool
}

func NewConsoleWriter(options ...func(w *ConsoleWriter)) ConsoleWriter {
	return ConsoleWriter{}
}

func (w ConsoleWriter) Write(p []byte) (n int, err error) {
	return
}

type Context struct{}

func (c Context) Logger() Logger {
	return Logger{}
}

func (c Context) Fields(fields interface{}) Context {
	return c
}

func (c Context) Dict(key string, dict *Event) Context {
	return c
}

func (c Context) Array(key string, arr LogArrayMarshaler) Context {
	return c
}

func (c Context) Object(key string, obj LogObjectMarshaler) Context {
	return c
}

func (c Context) EmbedObject(obj LogObjectMarshaler) Context {
	return c
}

func (c Context) Str(key, val string) Context {
	return c
}

func (c Context) Strs(key string, vals []string) Context {
	return c
}

func (c Context) Stringer(key string, val fmt.Stringer) Context {
	return c
}

func (c Context) Bytes(key string, val []byte) Context {
	return c
}

func (c Context) Hex(key string, val []byte) Context {
	return c
}

func (c Context) RawJSON(key string, b []byte) Context {
	return c
}

func (c Context) AnErr(key string, err error) Context {
	return c
}

func (c Context) Errs(key string, errs []error) Context {
	return c
}

func (c Context) Err(err error) Context {
	return c
}

func (c Context) Bool(key string, b bool) Context {
	return c
}

func (c Context) Bools(key string, b []bool) Context {
	return c
}

func (c Context) Int(key string, i int) Context {
	return c
}

func (c Context) Ints(key string, i []int) Context {
	return c
}

func (c Context) Int8(key string, i int8) Context {
	return c
}

func (c Context) Ints8(key string, i []int8) Context {
	return c
}

func (c Context) Int16(key string, i int16) Context {
	return c
}

func (c Context) Ints16(key string, i []int16) Context {
	return c
}

func (c Context) Int32(key string, i int32) Context {
	return c
}

func (c Context) Ints32(key string, i []int32) Context {
	return c
}

func (c Context) Int64(key string, i int64) Context {
	return c
}

func (c Context) Ints64(key string, i []int64) Context {
	return c
}

func (c Context) Uint(key string, i uint) Context {
	return c
}

func (c Context) Uints(key string, i []uint) Context {
	return c
}

func (c Context) Uint8(key string, i uint8) Context {
	return c
}

func (c Context) Uints8(key string, i []uint8) Context {
	return c
}

func (c Context) Uint16(key string, i uint16) Context {
	return c
}

func (c Context) Uints16(key string, i []uint16) Context {
	return c
}

func (c Context) Uint32(key string, i uint32) Context {
	return c
}

func (c Context) Uints32(key string, i []uint32) Context {
	return c
}

func (c Context) Uint64(key string, i uint64) Context {
	return c
}

func (c Context) Uints64(key string, i []uint64) Context {
	return c
}

func (c Context) Float32(key string, f float32) Context {
	return c
}

func (c Context) Floats32(key string, f []float32) Context {
	return c
}

func (c Context) Float64(key string, f float64) Context {
	return c
}

func (c Context) Floats64(key string, f []float64) Context {
	return c
}

func (c Context) Timestamp() Context {
	return c
}

func (c Context) Time(key string, t time.Time) Context {
	return c
}

func (c Context) Times(key string, t []time.Time) Context {
	return c
}

func (c Context) Dur(key string, d time.Duration) Context {
	return c
}

func (c Context) Durs(key string, d []time.Duration) Context {
	return c
}

func (c Context) Interface(key string, i interface{}) Context {
	return c
}

func (c Context) Caller() Context {
	return c
}

func (c Context) CallerWithSkipFrameCount(skipFrameCount int) Context {
	return c
}

func (c Context) Stack() Context {
	return c
}

func (c Context) IPAddr(key string, ip net.IP) Context {
	return c
}

func (c Context) IPPrefix(key string, pfx net.IPNet) Context {
	return c
}

func (c Context) MACAddr(key string, ha net.HardwareAddr) Context {
	return c
}

func (l Logger) WithContext(ctx context.Context) context.Context {
	return context.TODO()
}

func Ctx(ctx context.Context) *Logger {
	return nil
}

type Event struct{}

type LogObjectMarshaler interface {
	MarshalZerologObject(e *Event)
}

type LogArrayMarshaler interface {
	MarshalZerologArray(a *Array)
}

func (e *Event) Enabled() bool {
	return false
}

func (e *Event) Discard() *Event {
	return nil
}

func (e *Event) Msg(msg string) {
}

func (e *Event) Send() {
}

func (e *Event) Msgf(format string, v ...interface{}) {
}

func (e *Event) MsgFunc(createMsg func() string) {
}

func (e *Event) Fields(fields interface{}) *Event {
	return e
}

func (e *Event) Dict(key string, dict *Event) *Event {
	return e
}

func Dict() *Event {
	return nil
}

func (e *Event) Array(key string, arr LogArrayMarshaler) *Event {
	return e
}

func (e *Event) Object(key string, obj LogObjectMarshaler) *Event {
	return e
}

func (e *Event) Func(f func(e *Event)) *Event {
	return e
}

func (e *Event) EmbedObject(obj LogObjectMarshaler) *Event {
	return e
}

func (e *Event) Str(key, val string) *Event {
	return e
}

func (e *Event) Strs(key string, vals []string) *Event {
	return e
}

func (e *Event) Stringer(key string, val fmt.Stringer) *Event {
	return e
}

func (e *Event) Stringers(key string, vals []fmt.Stringer) *Event {
	return e
}

func (e *Event) Bytes(key string, val []byte) *Event {
	return e
}

func (e *Event) Hex(key string, val []byte) *Event {
	return e
}

func (e *Event) RawJSON(key string, b []byte) *Event {
	return e
}

func (e *Event) AnErr(key string, err error) *Event {
	return e
}

func (e *Event) Errs(key string, errs []error) *Event {
	return e
}

func (e *Event) Err(err error) *Event {
	return e
}

func (e *Event) Stack() *Event {
	return e
}

func (e *Event) Bool(key string, b bool) *Event {
	return e
}

func (e *Event) Bools(key string, b []bool) *Event {
	return e
}

func (e *Event) Int(key string, i int) *Event {
	return e
}

func (e *Event) Ints(key string, i []int) *Event {
	return e
}

func (e *Event) Int8(key string, i int8) *Event {
	return e
}

func (e *Event) Ints8(key string, i []int8) *Event {
	return e
}

func (e *Event) Int16(key string, i int16) *Event {
	return e
}

func (e *Event) Ints16(key string, i []int16) *Event {
	return e
}

func (e *Event) Int32(key string, i int32) *Event {
	return e
}

func (e *Event) Ints32(key string, i []int32) *Event {
	return e
}

func (e *Event) Int64(key string, i int64) *Event {
	return e
}

func (e *Event) Ints64(key string, i []int64) *Event {
	return e
}

func (e *Event) Uint(key string, i uint) *Event {
	return e
}

func (e *Event) Uints(key string, i []uint) *Event {
	return e
}

func (e *Event) Uint8(key string, i uint8) *Event {
	return e
}

func (e *Event) Uints8(key string, i []uint8) *Event {
	return e
}

func (e *Event) Uint16(key string, i uint16) *Event {
	return e
}

func (e *Event) Uints16(key string, i []uint16) *Event {
	return e
}

func (e *Event) Uint32(key string, i uint32) *Event {
	return e
}

func (e *Event) Uints32(key string, i []uint32) *Event {
	return e
}

func (e *Event) Uint64(key string, i uint64) *Event {
	return e
}

func (e *Event) Uints64(key string, i []uint64) *Event {
	return e
}

func (e *Event) Float32(key string, f float32) *Event {
	return e
}

func (e *Event) Floats32(key string, f []float32) *Event {
	return e
}

func (e *Event) Float64(key string, f float64) *Event {
	return e
}

func (e *Event) Floats64(key string, f []float64) *Event {
	return e
}

func (e *Event) Timestamp() *Event {
	return e
}

func (e *Event) Time(key string, t time.Time) *Event {
	return e
}

func (e *Event) Times(key string, t []time.Time) *Event {
	return e
}

func (e *Event) Dur(key string, d time.Duration) *Event {
	return e
}

func (e *Event) Durs(key string, d []time.Duration) *Event {
	return e
}

func (e *Event) TimeDiff(key string, t time.Time, start time.Time) *Event {
	return e
}

func (e *Event) Interface(key string, i interface{}) *Event {
	return e
}

func (e *Event) CallerSkipFrame(skip int) *Event {
	return e
}

func (e *Event) Caller(skip ...int) *Event {
	return e
}

func (e *Event) IPAddr(key string, ip net.IP) *Event {
	return e
}

func (e *Event) IPPrefix(key string, pfx net.IPNet) *Event {
	return e
}

func (e *Event) MACAddr(key string, ha net.HardwareAddr) *Event {
	return e
}

const (
	TimeFormatUnix = ""

	TimeFormatUnixMs = ""

	TimeFormatUnixMicro = ""

	TimeFormatUnixNano = ""
)

func SetGlobalLevel(l Level) {
}

func GlobalLevel() Level {
	return Disabled
}

func DisableSampling(v bool) {
}

type Hook interface {
	Run(e *Event, level Level, message string)
}

type HookFunc func(e *Event, level Level, message string)

func (h HookFunc) Run(e *Event, level Level, message string) {
}

type LevelHook struct {
	NoLevelHook, TraceHook, DebugHook, InfoHook, WarnHook, ErrorHook, FatalHook, PanicHook Hook
}

func (h LevelHook) Run(e *Event, level Level, message string) {
}

func NewLevelHook() LevelHook {
	return LevelHook{}
}

type Level int8

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

func (l Level) String() string {
	return ""
}

func ParseLevel(levelStr string) (Level, error) {
	return 0, nil
}

func (l *Level) UnmarshalText(text []byte) error {
	return nil
}

func (l Level) MarshalText() ([]byte, error) {
	return nil, nil
}

type Logger struct{}

func New(w io.Writer) Logger {
	return Logger{}
}

func Nop() Logger {
	return New(nil).Level(Disabled)
}

func (l Logger) Output(w io.Writer) Logger {
	return l
}

func (l Logger) With() Context {
	return Context{}
}

func (l *Logger) UpdateContext(update func(c Context) Context) {
}

func (l Logger) Level(lvl Level) Logger {
	return l
}

func (l Logger) GetLevel() Level {
	return Disabled
}

func (l Logger) Sample(s Sampler) Logger {
	return l
}

func (l Logger) Hook(h Hook) Logger {
	return l
}

func (l *Logger) Trace() *Event {
	return nil
}

func (l *Logger) Debug() *Event {
	return nil
}

func (l *Logger) Info() *Event {
	return nil
}

func (l *Logger) Warn() *Event {
	return nil
}

func (l *Logger) Error() *Event {
	return nil
}

func (l *Logger) Err(err error) *Event {
	return nil
}

func (l *Logger) Fatal() *Event {
	return nil
}

func (l *Logger) Panic() *Event {
	return nil
}

func (l *Logger) WithLevel(level Level) *Event {
	return nil
}

func (l *Logger) Log() *Event {
	return nil
}

func (l *Logger) Print(v ...interface{}) {
}

func (l *Logger) Printf(format string, v ...interface{}) {
}

func (l Logger) Write(p []byte) (n int, err error) {
	return
}

type Sampler interface {
	Sample(lvl Level) bool
}

type RandomSampler uint32

func (s RandomSampler) Sample(lvl Level) bool {
	return true
}

type BasicSampler struct {
	N uint32
}

func (s *BasicSampler) Sample(lvl Level) bool {
	return false
}

type BurstSampler struct {
	NextSampler Sampler
	Period      time.Duration
	Burst       uint32
}

func (s *BurstSampler) Sample(lvl Level) bool {
	return false
}

type LevelSampler struct {
	TraceSampler, DebugSampler, InfoSampler, WarnSampler, ErrorSampler Sampler
}

func (s LevelSampler) Sample(lvl Level) bool {
	return false
}

type LevelWriter interface {
	io.Writer
	WriteLevel(level Level, p []byte) (n int, err error)
}

func SyncWriter(w io.Writer) io.Writer {
	return nil
}

func MultiLevelWriter(writers ...io.Writer) LevelWriter {
	return nil
}

type TestingLog interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Helper()
}

type TestWriter struct {
	T TestingLog

	Frame int
}

func NewTestWriter(t TestingLog) TestWriter {
	return TestWriter{}
}

func (t TestWriter) Write(p []byte) (n int, err error) {
	return n, err
}

func ConsoleTestWriter(t TestingLog) func(w *ConsoleWriter) {
	return func(w *ConsoleWriter) {
	}
}
