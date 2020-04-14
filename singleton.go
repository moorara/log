package log

import (
	kitlog "github.com/go-kit/kit/log"
	zaplog "go.uber.org/zap"
)

// The singleton logger
var singleton Logger

// SetSingleton updates the singleton logger.
func SetSingleton(l Logger) {
	switch v := l.(type) {
	case *kit:
		base := kitlog.With(v.base, "caller", kitlog.Caller(singletonCallerDepth))
		logger := new(kitlog.SwapLogger)
		filtered := createFilteredLogger(base, v.level)
		logger.Swap(filtered)

		singleton = &kit{
			level:  v.level,
			base:   base,
			logger: logger,
		}

	case *zap:
		logger, _ := v.config.Build(
			zaplog.AddCaller(),
			zaplog.AddCallerSkip(singletonCallerSkip),
		)

		singleton = &zap{
			config:        v.config,
			logger:        logger,
			sugaredLogger: logger.Sugar(),
		}

	default:
		singleton = l
	}
}

// GetLevel returns the current logging level of the singleton logger.
func GetLevel() Level {
	if singleton != nil {
		return singleton.GetLevel()
	}
	return LevelNone
}

// SetLevel changes the logging level of the singleton logger.
func SetLevel(level string) {
	if singleton != nil {
		singleton.SetLevel(level)
	}
}

// Debug logs a message and a list of key-value pairs in debug level using the singleton logger.
func Debug(message string, kv ...interface{}) {
	if singleton != nil {
		singleton.Debug(message, kv...)
	}
}

// Debugf formats and logs a message in debug level using the singleton logger.
// It uses fmt.Sprintf() to log a message.
func Debugf(format string, v ...interface{}) {
	if singleton != nil {
		singleton.Debugf(format, v...)
	}
}

// Info logs a message and a list of key-value pairs in info level using the singleton logger.
func Info(message string, kv ...interface{}) {
	if singleton != nil {
		singleton.Info(message, kv...)
	}
}

// Infof formats and logs a message in info level using the singleton logger.
// It uses fmt.Sprintf() to log a message.
func Infof(format string, v ...interface{}) {
	if singleton != nil {
		singleton.Infof(format, v...)
	}
}

// Warn logs a message and a list of key-value pairs in warn level using the singleton logger.
func Warn(message string, kv ...interface{}) {
	if singleton != nil {
		singleton.Warn(message, kv...)
	}
}

// Warnf formats and logs a message in warn level using the singleton logger.
// It uses fmt.Sprintf() to log a message.
func Warnf(format string, v ...interface{}) {
	if singleton != nil {
		singleton.Warnf(format, v...)
	}
}

// Error logs a message and a list of key-value pairs in error level using the singleton logger.
func Error(message string, kv ...interface{}) {
	if singleton != nil {
		singleton.Error(message, kv...)
	}
}

// Errorf formats and logs a message in error level using the singleton logger.
// It uses fmt.Sprintf() to log a message.
func Errorf(format string, v ...interface{}) {
	if singleton != nil {
		singleton.Errorf(format, v...)
	}
}

// Close flushes the singleton logger.
func Close() error {
	if singleton != nil {
		return singleton.Close()
	}
	return nil
}
