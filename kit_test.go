package log

import (
	"errors"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

// mockKitLogger is a mock implementation of go-kit Logger
type mockKitLogger struct {
	LogInKV     []interface{}
	LogOutError error
}

func (m *mockKitLogger) Log(kv ...interface{}) error {
	m.LogInKV = kv
	return m.LogOutError
}

func TestCreateBaseLogger(t *testing.T) {
	tests := []struct {
		name string
		opts Options
	}{
		{
			"NoOption",
			Options{},
		},
		{
			"WithName",
			Options{
				Name: "test",
			},
		},
		{
			"WithMetadata",
			Options{
				Name:        "test",
				Environment: "local",
				Region:      "local",
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
		t.Run(tc.name, func(t *testing.T) {
			base := createBaseLogger(tc.opts)

			assert.NotNil(t, base)
		})
	}
}

func TestCreateFilteredLogger(t *testing.T) {
	tests := []struct {
		name  string
		base  kitlog.Logger
		level Level
	}{
		{
			"None",
			kitlog.NewNopLogger(),
			LevelNone,
		},
		{
			"Error",
			kitlog.NewNopLogger(),
			LevelError,
		},
		{
			"Warn",
			kitlog.NewNopLogger(),
			LevelWarn,
		},
		{
			"Info",
			kitlog.NewNopLogger(),
			LevelInfo,
		},
		{
			"Debug",
			kitlog.NewNopLogger(),
			LevelDebug,
		},
		{
			"InvalidLevel",
			kitlog.NewNopLogger(),
			Level(99),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filtered := createFilteredLogger(tc.base, tc.level)

			assert.NotNil(t, filtered)
		})
	}
}

func TestNewKit(t *testing.T) {
	tests := []struct {
		name string
		opts Options
	}{
		{
			"NoOption",
			Options{},
		},
		{
			"WithMetadata",
			Options{
				Name:        "test",
				Environment: "local",
				Region:      "local",
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
			logger := NewKit(tc.opts)

			assert.NotNil(t, logger)
			assert.IsType(t, &kit{}, logger)
		})
	}
}

func TestKitWith(t *testing.T) {
	tests := []struct {
		name   string
		logger *kit
		kv     []interface{}
	}{
		{
			"OK",
			&kit{
				level:  LevelInfo,
				base:   kitlog.NewNopLogger(),
				logger: &kitlog.SwapLogger{},
			},
			[]interface{}{
				"version", "0.1.0",
				"revision", "1234567",
				"context", "test",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := tc.logger.With(tc.kv...)

			assert.NotNil(t, logger)
			assert.IsType(t, &kit{}, logger)
		})
	}
}

func TestKitGetLevel(t *testing.T) {
	tests := []struct {
		name          string
		logger        *kit
		expectedLevel Level
	}{
		{
			"None",
			&kit{
				level: LevelNone,
			},
			LevelNone,
		},
		{
			"Error",
			&kit{
				level: LevelError,
			},
			LevelError,
		},
		{
			"Warn",
			&kit{
				level: LevelWarn,
			},
			LevelWarn,
		},
		{
			"Info",
			&kit{
				level: LevelInfo,
			},
			LevelInfo,
		},
		{
			"Debug",
			&kit{
				level: LevelDebug,
			},
			LevelDebug,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			level := tc.logger.GetLevel()

			assert.Equal(t, tc.expectedLevel, level)
		})
	}
}

func TestKitSetLevel(t *testing.T) {
	tests := []struct {
		name          string
		logger        *kit
		level         string
		expectedLevel Level
	}{
		{
			"None",
			&kit{
				base:   kitlog.NewNopLogger(),
				logger: &kitlog.SwapLogger{},
			},
			"none",
			LevelNone,
		},
		{
			"Error",
			&kit{
				base:   kitlog.NewNopLogger(),
				logger: &kitlog.SwapLogger{},
			},
			"error",
			LevelError,
		},
		{
			"Warn",
			&kit{
				base:   kitlog.NewNopLogger(),
				logger: &kitlog.SwapLogger{},
			},
			"warn",
			LevelWarn,
		},
		{
			"Info",
			&kit{
				base:   kitlog.NewNopLogger(),
				logger: &kitlog.SwapLogger{},
			},
			"info",
			LevelInfo,
		},
		{
			"Debug",
			&kit{
				base:   kitlog.NewNopLogger(),
				logger: &kitlog.SwapLogger{},
			},
			"debug",
			LevelDebug,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.logger.SetLevel(tc.level)

			assert.Equal(t, tc.expectedLevel, tc.logger.level)
		})
	}
}

func TestKitLog(t *testing.T) {
	tests := []struct {
		name          string
		mockKitLogger *mockKitLogger
		message       string
		kv            []interface{}
		expectedKV    []interface{}
	}{
		{
			name: "Error",
			mockKitLogger: &mockKitLogger{
				LogOutError: errors.New("log error"),
			},
			message:    "operation failed",
			kv:         []interface{}{"reason", "no capacity"},
			expectedKV: []interface{}{"message", "operation failed", "reason", "no capacity"},
		},
		{
			name:          "Success",
			mockKitLogger: &mockKitLogger{},
			message:       "operation succeeded",
			kv:            []interface{}{"operation", "test"},
			expectedKV:    []interface{}{"message", "operation succeeded", "operation", "test"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			kl := &kit{logger: &kitlog.SwapLogger{}}
			kl.logger.Swap(tc.mockKitLogger)

			t.Run("Debug", func(t *testing.T) {
				kl.Debug(tc.message, tc.kv...)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockKitLogger.LogInKV, val)
				}
			})

			t.Run("Info", func(t *testing.T) {
				kl.Info(tc.message, tc.kv...)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockKitLogger.LogInKV, val)
				}
			})

			t.Run("Warn", func(t *testing.T) {
				kl.Warn(tc.message, tc.kv...)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockKitLogger.LogInKV, val)
				}
			})

			t.Run("Error", func(t *testing.T) {
				kl.Error(tc.message, tc.kv...)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockKitLogger.LogInKV, val)
				}
			})
		})
	}
}

func TestKitLogf(t *testing.T) {
	tests := []struct {
		name          string
		mockKitLogger *mockKitLogger
		format        string
		args          []interface{}
		expectedKV    []interface{}
	}{
		{
			name: "Error",
			mockKitLogger: &mockKitLogger{
				LogOutError: errors.New("log error"),
			},
			format:     "operation failed: %s",
			args:       []interface{}{"no capacity"},
			expectedKV: []interface{}{"message", "operation failed: no capacity"},
		},
		{
			name:          "Success",
			mockKitLogger: &mockKitLogger{},
			format:        "operation succeeded: %s",
			args:          []interface{}{"test"},
			expectedKV:    []interface{}{"message", "operation succeeded: test"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			kl := &kit{logger: &kitlog.SwapLogger{}}
			kl.logger.Swap(tc.mockKitLogger)

			t.Run("Debugf", func(t *testing.T) {
				kl.Debugf(tc.format, tc.args...)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockKitLogger.LogInKV, val)
				}
			})

			t.Run("Infof", func(t *testing.T) {
				kl.Infof(tc.format, tc.args...)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockKitLogger.LogInKV, val)
				}
			})

			t.Run("Warnf", func(t *testing.T) {
				kl.Warnf(tc.format, tc.args...)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockKitLogger.LogInKV, val)
				}
			})

			t.Run("Errorf", func(t *testing.T) {
				kl.Errorf(tc.format, tc.args...)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockKitLogger.LogInKV, val)
				}
			})
		})
	}
}

func TestKitClose(t *testing.T) {
	tests := []struct {
		name          string
		logger        *kit
		expectedError error
	}{
		{
			name:          "NoError",
			logger:        &kit{},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.logger.Close()

			assert.Equal(t, tc.expectedError, err)
		})
	}
}
