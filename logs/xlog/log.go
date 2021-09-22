package xlog

import (
	"errors"
	"io"
	"os"
	"strings"
	"sync/atomic"
)

const (
	FmtEmptySeparate = ""
)

// log level
type Level uint32

// const log level
const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = iota
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)

var errUnmarshalNilLevel = errors.New("can't unmarshal a nil *Level")

func String(l Level) string {
	switch l {
	default:
		return "UNKNOWN"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	}
}

func (l *Level) UnmarshalText(text []byte) bool {
	switch strings.ToUpper(string(text)) {
	case "DEBUG":
		*l = DebugLevel
	case "INFO": // make the zero value useful
		*l = InfoLevel
	case "WARN":
		*l = WarnLevel
	case "ERROR":
		*l = ErrorLevel
	case "PANIC":
		*l = PanicLevel
	case "FATAL":
		*l = FatalLevel
	default:
		return false
	}
	return true
}

type logoptions struct {
	output        io.Writer
	level         Level
	stdLevel      Level
	formatter     Formatter
	disableCaller bool
}

type LogOption func(*logoptions)

func initOptions(opts ...LogOption) (o *logoptions) {
	o = &logoptions{}
	for _, opt := range opts {
		opt(o)
	}
	//The default value is stderr
	if o.output == nil {
		o.output = os.Stderr
	}
	//The default value is textFormatter
	if o.formatter == nil {
		o.formatter = &TextFormatter{}
	}

	return
}

func WithOutPut(p io.Writer) LogOption {
	return func(o *logoptions) {
		o.output = p
	}
}

func WithLevel(level Level) LogOption {
	return func(o *logoptions) {
		atomic.StoreUint32((*uint32)(&o.level), uint32(level))
	}
}

func WithStdLevel(level Level) LogOption {
	return func(o *logoptions) {
		atomic.StoreUint32((*uint32)(&o.stdLevel), uint32(level))
	}
}

func WithFormatter(formatter Formatter) LogOption {
	return func(o *logoptions) {
		o.formatter = formatter
	}
}

func WithDisableCaller(caller bool) LogOption {
	return func(o *logoptions) {
		o.disableCaller = caller
	}
}
