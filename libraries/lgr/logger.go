package lgr

import (
	"time"
)

// Level indicates the log priority. Higher levels are more important and require more immediate attention.
type Level int8

const (
	// DebugLevel should be used for information that may be needed when diagnosing and troubleshooting issues
	// or when running the application in a test environment to make sure everything is running correctly.
	DebugLevel Level = iota - 1 // Start at -1 so the default value (0) is InfoLevel.
	// InfoLevel is the default log level and should be used to indicate that something happened. Info logged
	// at this level should be purely informative and should not require monitoring or attention.
	InfoLevel
	// WarnLevel should be used to indicate that something unexpected has happened in the application but
	// the application can continue to function and deal with the issue internally. Logs reported at this level
	// should be monitored but may not require urgent attention.
	WarnLevel
	// ErrorLevel should be used to indicate that something unexpected has happened in the application that
	// is preventing one or more functionalities from functioning properly. Logs reported at this level should
	// be monitored and addressed with urgency.
	ErrorLevel
)

// Adapter is an interface to the underlying logger/log sink so that we can be vendor agnostic
// moving forward.
type Adapter interface {
	// Adapt should immediately write a log to the underlying logger.
	Adapt(level Level, message string, fields ...Field)
}

// TimestampFactoryFunc represents a function knows how to create time values
// that will be used by the logger to set the timestamp field in the log context.
type TimestampFactoryFunc func() time.Time

// Logger represents an active logging object that wraps an underlying logger (such as Uber's zap)
// via an Adapter. Each logging operation makes a single call Adapter.Adapt where the log is written
// to the underlying, wrapped logger, immediately.
type Logger struct {
	adapter          Adapter
	minLevel         Level
	outputPath       string
	timestampFactory TimestampFactoryFunc
}

// Option allows a user to configure the logger without exposing the internals of the Logger
// in the public API.
type Option func(l *Logger)

// WithMinLevel will set the minimum level that logs will be written at. This is useful for disabling
// DebugLevel logs when in production by setting WithMinLevel(InfoLevel) etc.
func WithMinLevel(level Level) Option {
	return func(l *Logger) {
		l.minLevel = level
	}
}

// WithOutputPath allows setting the path that logs will be written to. This would typically
// be set to stdout or stderr and the default is set to stderr.
func WithOutputPath(outputPath string) Option {
	return func(l *Logger) {
		l.outputPath = outputPath
	}
}

// WithTimestampFactory allows setting the TimestampFactoryFunc on the logger that will be used
// to determine the timestamps of the logs. This is useful when doing test automation so that
// the timestamp is not constantly changing when the log is created.
func WithTimestampFactory(factory TimestampFactoryFunc) Option {
	return func(l *Logger) {
		l.timestampFactory = factory
	}
}

// New creates a new Logger instance and sets any optional configuration before returning the
// Logger.
func New(opts ...Option) (*Logger, error) {
	logger := FromAdapter(nil, opts...)
	bindZerologAdapter(logger)

	return logger, nil
}

// NewNop creates a new no-op logger that can be used to disable logs in an application without having
// to remove all log instances. It will return a nil *Logger and each of the level methods will return
// before trying to call the Adapter when called.
func NewNop() *Logger {
	return nil
}

// FromAdapter will create a new logger with the given adapter and the passed options applied.
// Passing nil as the adapter will have the same effect as NewNop.
func FromAdapter(adapter Adapter, opts ...Option) *Logger {
	logger := &Logger{
		outputPath:       "stderr",
		minLevel:         DebugLevel,
		timestampFactory: time.Now,
	}

	for _, opt := range opts {
		opt(logger)
	}

	logger.adapter = adapter

	return logger
}

// Debug will write a log at DebugLevel with the given msg and fields as context. See the Level constants
// for information on when the level should be used.
func (l *Logger) Debug(msg string, fields ...Field) {
	if l == nil || l.adapter == nil {
		return
	}

	l.adapter.Adapt(DebugLevel, msg, fields...)
}

// Info will write a log at InfoLevel with the given msg and fields as context. See the Level constants
// for information on when the level should be used.
func (l *Logger) Info(msg string, fields ...Field) {
	if l == nil || l.adapter == nil {
		return
	}

	l.adapter.Adapt(InfoLevel, msg, fields...)
}

// Warn will write a log at WarnLevel with the given msg and fields as context. See the Level constants
// for information on when the level should be used.
func (l *Logger) Warn(msg string, fields ...Field) {
	if l == nil || l.adapter == nil {
		return
	}

	l.adapter.Adapt(WarnLevel, msg, fields...)
}

// Error will write a log at ErrorLevel with the given msg and fields as context. See the Level constants
// for information on when the level should be used.
func (l *Logger) Error(msg string, fields ...Field) {
	if l == nil || l.adapter == nil {
		return
	}

	l.adapter.Adapt(ErrorLevel, msg, fields...)
}
