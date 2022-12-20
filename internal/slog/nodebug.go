//go:build !debug

package slog

import (
	"context"
	"io"
	"time"
)

const (
	TimeKey    = ""
	LevelKey   = ""
	MessageKey = ""
	SourceKey  = ""
	ErrorKey   = ""
)

func Debug(msg string, args ...any) {
}

func Error(msg string, err error, args ...any) {
}

func Info(msg string, args ...any) {
}

func Log(level Level, msg string, args ...any) {
}

func LogAttrs(level Level, msg string, attrs ...Attr) {
}

func NewContext(ctx context.Context, l *Logger) (_ context.Context) {
	return
}

func SetDefault(l *Logger) {
}

func Warn(msg string, args ...any) {
}

type Attr struct {
	Value Value
	Key   string
}

func Any(key string, value any) (_ Attr) {
	return
}

func Bool(key string, v bool) (_ Attr) {
	return
}

func Duration(key string, v time.Duration) (_ Attr) {
	return
}

func Float64(key string, v float64) (_ Attr) {
	return
}

func Group(key string, as ...Attr) (_ Attr) {
	return
}

func Int(key string, value int) (_ Attr) {
	return
}

func Int64(key string, value int64) (_ Attr) {
	return
}

func String(key, value string) (_ Attr) {
	return
}

func Time(key string, v time.Time) (_ Attr) {
	return
}

func Uint64(key string, v uint64) (_ Attr) {
	return
}

func (a Attr) Equal(b Attr) (_ bool) {
	return
}

func (a Attr) String() (_ string) {
	return
}

type Handler interface {
	Enabled(Level) bool
	Handle(r Record) error
	WithAttrs(attrs []Attr) Handler
	WithGroup(name string) Handler
}

type HandlerOptions struct {
	Level       Leveler
	ReplaceAttr func(groups []string, a Attr) Attr
	AddSource   bool
}

func (opts HandlerOptions) NewJSONHandler(w io.Writer) (_ *JSONHandler) {
	return
}

func (opts HandlerOptions) NewTextHandler(w io.Writer) (_ *TextHandler) {
	return
}

type JSONHandler struct{}

// NewJSONHandler creates a JSONHandler that writes to w,
// using the default options.
func NewJSONHandler(w io.Writer) (_ *JSONHandler) {
	return
}

func (h *JSONHandler) Enabled(level Level) (_ bool) {
	return
}

func (h *JSONHandler) Handle(r Record) (_ error) {
	return
}

func (h *JSONHandler) WithAttrs(attrs []Attr) (_ Handler) {
	return
}

func (h *JSONHandler) WithGroup(name string) (_ Handler) {
	return
}

type Kind int

const (
	AnyKind Kind = iota
	BoolKind
	DurationKind
	Float64Kind
	Int64Kind
	StringKind
	TimeKind
	Uint64Kind
	GroupKind
	LogValuerKind
)

func (k Kind) String() (_ string) {
	return
}

type Level int

const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

func (l Level) Level() (_ Level) { return }

func (l Level) MarshalJSON() (_ []byte, _ error) {
	return
}

func (l Level) String() (_ string) {
	return
}

type LevelVar struct{}

func (v *LevelVar) Level() (_ Level) {
	return
}

func (v *LevelVar) Set(l Level) {
}

func (v *LevelVar) String() (_ string) {
	return
}

type Leveler interface {
	Level() Level
}

type LogValuer interface {
	LogValue() Value
}

type Logger struct{}

func Ctx(ctx context.Context) (_ *Logger) {
	return
}

func Default() (_ *Logger) { return }

func FromContext(ctx context.Context) (_ *Logger) {
	return
}

func New(h Handler) (_ *Logger) { return }

func With(args ...any) (_ *Logger) {
	return
}

func (l *Logger) Context() (_ context.Context) { return }

func (l *Logger) Debug(msg string, args ...any) {
}

func (l *Logger) Enabled(level Level) (_ bool) {
	return
}

func (l *Logger) Error(msg string, err error, args ...any) {
}

func (l *Logger) Handler() (_ Handler) { return }

func (l *Logger) Info(msg string, args ...any) {
}

func (l *Logger) Log(level Level, msg string, args ...any) {
}

func (l *Logger) LogAttrs(level Level, msg string, attrs ...Attr) {
}

func (l *Logger) LogAttrsDepth(calldepth int, level Level, msg string, attrs ...Attr) {
}

func (l *Logger) LogDepth(calldepth int, level Level, msg string, args ...any) {
}

func (l *Logger) Warn(msg string, args ...any) {
}

func (l *Logger) With(args ...any) (_ *Logger) {
	return
}

func (l *Logger) WithContext(ctx context.Context) (_ *Logger) {
	return
}

func (l *Logger) WithGroup(name string) (_ *Logger) {
	return
}

type Record struct {
	Time    time.Time
	Context context.Context
	Message string
	Level   Level
}

func NewRecord(t time.Time, level Level, msg string, calldepth int, ctx context.Context) (_ Record) {
	return
}

func (r *Record) AddAttrs(attrs ...Attr) {
}

func (r Record) Attrs(f func(Attr)) {
}

func (r Record) Clone() (_ Record) {
	return
}

func (r Record) NumAttrs() (_ int) {
	return
}

func (r Record) SourceLine() (file string, line int) {
	return
}

type TextHandler struct{}

func NewTextHandler(w io.Writer) (_ *TextHandler) {
	return
}

func (h *TextHandler) Enabled(level Level) (_ bool) {
	return
}

func (h *TextHandler) Handle(r Record) (_ error) {
	return
}

func (h *TextHandler) WithAttrs(attrs []Attr) (_ Handler) {
	return
}

func (h *TextHandler) WithGroup(name string) (_ Handler) {
	return
}

type Value struct{}

func AnyValue(v any) (_ Value) {
	return
}

// BoolValue returns a Value for a bool.
func BoolValue(v bool) (_ Value) {
	return
}

func DurationValue(v time.Duration) (_ Value) {
	return
}

func Float64Value(v float64) (_ Value) {
	return
}

func GroupValue(as ...Attr) (_ Value) {
	return
}

func Int64Value(v int64) (_ Value) {
	return
}

func IntValue(v int) (_ Value) {
	return
}

func StringValue(value string) (_ Value) {
	return
}

func TimeValue(v time.Time) (_ Value) {
	return
}

func Uint64Value(v uint64) (_ Value) {
	return
}

func (v Value) Any() (_ any) {
	return
}

func (v Value) Bool() (_ bool) {
	return
}

func (v Value) Duration() (_ time.Duration) {
	return
}

func (v Value) Equal(w Value) (_ bool) {
	return
}

func (v Value) Float64() (_ float64) {
	return
}

func (v Value) Group() (_ []Attr) {
	return
}

func (v Value) Int64() (_ int64) {
	return
}

func (v Value) Kind() (_ Kind) {
	return
}

func (v Value) LogValuer() (_ LogValuer) {
	return
}

func (v Value) Resolve() (_ Value) {
	return
}

func (v Value) String() (_ string) {
	return
}

func (v Value) Time() (_ time.Time) {
	return
}

func (v Value) Uint64() (_ uint64) {
	return
}
