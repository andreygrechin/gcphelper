package logger_test

import (
	"bytes"
	"testing"

	"github.com/andreygrechin/gcphelper/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewDevelopmentLogger(t *testing.T) {
	log, err := logger.NewDevelopmentLogger()
	require.NoError(t, err)
	require.NotNil(t, log)

	// test that we can close the logger (might fail in test environment, which is ok)
	_ = log.Close()
}

func TestZapLogger_LoggingMethods(t *testing.T) {
	// create a logger that writes to a buffer for testing
	var buf bytes.Buffer
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	config.OutputPaths = []string{"stdout"}

	// create a custom logger that writes to our buffer
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config.EncoderConfig),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)
	zapLogger := zap.New(core)

	// create logger for testing
	log := logger.NewZapLoggerForTesting(zapLogger)

	tests := map[string]struct {
		logFunc     func(string, ...zap.Field)
		message     string
		fields      []zap.Field
		expectLevel string
	}{
		"debug message": {
			logFunc:     log.Debug,
			message:     "debug message",
			expectLevel: "DEBUG",
		},
		"info message": {
			logFunc:     log.Info,
			message:     "info message",
			expectLevel: "INFO",
		},
		"warn message": {
			logFunc:     log.Warn,
			message:     "warning message",
			expectLevel: "WARN",
		},
		"error message": {
			logFunc:     log.Error,
			message:     "error message",
			expectLevel: "ERROR",
		},
		"debug with fields": {
			logFunc:     log.Debug,
			message:     "debug with fields",
			fields:      []zap.Field{zap.String("key", "value"), zap.Int("count", 42)},
			expectLevel: "DEBUG",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.message, tt.fields...)

			output := buf.String()
			assert.Contains(t, output, tt.expectLevel)
			assert.Contains(t, output, tt.message)

			// check fields if present
			for _, field := range tt.fields {
				if field.Type == zapcore.StringType {
					assert.Contains(t, output, field.String)
				}
			}
		})
	}
}

func TestZapLogger_With(t *testing.T) {
	log, err := logger.NewDevelopmentLogger()
	require.NoError(t, err)

	// create a new logger with additional fields
	newLog := log.With(zap.String("service", "test"), zap.Int("version", 1))
	require.NotNil(t, newLog)

	// verify it's a different instance
	assert.NotEqual(t, log, newLog)

	// test that we can close the loggers (might fail in test environment, which is ok)
	_ = log.Close()
	_ = newLog.Close()
}

func TestZapLogger_WithField(t *testing.T) {
	log, err := logger.NewDevelopmentLogger()
	require.NoError(t, err)

	tests := map[string]struct {
		key   string
		value interface{}
	}{
		"string value": {
			key:   "service",
			value: "test",
		},
		"int value": {
			key:   "port",
			value: 8080,
		},
		"bool value": {
			key:   "enabled",
			value: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			newLog := log.WithField(tt.key, tt.value)
			require.NotNil(t, newLog)
			assert.NotEqual(t, log, newLog)

			// test that we can close the logger (might fail in test environment, which is ok)
			_ = newLog.Close()
		})
	}

	// test that we can close the logger (might fail in test environment, which is ok)
	_ = log.Close()
}

func TestZapLogger_Close(t *testing.T) {
	log, err := logger.NewDevelopmentLogger()
	require.NoError(t, err)

	// test that close doesn't panic
	_ = log.Close()

	// closing again should not panic (sync is idempotent)
	_ = log.Close()
}
