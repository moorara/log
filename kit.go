package log

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const (
	instanceCallerDepth  = 7
	singletonCallerDepth = 8
)

// kit is an implementation of Logger using go-kit
type kit struct {
	level  Level
	base   log.Logger
	logger *log.SwapLogger
}

func createBaseLogger(opts Options) log.Logger {
	var base log.Logger

	switch opts.Format {
	case FormatConsole:
		base = log.NewLogfmtLogger(os.Stdout)
	case FormatJSON:
		fallthrough
	default:
		base = log.NewJSONLogger(os.Stdout)
	}

	// This is not required since SwapLogger uses a SyncLogger and can be used concurrently
	// base = log.NewSyncLogger(base)

	context := []interface{}{
		"timestamp", log.DefaultTimestamp,
		"caller", log.Caller(instanceCallerDepth),
	}

	if opts.Name != "" {
		context = append(context, "logger", opts.Name)
	}

	if opts.Environment != "" {
		context = append(context, "environment", opts.Environment)
	}

	if opts.Region != "" {
		context = append(context, "region", opts.Region)
	}

	base = log.With(base, context...)

	return base
}

func createFilteredLogger(base log.Logger, l Level) log.Logger {
	switch l {
	case LevelDebug:
		return level.NewFilter(base, level.AllowDebug())
	case LevelInfo:
		return level.NewFilter(base, level.AllowInfo())
	case LevelWarn:
		return level.NewFilter(base, level.AllowWarn())
	case LevelError:
		return level.NewFilter(base, level.AllowError())
	case LevelNone:
		return level.NewFilter(base, level.AllowNone())
	default:
		return base
	}
}

// NewKit creates a new logger based on go-kit logger.
func NewKit(opts Options) Logger {
	level := parseLevel(opts.Level)
	base := createBaseLogger(opts)
	logger := new(log.SwapLogger)

	filtered := createFilteredLogger(base, level)
	logger.Swap(filtered)

	return &kit{
		level:  level,
		base:   base,
		logger: logger,
	}
}

// With returns a new logger that automatically logs the given set of key-value pairs.
// This can be used for creating a contextualized logger.
func (k *kit) With(kv ...interface{}) Logger {
	level := k.level
	base := log.With(k.base, kv...)
	logger := new(log.SwapLogger)

	filtered := createFilteredLogger(base, level)
	logger.Swap(filtered)

	return &kit{
		level:  level,
		base:   base,
		logger: logger,
	}
}

// GetLevel returns the current logging level.
func (k *kit) GetLevel() Level {
	return k.level
}

// SetLevel changes the logging level.
func (k *kit) SetLevel(level string) {
	k.level = parseLevel(level)
	filtered := createFilteredLogger(k.base, k.level)
	k.logger.Swap(filtered)
}

// Debug logs a message and a list of key-value pairs in debug level.
func (k *kit) Debug(message string, kv ...interface{}) {
	kv = append(kv, "message", message)
	_ = level.Debug(k.logger).Log(kv...)
}

// Debugf formats and logs a message in debug level.
// It uses fmt.Sprintf() to log a message.
func (k *kit) Debugf(format string, v ...interface{}) {
	_ = level.Debug(k.logger).Log("message", fmt.Sprintf(format, v...))
}

// Info logs a message and a list of key-value pairs in info level.
func (k *kit) Info(message string, kv ...interface{}) {
	kv = append(kv, "message", message)
	_ = level.Info(k.logger).Log(kv...)
}

// Infof formats and logs a message in info level.
// It uses fmt.Sprintf() to log a message.
func (k *kit) Infof(format string, v ...interface{}) {
	_ = level.Info(k.logger).Log("message", fmt.Sprintf(format, v...))
}

// Warn logs a message and a list of key-value pairs in warn level.
func (k *kit) Warn(message string, kv ...interface{}) {
	kv = append(kv, "message", message)
	_ = level.Warn(k.logger).Log(kv...)
}

// Warnf formats and logs a message in warn level.
// It uses fmt.Sprintf() to log a message.
func (k *kit) Warnf(format string, v ...interface{}) {
	_ = level.Warn(k.logger).Log("message", fmt.Sprintf(format, v...))
}

// Error logs a message and a list of key-value pairs in error level.
func (k *kit) Error(message string, kv ...interface{}) {
	kv = append(kv, "message", message)
	_ = level.Error(k.logger).Log(kv...)
}

// Errorf formats and logs a message in error level.
// It uses fmt.Sprintf() to log a message.
func (k *kit) Errorf(format string, v ...interface{}) {
	_ = level.Error(k.logger).Log("message", fmt.Sprintf(format, v...))
}

// Close flushes the logger.
func (k *kit) Close() error {
	return nil
}
