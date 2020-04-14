package log

import (
	"strings"

	zaplog "go.uber.org/zap"
	zapcore "go.uber.org/zap/zapcore"
)

const (
	instanceCallerSkip  = 1
	singletonCallerSkip = 2
)

// zapLogger is an interface for zap.Logger struct.
type zapLogger interface {
	Sugar() *zaplog.SugaredLogger
}

// zapSugaredLogger is an interface for zap.SugaredLogger struct.
type zapSugaredLogger interface {
	Sync() error
	Desugar() *zaplog.Logger
	With(...interface{}) *zaplog.SugaredLogger
	Debugw(string, ...interface{})
	Debugf(string, ...interface{})
	Infow(string, ...interface{})
	Infof(string, ...interface{})
	Warnw(string, ...interface{})
	Warnf(string, ...interface{})
	Errorw(string, ...interface{})
	Errorf(string, ...interface{})
}

// zap is an implementation of Logger using zap.
type zap struct {
	config        *zaplog.Config
	logger        zapLogger
	sugaredLogger zapSugaredLogger
}

// NewZap creates a new logger based on zap logger.
func NewZap(opts Options) Logger {
	config := zaplog.NewProductionConfig()
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.NameKey = "logger"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	config.OutputPaths = []string{"stdout"}
	config.InitialFields = make(map[string]interface{})

	if opts.Name != "" {
		config.InitialFields["logger"] = opts.Name
	}

	if opts.Environment != "" {
		config.InitialFields["environment"] = opts.Environment
	}

	if opts.Region != "" {
		config.InitialFields["region"] = opts.Region
	}

	switch strings.ToLower(opts.Level) {
	case "debug":
		config.Level = zaplog.NewAtomicLevelAt(zapcore.DebugLevel)
	case "": // default
		fallthrough
	case "info":
		config.Level = zaplog.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zaplog.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zaplog.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "none":
		fallthrough
	default:
		config.Level = zaplog.NewAtomicLevelAt(zapcore.Level(99))
	}

	switch opts.Format {
	case FormatJSON:
		config.Encoding = "json"
	case FormatConsole:
		config.Encoding = "console"
	}

	logger, _ := config.Build(
		zaplog.AddCaller(),
		zaplog.AddCallerSkip(instanceCallerSkip),
	)

	return &zap{
		config:        &config,
		logger:        logger,
		sugaredLogger: logger.Sugar(),
	}
}

// With returns a new logger that automatically logs the given set of key-value pairs.
// This can be used for creating a contextualized logger.
func (z *zap) With(kv ...interface{}) Logger {
	sugaredLogger := z.sugaredLogger.With(kv...)

	return &zap{
		config:        z.config,
		logger:        sugaredLogger.Desugar(),
		sugaredLogger: sugaredLogger,
	}
}

// GetLevel returns the current logging level.
func (z *zap) GetLevel() Level {
	switch z.config.Level.Level() {
	case zapcore.DebugLevel:
		return LevelDebug
	case zapcore.InfoLevel:
		return LevelInfo
	case zapcore.WarnLevel:
		return LevelWarn
	case zapcore.ErrorLevel:
		return LevelError
	default:
		return LevelNone
	}
}

// SetLevel changes the logging level.
func (z *zap) SetLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		z.config.Level.SetLevel(zapcore.DebugLevel)
	case "info":
		z.config.Level.SetLevel(zapcore.InfoLevel)
	case "warn":
		z.config.Level.SetLevel(zapcore.WarnLevel)
	case "error":
		z.config.Level.SetLevel(zapcore.ErrorLevel)
	}
}

// Debug logs a message and a list of key-value pairs in debug level.
func (z *zap) Debug(message string, kv ...interface{}) {
	z.sugaredLogger.Debugw(message, kv...)
}

// Debugf formats and logs a message in debug level.
// It uses fmt.Sprintf() to log a message.
func (z *zap) Debugf(format string, args ...interface{}) {
	z.sugaredLogger.Debugf(format, args...)
}

// Info logs a message and a list of key-value pairs in info level.
func (z *zap) Info(message string, kv ...interface{}) {
	z.sugaredLogger.Infow(message, kv...)
}

// Infof formats and logs a message in info level.
// It uses fmt.Sprintf() to log a message.
func (z *zap) Infof(format string, args ...interface{}) {
	z.sugaredLogger.Infof(format, args...)
}

// Warn logs a message and a list of key-value pairs in warn level.
func (z *zap) Warn(message string, kv ...interface{}) {
	z.sugaredLogger.Warnw(message, kv...)
}

// Warnf formats and logs a message in warn level.
// It uses fmt.Sprintf() to log a message.
func (z *zap) Warnf(format string, args ...interface{}) {
	z.sugaredLogger.Warnf(format, args...)
}

// Error logs a message and a list of key-value pairs in error level.
func (z *zap) Error(message string, kv ...interface{}) {
	z.sugaredLogger.Errorw(message, kv...)
}

// Errorf formats and logs a message in error level.
// It uses fmt.Sprintf() to log a message.
func (z *zap) Errorf(format string, args ...interface{}) {
	z.sugaredLogger.Errorf(format, args...)
}

// Close flushes the logger.
func (z *zap) Close() error {
	return z.sugaredLogger.Sync()
}
