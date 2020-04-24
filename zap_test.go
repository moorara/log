package log

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	zaplog "go.uber.org/zap"
	zapcore "go.uber.org/zap/zapcore"
)

// mockZapLogger is a mock implementation of zapLogger.
type mockZapLogger struct {
	SugarOutSugaredLogger *zaplog.SugaredLogger
}

func (m *mockZapLogger) Sugar() *zaplog.SugaredLogger {
	return m.SugarOutSugaredLogger
}

// mockZapSugaredLogger is a mock implementation of zapSugaredLogger.
type mockZapSugaredLogger struct {
	SyncOutError         error
	DesugarOutLogger     *zaplog.Logger
	WithInArgs           []interface{}
	WithOutSugaredLogger *zaplog.SugaredLogger
	DebugwInMsg          string
	DebugwInKV           []interface{}
	DebugfInTemplate     string
	DebugfInArgs         []interface{}
	InfowInMsg           string
	InfowInKV            []interface{}
	InfofInTemplate      string
	InfofInArgs          []interface{}
	WarnwInMsg           string
	WarnwInKV            []interface{}
	WarnfInTemplate      string
	WarnfInArgs          []interface{}
	ErrorwInMsg          string
	ErrorwInKV           []interface{}
	ErrorfInTemplate     string
	ErrorfInArgs         []interface{}
}

func (m *mockZapSugaredLogger) Sync() error {
	return m.SyncOutError
}

func (m *mockZapSugaredLogger) Desugar() *zaplog.Logger {
	return m.DesugarOutLogger
}

func (m *mockZapSugaredLogger) With(args ...interface{}) *zaplog.SugaredLogger {
	m.WithInArgs = args
	return m.WithOutSugaredLogger
}

func (m *mockZapSugaredLogger) Debugw(msg string, kv ...interface{}) {
	m.DebugwInMsg, m.DebugwInKV = msg, kv
}

func (m *mockZapSugaredLogger) Debugf(template string, args ...interface{}) {
	m.DebugfInTemplate, m.DebugfInArgs = template, args
}

func (m *mockZapSugaredLogger) Infow(msg string, kv ...interface{}) {
	m.InfowInMsg, m.InfowInKV = msg, kv
}

func (m *mockZapSugaredLogger) Infof(template string, args ...interface{}) {
	m.InfofInTemplate, m.InfofInArgs = template, args
}

func (m *mockZapSugaredLogger) Warnw(msg string, kv ...interface{}) {
	m.WarnwInMsg, m.WarnwInKV = msg, kv
}

func (m *mockZapSugaredLogger) Warnf(template string, args ...interface{}) {
	m.WarnfInTemplate, m.WarnfInArgs = template, args
}

func (m *mockZapSugaredLogger) Errorw(msg string, kv ...interface{}) {
	m.ErrorwInMsg, m.ErrorwInKV = msg, kv
}

func (m *mockZapSugaredLogger) Errorf(template string, args ...interface{}) {
	m.ErrorfInTemplate, m.ErrorfInArgs = template, args
}

func TestNewZap(t *testing.T) {
	tests := []struct {
		name string
		opts Options
	}{
		{
			"Default",
			Options{},
		},
		{
			"Production",
			Options{
				Name:        "my-service",
				Version:     "0.1.0",
				Environment: "production",
				Region:      "ca-central-1",
				Tags: map[string]string{
					"domain": "auth",
				},
			},
		},
		{
			"LevelNone",
			Options{
				Name:  "test",
				Level: "none",
			},
		},
		{
			"LevelError",
			Options{
				Name:  "test",
				Level: "error",
			},
		},
		{
			"LevelWarn",
			Options{
				Name:  "test",
				Level: "warn",
			},
		},
		{
			"LevelInfo",
			Options{
				Name:  "test",
				Level: "info",
			},
		},
		{
			"LevelDebug",
			Options{
				Name:  "test",
				Level: "debug",
			},
		},
		{
			"JSONLogger",
			Options{
				Format: FormatJSON,
			},
		},
		{
			"ConsoleLogger",
			Options{
				Format: FormatConsole,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(T *testing.T) {
			logger := NewZap(tc.opts)

			assert.NotNil(t, logger)
			assert.IsType(t, &zap{}, logger)
		})
	}
}

func TestZapWith(t *testing.T) {
	zlogger := zaplog.NewNop()

	tests := []struct {
		name                 string
		mockZapSugaredLogger *mockZapSugaredLogger
		kv                   []interface{}
		expectedLogger       Logger
	}{
		{
			name: "OK",
			mockZapSugaredLogger: &mockZapSugaredLogger{
				DesugarOutLogger:     zlogger,
				WithOutSugaredLogger: zlogger.Sugar(),
			},
			kv: []interface{}{
				"version", "0.1.0",
				"revision", "1234567",
				"context", "test",
			},
			expectedLogger: &zap{
				config:        nil,
				logger:        zlogger,
				sugaredLogger: zlogger.Sugar(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			zl := &zap{
				sugaredLogger: tc.mockZapSugaredLogger,
			}

			logger := zl.With(tc.kv...)

			assert.Equal(t, tc.expectedLogger, logger)
		})
	}
}

func TestZapGetLevel(t *testing.T) {
	tests := []struct {
		name          string
		config        *zaplog.Config
		expectedLevel Level
	}{
		{
			name: "Debug",
			config: &zaplog.Config{
				Level: zaplog.NewAtomicLevelAt(zapcore.DebugLevel),
			},
			expectedLevel: LevelDebug,
		},
		{
			name: "Info",
			config: &zaplog.Config{
				Level: zaplog.NewAtomicLevelAt(zapcore.InfoLevel),
			},
			expectedLevel: LevelInfo,
		},
		{
			name: "Warn",
			config: &zaplog.Config{
				Level: zaplog.NewAtomicLevelAt(zapcore.WarnLevel),
			},
			expectedLevel: LevelWarn,
		},
		{
			name: "Error",
			config: &zaplog.Config{
				Level: zaplog.NewAtomicLevelAt(zapcore.ErrorLevel),
			},
			expectedLevel: LevelError,
		},
		{
			name: "Unsupported",
			config: &zaplog.Config{
				Level: zaplog.NewAtomicLevelAt(zapcore.FatalLevel),
			},
			expectedLevel: LevelNone,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			zl := &zap{
				config: tc.config,
			}

			level := zl.GetLevel()

			assert.Equal(t, tc.expectedLevel, level)
		})
	}
}

func TestZapSetLevel(t *testing.T) {
	tests := []struct {
		name          string
		config        *zaplog.Config
		level         string
		expectedLevel zapcore.Level
	}{
		{
			name: "Debug",
			config: &zaplog.Config{
				Level: zaplog.NewAtomicLevel(),
			},
			level:         "debug",
			expectedLevel: zapcore.DebugLevel,
		},
		{
			name: "Info",
			config: &zaplog.Config{
				Level: zaplog.NewAtomicLevel(),
			},
			level:         "info",
			expectedLevel: zapcore.InfoLevel,
		},
		{
			name: "Warn",
			config: &zaplog.Config{
				Level: zaplog.NewAtomicLevel(),
			},
			level:         "warn",
			expectedLevel: zapcore.WarnLevel,
		},
		{
			name: "Error",
			config: &zaplog.Config{
				Level: zaplog.NewAtomicLevel(),
			},
			level:         "error",
			expectedLevel: zapcore.ErrorLevel,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			zl := &zap{
				config: tc.config,
			}

			zl.SetLevel(tc.level)

			level := tc.config.Level.Level()
			assert.Equal(t, tc.expectedLevel, level)
		})
	}
}

func TestZapLog(t *testing.T) {
	tests := []struct {
		name                 string
		mockZapSugaredLogger *mockZapSugaredLogger
		message              string
		kv                   []interface{}
	}{
		{
			name:                 "OK",
			mockZapSugaredLogger: &mockZapSugaredLogger{},
			message:              "operation succeeded",
			kv:                   []interface{}{"operation", "test"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			zl := &zap{
				sugaredLogger: tc.mockZapSugaredLogger,
			}

			t.Run("Debug", func(t *testing.T) {
				zl.Debug(tc.message, tc.kv...)

				assert.Equal(t, tc.message, tc.mockZapSugaredLogger.DebugwInMsg)
				for _, val := range tc.kv {
					assert.Contains(t, tc.mockZapSugaredLogger.DebugwInKV, val)
				}
			})

			t.Run("Info", func(t *testing.T) {
				zl.Info(tc.message, tc.kv...)

				assert.Equal(t, tc.message, tc.mockZapSugaredLogger.InfowInMsg)
				for _, val := range tc.kv {
					assert.Contains(t, tc.mockZapSugaredLogger.InfowInKV, val)
				}
			})

			t.Run("Warn", func(t *testing.T) {
				zl.Warn(tc.message, tc.kv...)

				assert.Equal(t, tc.message, tc.mockZapSugaredLogger.WarnwInMsg)
				for _, val := range tc.kv {
					assert.Contains(t, tc.mockZapSugaredLogger.WarnwInKV, val)
				}
			})

			t.Run("Error", func(t *testing.T) {
				zl.Error(tc.message, tc.kv...)

				assert.Equal(t, tc.message, tc.mockZapSugaredLogger.ErrorwInMsg)
				for _, val := range tc.kv {
					assert.Contains(t, tc.mockZapSugaredLogger.ErrorwInKV, val)
				}
			})
		})
	}
}

func TestZapLogf(t *testing.T) {
	tests := []struct {
		name                 string
		mockZapSugaredLogger *mockZapSugaredLogger
		format               string
		args                 []interface{}
	}{
		{
			name:                 "OK",
			mockZapSugaredLogger: &mockZapSugaredLogger{},
			format:               "operation succeeded: %s",
			args:                 []interface{}{"test"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			zl := &zap{
				sugaredLogger: tc.mockZapSugaredLogger,
			}

			t.Run("Debugf", func(t *testing.T) {
				zl.Debugf(tc.format, tc.args...)

				assert.Equal(t, tc.format, tc.mockZapSugaredLogger.DebugfInTemplate)
				for _, val := range tc.args {
					assert.Contains(t, tc.mockZapSugaredLogger.DebugfInArgs, val)
				}
			})

			t.Run("Infof", func(t *testing.T) {
				zl.Infof(tc.format, tc.args...)

				assert.Equal(t, tc.format, tc.mockZapSugaredLogger.InfofInTemplate)
				for _, val := range tc.args {
					assert.Contains(t, tc.mockZapSugaredLogger.InfofInArgs, val)
				}
			})

			t.Run("Warnf", func(t *testing.T) {
				zl.Warnf(tc.format, tc.args...)

				assert.Equal(t, tc.format, tc.mockZapSugaredLogger.WarnfInTemplate)
				for _, val := range tc.args {
					assert.Contains(t, tc.mockZapSugaredLogger.WarnfInArgs, val)
				}
			})

			t.Run("Errorf", func(t *testing.T) {
				zl.Errorf(tc.format, tc.args...)

				assert.Equal(t, tc.format, tc.mockZapSugaredLogger.ErrorfInTemplate)
				for _, val := range tc.args {
					assert.Contains(t, tc.mockZapSugaredLogger.ErrorfInArgs, val)
				}
			})
		})
	}
}

func TestZapClose(t *testing.T) {
	tests := []struct {
		name                 string
		mockZapSugaredLogger *mockZapSugaredLogger
		expectedError        error
	}{
		{
			name:                 "NoError",
			mockZapSugaredLogger: &mockZapSugaredLogger{},
			expectedError:        nil,
		},
		{
			name: "WithError",
			mockZapSugaredLogger: &mockZapSugaredLogger{
				SyncOutError: errors.New("sync error"),
			},
			expectedError: errors.New("sync error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			zl := &zap{
				sugaredLogger: tc.mockZapSugaredLogger,
			}

			err := zl.Close()

			assert.Equal(t, tc.expectedError, err)
		})
	}
}
