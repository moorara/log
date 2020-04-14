// Package log provides leveled structured loggers.
// It hides the complexity of configuring and using the state-of-the-arts loggers by providing a single interface.
//
// You can either log using an instance logger or the singleton logger.
// After creating an instance logger, you need to set the singleton logger using the SetSingleton method once.
// The instance logger can be further used to create more contextualized loggers as the children of the root logger.
package log

import "strings"

// Format is the logging format.
type Format int

// Logging format
const (
	FormatJSON Format = iota
	FormatConsole
)

// Level is the logging level.
type Level int

// Logging level
const (
	LevelNone Level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

func parseLevel(level string) Level {
	switch strings.ToLower(level) {
	case "debug":
		return LevelDebug
	case "": // default
		fallthrough
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "none":
		fallthrough
	default:
		return LevelNone
	}
}

// Options are optional configurations for creating a logger.
type Options struct {
	Name        string
	Environment string
	Region      string
	Level       string
	Format      Format
}

// Logger is a leveled structured logger.
// It is concurrently safe to be used by multiple goroutines.
type Logger interface {
	With(kv ...interface{}) Logger
	GetLevel() Level
	SetLevel(level string)
	Debug(message string, kv ...interface{})
	Debugf(format string, args ...interface{})
	Info(message string, kv ...interface{})
	Infof(format string, args ...interface{})
	Warn(message string, kv ...interface{})
	Warnf(format string, args ...interface{})
	Error(message string, kv ...interface{})
	Errorf(format string, args ...interface{})
	Close() error
}

type nopLogger struct{}

// NewNopLogger creates a logger that never logs anything to anywhere!
// It can be used for testing purposes.
func NewNopLogger() Logger {
	return &nopLogger{}
}

func (l *nopLogger) With(kv ...interface{}) Logger             { return l }
func (l *nopLogger) GetLevel() Level                           { return LevelNone }
func (l *nopLogger) SetLevel(level string)                     {}
func (l *nopLogger) Debug(message string, kv ...interface{})   {}
func (l *nopLogger) Debugf(format string, args ...interface{}) {}
func (l *nopLogger) Info(message string, kv ...interface{})    {}
func (l *nopLogger) Infof(format string, args ...interface{})  {}
func (l *nopLogger) Warn(message string, kv ...interface{})    {}
func (l *nopLogger) Warnf(format string, args ...interface{})  {}
func (l *nopLogger) Error(message string, kv ...interface{})   {}
func (l *nopLogger) Errorf(format string, args ...interface{}) {}
func (l *nopLogger) Close() error                              { return nil }
