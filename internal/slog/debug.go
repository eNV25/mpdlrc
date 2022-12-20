//go:build debug

package slog

import (
	"context"
	"io"
	"time"

	"golang.org/x/exp/slog"
)

// Keys for "built-in" attributes.
const (
	// TimeKey is the key used by the built-in handlers for the time
	// when the log method is called. The associated Value is a [time.Time].
	TimeKey = slog.TimeKey
	// LevelKey is the key used by the built-in handlers for the level
	// of the log call. The associated value is a [Level].
	LevelKey = slog.LevelKey
	// MessageKey is the key used by the built-in handlers for the
	// message of the log call. The associated value is a string.
	MessageKey = slog.MessageKey
	// SourceKey is the key used by the built-in handlers for the source file
	// and line of the log call. The associated value is a string.
	SourceKey = slog.SourceKey
	// ErrorKey is the key used for errors by Logger.Error.
	// The associated value is an [error].
	ErrorKey = slog.ErrorKey
)

// Debug calls Logger.Debug on the default logger.
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// Error calls Logger.Error on the default logger.
func Error(msg string, err error, args ...any) {
	slog.Error(msg, err, args...)
}

// Info calls Logger.Info on the default logger.
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

// Log calls Logger.Log on the default logger.
func Log(level Level, msg string, args ...any) {
	slog.Log(level, msg, args...)
}

// LogAttrs calls Logger.LogAttrs on the default logger.
func LogAttrs(level Level, msg string, attrs ...Attr) {
	slog.LogAttrs(level, msg, attrs...)
}

// NewContext returns a context that contains the given Logger.
// Use FromContext to retrieve the Logger.
func NewContext(ctx context.Context, l *Logger) context.Context {
	return slog.NewContext(ctx, l)
}

// SetDefault makes l the default Logger.
// After this call, output from the log package's default Logger
// (as with [log.Print], etc.) will be logged at LevelInfo using l's Handler.
func SetDefault(l *Logger) {
	slog.SetDefault(l)
}

// Warn calls Logger.Warn on the default logger.
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// An Attr is a key-value pair.
type Attr = slog.Attr

// Any returns an Attr for the supplied value.
// See [Value.AnyValue] for how values are treated.
func Any(key string, value any) Attr {
	return slog.Any(key, value)
}

// Bool returns an Attr for a bool.
func Bool(key string, v bool) Attr {
	return slog.Bool(key, v)
}

// Duration returns an Attr for a time.Duration.
func Duration(key string, v time.Duration) Attr {
	return slog.Duration(key, v)
}

// Float64 returns an Attr for a floating-point number.
func Float64(key string, v float64) Attr {
	return slog.Float64(key, v)
}

// Group returns an Attr for a Group Value.
// The caller must not subsequently mutate the
// argument slice.
//
// Use Group to collect several Attrs under a single
// key on a log line, or as the result of LogValue
// in order to log a single value as multiple Attrs.
func Group(key string, as ...Attr) Attr {
	return slog.Group(key, as...)
}

// Int converts an int to an int64 and returns
// an Attr with that value.
func Int(key string, value int) Attr {
	return slog.Int(key, value)
}

// Int64 returns an Attr for an int64.
func Int64(key string, value int64) Attr {
	return slog.Int64(key, value)
}

// String returns an Attr for a string value.
func String(key, value string) Attr {
	return slog.String(key, value)
}

// Time returns an Attr for a time.Time.
// It discards the monotonic portion.
func Time(key string, v time.Time) Attr {
	return slog.Time(key, v)
}

// Uint64 returns an Attr for a uint64.
func Uint64(key string, v uint64) Attr {
	return slog.Uint64(key, v)
}

// A Handler handles log records produced by a Logger..
//
// A typical handler may print log records to standard error,
// or write them to a file or database, or perhaps augment them
// with additional attributes and pass them on to another handler.
//
// Any of the Handler's methods may be called concurrently with itself
// or with other methods. It is the responsibility of the Handler to
// manage this concurrency.
type Handler = slog.Handler

// HandlerOptions are options for a TextHandler or JSONHandler.
// A zero HandlerOptions consists entirely of default values.
type HandlerOptions = slog.HandlerOptions

// JSONHandler is a Handler that writes Records to an io.Writer as
// line-delimited JSON objects.
type JSONHandler = slog.JSONHandler

// NewJSONHandler creates a JSONHandler that writes to w,
// using the default options.
func NewJSONHandler(w io.Writer) *JSONHandler {
	return slog.NewJSONHandler(w)
}

// Kind is the kind of a Value.
type Kind = slog.Kind

const (
	AnyKind       = slog.AnyKind
	BoolKind      = slog.BoolKind
	DurationKind  = slog.DurationKind
	Float64Kind   = slog.Float64Kind
	Int64Kind     = slog.Int64Kind
	StringKind    = slog.StringKind
	TimeKind      = slog.TimeKind
	Uint64Kind    = slog.Uint64Kind
	GroupKind     = slog.GroupKind
	LogValuerKind = slog.LogValuerKind
)

// A Level is the importance or severity of a log event.
// The higher the level, the more important or severe the event.
type Level = slog.Level

// Second, we wanted to make it easy to use levels to specify logger verbosity.
// Since a larger level means a more severe event, a logger that accepts events
// with smaller (or more negative) level means a more verbose logger. Logger
// verbosity is thus the negation of event severity, and the default verbosity
// of 0 accepts all events at least as severe as INFO.
//
// Third, we wanted some room between levels to accommodate schemes with named
// levels between ours. For example, Google Cloud Logging defines a Notice level
// between Info and Warn. Since there are only a few of these intermediate
// levels, the gap between the numbers need not be large. Our gap of 4 matches
// OpenTelemetry's mapping. Subtracting 9 from an OpenTelemetry level in the
// DEBUG, INFO, WARN and ERROR ranges converts it to the corresponding slog
// Level range. OpenTelemetry also has the names TRACE and FATAL, which slog
// does not. But those OpenTelemetry levels can still be represented as slog
// Levels by using the appropriate integers.
//
// Names for common levels.
const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// A LevelVar is a Level variable, to allow a Handler level to change
// dynamically.
// It implements Leveler as well as a Set method,
// and it is safe for use by multiple goroutines.
// The zero LevelVar corresponds to LevelInfo.
type LevelVar = slog.LevelVar

// A Leveler provides a Level value.
//
// As Level itself implements Leveler, clients typically supply
// a Level value wherever a Leveler is needed, such as in HandlerOptions.
// Clients who need to vary the level dynamically can provide a more complex
// Leveler implementation such as *LevelVar.
type Leveler = slog.Leveler

// A LogValuer is any Go value that can convert itself into a Value for logging.
//
// This mechanism may be used to defer expensive operations until they are
// needed, or to expand a single value into a sequence of components.
type LogValuer = slog.LogValuer

// A Logger records structured information about each call to its
// Log, Debug, Info, Warn, and Error methods.
// For each call, it creates a Record and passes it to a Handler.
//
// To create a new Logger, call [New] or a Logger method
// that begins "With".
type Logger = slog.Logger

// Ctx retrieves a Logger from the given context using FromContext. Then it adds
// the given context to the Logger using WithContext and returns the result.
func Ctx(ctx context.Context) *Logger {
	return slog.FromContext(ctx)
}

// Default returns the default Logger.
func Default() *Logger { return slog.Default() }

// FromContext returns the Logger stored in ctx by NewContext, or the default
// Logger if there is none.
func FromContext(ctx context.Context) *Logger {
	return slog.FromContext(ctx)
}

// New creates a new Logger with the given Handler.
func New(h Handler) *Logger { return slog.New(h) }

// With calls Logger.With on the default logger.
func With(args ...any) *Logger {
	return slog.With(args...)
}

// A Record holds information about a log event.
// Copies of a Record share state.
// Do not modify a Record after handing out a copy to it.
// Use [Record.Clone] to create a copy with no shared state.
type Record = slog.Record

// NewRecord creates a Record from the given arguments.
// Use [Record.AddAttrs] to add attributes to the Record.
// If calldepth is greater than zero, [Record.SourceLine] will
// return the file and line number at that depth,
// where 1 means the caller of NewRecord.
//
// NewRecord is intended for logging APIs that want to support a [Handler] as
// a backend.
func NewRecord(t time.Time, level Level, msg string, calldepth int, ctx context.Context) Record {
	return slog.NewRecord(t, level, msg, calldepth, ctx)
}

// TextHandler is a Handler that writes Records to an io.Writer as a
// sequence of key=value pairs separated by spaces and followed by a newline.
type TextHandler = slog.TextHandler

// NewTextHandler creates a TextHandler that writes to w,
// using the default options.
func NewTextHandler(w io.Writer) *TextHandler {
	return slog.NewTextHandler(w)
}

// A Value can represent any Go value, but unlike type any,
// it can represent most small values without an allocation.
// The zero Value corresponds to nil.
type Value = slog.Value

// AnyValue returns a Value for the supplied value.
//
// If the supplied value is of type Value, it is returned
// unmodified.
//
// Given a value of one of Go's predeclared string, bool, or
// (non-complex) numeric types, AnyValue returns a Value of kind
// String, Bool, Uint64, Int64, or Float64. The width of the
// original numeric type is not preserved.
//
// Given a time.Time or time.Duration value, AnyValue returns a Value of kind
// TimeKind or DurationKind. The monotonic time is not preserved.
//
// For nil, or values of all other types, including named types whose
// underlying type is numeric, AnyValue returns a value of kind AnyKind.
func AnyValue(v any) Value {
	return slog.AnyValue(v)
}

// BoolValue returns a Value for a bool.
func BoolValue(v bool) Value {
	return slog.BoolValue(v)
}

// DurationValue returns a Value for a time.Duration.
func DurationValue(v time.Duration) Value {
	return slog.DurationValue(v)
}

// Float64Value returns a Value for a floating-point number.
func Float64Value(v float64) Value {
	return slog.Float64Value(v)
}

// GroupValue returns a new Value for a list of Attrs.
// The caller must not subsequently mutate the argument slice.
func GroupValue(as ...Attr) Value {
	return slog.GroupValue(as...)
}

// Int64Value returns a Value for an int64.
func Int64Value(v int64) Value {
	return slog.Int64Value(v)
}

// IntValue returns a Value for an int.
func IntValue(v int) Value {
	return slog.IntValue(v)
}

// String returns a new Value for a string.
func StringValue(value string) Value {
	return slog.StringValue(value)
}

// TimeValue returns a Value for a time.Time.
// It discards the monotonic portion.
func TimeValue(v time.Time) Value {
	return slog.TimeValue(v)
}

// Uint64Value returns a Value for a uint64.
func Uint64Value(v uint64) Value {
	return slog.Uint64Value(v)
}
