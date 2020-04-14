package log

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetSingleton(t *testing.T) {
	tests := []struct {
		name   string
		logger Logger
	}{
		{
			name:   "NilLogger",
			logger: nil,
		},
		{
			name:   "KitLogger",
			logger: NewKit(Options{}),
		},
		{
			name:   "ZapLogger",
			logger: NewZap(Options{}),
		},
		{
			name:   "MockLogger",
			logger: &mockLogger{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			SetSingleton(tc.logger)
		})
	}
}

func TestGetLevel(t *testing.T) {
	tests := []struct {
		name          string
		mockLogger    *mockLogger
		level         Level
		expectedLevel Level
	}{
		{
			name:          "NoSingleton",
			mockLogger:    nil,
			expectedLevel: LevelNone,
		},
		{
			name: "None",
			mockLogger: &mockLogger{
				GetLevelOutLevel: LevelNone,
			},
			expectedLevel: LevelNone,
		},
		{
			name: "Error",
			mockLogger: &mockLogger{
				GetLevelOutLevel: LevelError,
			},
			expectedLevel: LevelError,
		},
		{
			name: "Warn",
			mockLogger: &mockLogger{
				GetLevelOutLevel: LevelWarn,
			},
			expectedLevel: LevelWarn,
		},
		{
			name: "Info",
			mockLogger: &mockLogger{
				GetLevelOutLevel: LevelInfo,
			},
			expectedLevel: LevelInfo,
		},
		{
			name: "Debug",
			mockLogger: &mockLogger{
				GetLevelOutLevel: LevelDebug,
			},
			expectedLevel: LevelDebug,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockLogger == nil {
				singleton = nil
			} else {
				singleton = tc.mockLogger
			}

			level := GetLevel()

			assert.Equal(t, tc.expectedLevel, level)
		})
	}
}

func TestSetLevel(t *testing.T) {
	tests := []struct {
		name       string
		mockLogger *mockLogger
		level      string
	}{
		{
			name:       "NoSingleton",
			mockLogger: nil,
		},
		{
			name:       "None",
			mockLogger: &mockLogger{},
			level:      "none",
		},
		{
			name:       "Error",
			mockLogger: &mockLogger{},
			level:      "error",
		},
		{
			name:       "Warn",
			mockLogger: &mockLogger{},
			level:      "warn",
		},
		{
			name:       "Info",
			mockLogger: &mockLogger{},
			level:      "info",
		},
		{
			name:       "Debug",
			mockLogger: &mockLogger{},
			level:      "debug",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockLogger == nil {
				singleton = nil
			} else {
				singleton = tc.mockLogger
			}

			SetLevel(tc.level)

			if tc.mockLogger != nil {
				assert.Equal(t, tc.level, tc.mockLogger.SetLevelInLevel)
			}
		})
	}
}

func TestLog(t *testing.T) {
	tests := []struct {
		name       string
		mockLogger *mockLogger
		message    string
		kv         []interface{}
	}{
		{
			name:       "NoSingleton",
			mockLogger: nil,
		},
		{
			name:       "OK",
			mockLogger: &mockLogger{},
			message:    "operation succeeded",
			kv:         []interface{}{"operation", "test"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockLogger == nil {
				singleton = nil
			} else {
				singleton = tc.mockLogger
			}

			t.Run("Debug", func(t *testing.T) {
				Debug(tc.message, tc.kv...)

				if tc.mockLogger != nil {
					assert.Equal(t, tc.message, tc.mockLogger.DebugInMessage)
					assert.Equal(t, tc.kv, tc.mockLogger.DebugInKV)
				}
			})

			t.Run("Info", func(t *testing.T) {
				Info(tc.message, tc.kv...)

				if tc.mockLogger != nil {
					assert.Equal(t, tc.message, tc.mockLogger.InfoInMessage)
					assert.Equal(t, tc.kv, tc.mockLogger.InfoInKV)
				}
			})

			t.Run("Warn", func(t *testing.T) {
				Warn(tc.message, tc.kv...)

				if tc.mockLogger != nil {
					assert.Equal(t, tc.message, tc.mockLogger.WarnInMessage)
					assert.Equal(t, tc.kv, tc.mockLogger.WarnInKV)
				}
			})

			t.Run("Error", func(t *testing.T) {
				Error(tc.message, tc.kv...)

				if tc.mockLogger != nil {
					assert.Equal(t, tc.message, tc.mockLogger.ErrorInMessage)
					assert.Equal(t, tc.kv, tc.mockLogger.ErrorInKV)
				}
			})
		})
	}
}

func TestLogf(t *testing.T) {
	tests := []struct {
		name       string
		mockLogger *mockLogger
		format     string
		args       []interface{}
	}{
		{
			name:       "NoSingleton",
			mockLogger: nil,
		},
		{
			name:       "OK",
			mockLogger: &mockLogger{},
			format:     "operation succeeded: %s",
			args:       []interface{}{"test"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockLogger == nil {
				singleton = nil
			} else {
				singleton = tc.mockLogger
			}

			t.Run("Debugf", func(t *testing.T) {
				Debugf(tc.format, tc.args...)

				if tc.mockLogger != nil {
					assert.Equal(t, tc.format, tc.mockLogger.DebugfInFormat)
					assert.Equal(t, tc.args, tc.mockLogger.DebugfInArgs)
				}
			})

			t.Run("Infof", func(t *testing.T) {
				Infof(tc.format, tc.args...)

				if tc.mockLogger != nil {
					assert.Equal(t, tc.format, tc.mockLogger.InfofInFormat)
					assert.Equal(t, tc.args, tc.mockLogger.InfofInArgs)
				}
			})

			t.Run("Warnf", func(t *testing.T) {
				Warnf(tc.format, tc.args...)

				if tc.mockLogger != nil {
					assert.Equal(t, tc.format, tc.mockLogger.WarnfInFormat)
					assert.Equal(t, tc.args, tc.mockLogger.WarnfInArgs)
				}
			})

			t.Run("Errorf", func(t *testing.T) {
				Errorf(tc.format, tc.args...)

				if tc.mockLogger != nil {
					assert.Equal(t, tc.format, tc.mockLogger.ErrorfInFormat)
					assert.Equal(t, tc.args, tc.mockLogger.ErrorfInArgs)
				}
			})
		})
	}
}

func TestClose(t *testing.T) {
	tests := []struct {
		name          string
		mockLogger    *mockLogger
		expectedError error
	}{
		{
			name:          "NoSingleton",
			mockLogger:    nil,
			expectedError: nil,
		},
		{
			name:          "NoError",
			mockLogger:    &mockLogger{},
			expectedError: nil,
		},
		{
			name: "Error",
			mockLogger: &mockLogger{
				CloseOutError: errors.New("failed to close"),
			},
			expectedError: errors.New("failed to close"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockLogger == nil {
				singleton = nil
			} else {
				singleton = tc.mockLogger
			}

			err := Close()

			assert.Equal(t, tc.expectedError, err)
		})
	}
}
