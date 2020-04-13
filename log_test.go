package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockLogger is a mock implementation of Logger
type mockLogger struct {
	WithInKV         []interface{}
	WithOutLogger    Logger
	GetLevelOutLevel Level
	SetLevelInLevel  string
	DebugInMessage   string
	DebugInKV        []interface{}
	DebugfInFormat   string
	DebugfInArgs     []interface{}
	InfoInMessage    string
	InfoInKV         []interface{}
	InfofInFormat    string
	InfofInArgs      []interface{}
	WarnInMessage    string
	WarnInKV         []interface{}
	WarnfInFormat    string
	WarnfInArgs      []interface{}
	ErrorInMessage   string
	ErrorInKV        []interface{}
	ErrorfInFormat   string
	ErrorfInArgs     []interface{}
	CloseOutError    error
}

func (m *mockLogger) With(kv ...interface{}) Logger {
	m.WithInKV = kv
	return m.WithOutLogger
}

func (m *mockLogger) GetLevel() Level {
	return m.GetLevelOutLevel
}

func (m *mockLogger) SetLevel(level string) {
	m.SetLevelInLevel = level
}

func (m *mockLogger) Debug(message string, kv ...interface{}) {
	m.DebugInMessage, m.DebugInKV = message, kv
}

func (m *mockLogger) Debugf(format string, args ...interface{}) {
	m.DebugfInFormat, m.DebugfInArgs = format, args
}

func (m *mockLogger) Info(message string, kv ...interface{}) {
	m.InfoInMessage, m.InfoInKV = message, kv
}

func (m *mockLogger) Infof(format string, args ...interface{}) {
	m.InfofInFormat, m.InfofInArgs = format, args
}

func (m *mockLogger) Warn(message string, kv ...interface{}) {
	m.WarnInMessage, m.WarnInKV = message, kv
}

func (m *mockLogger) Warnf(format string, args ...interface{}) {
	m.WarnfInFormat, m.WarnfInArgs = format, args
}

func (m *mockLogger) Error(message string, kv ...interface{}) {
	m.ErrorInMessage, m.ErrorInKV = message, kv
}

func (m *mockLogger) Errorf(format string, args ...interface{}) {
	m.ErrorfInFormat, m.ErrorfInArgs = format, args
}

func (m *mockLogger) Close() error {
	return m.CloseOutError
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		expectedLevel Level
	}{
		{"Empty", "", LevelInfo},
		{"None", "none", LevelNone},
		{"Error", "error", LevelError},
		{"Warn", "warn", LevelWarn},
		{"Info", "info", LevelInfo},
		{"Debug", "debug", LevelDebug},
		{"Invalid", "invalid", LevelNone},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			level := parseLevel(tc.level)

			assert.Equal(t, tc.expectedLevel, level)
		})
	}
}

func TestNopLogger(t *testing.T) {
	logger := NewNopLogger()
	assert.NotNil(t, logger)

	logger.With()
	logger.GetLevel()
	logger.SetLevel("none")
	logger.Debug("debug", "key", "value")
	logger.Debugf("debug %s", "this")
	logger.Debug("info", "key", "value")
	logger.Infof("info %s", "this")
	logger.Debug("warn", "key", "value")
	logger.Warnf("warn %s", "this")
	logger.Debug("error", "key", "value")
	logger.Errorf("error %s", "this")
	logger.Close()
}
