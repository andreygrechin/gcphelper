package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the interface for logging operations.
type Logger interface {
	// Debug logs a debug message.
	Debug(msg string, fields ...zap.Field)

	// Info logs an info message.
	Info(msg string, fields ...zap.Field)

	// Warn logs a warning message.
	Warn(msg string, fields ...zap.Field)

	// Error logs an error message.
	Error(msg string, fields ...zap.Field)

	// Fatal logs a fatal message and calls os.Exit(1).
	Fatal(msg string, fields ...zap.Field)

	// With returns a new logger with additional fields.
	With(fields ...zap.Field) Logger

	// WithField returns a new logger with a single additional field.
	WithField(key string, value interface{}) Logger

	// Close flushes any buffered log entries.
	Close() error
}

// ZapLogger implements the Logger interface using zap.
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLoggerForTesting creates a ZapLogger from an existing zap.Logger for testing purposes.
func NewZapLoggerForTesting(zapLogger *zap.Logger) *ZapLogger {
	return &ZapLogger{logger: zapLogger}
}

// Debug logs a debug message.
func (z *ZapLogger) Debug(msg string, fields ...zap.Field) {
	z.logger.Debug(msg, fields...)
}

// Info logs an info message.
func (z *ZapLogger) Info(msg string, fields ...zap.Field) {
	z.logger.Info(msg, fields...)
}

// Warn logs a warning message.
func (z *ZapLogger) Warn(msg string, fields ...zap.Field) {
	z.logger.Warn(msg, fields...)
}

// Error logs an error message.
func (z *ZapLogger) Error(msg string, fields ...zap.Field) {
	z.logger.Error(msg, fields...)
}

// Fatal logs a fatal message and calls os.Exit(1).
func (z *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	z.logger.Fatal(msg, fields...)
}

// With returns a new logger with additional fields.
func (z *ZapLogger) With(fields ...zap.Field) Logger {
	return &ZapLogger{
		logger: z.logger.With(fields...),
	}
}

// WithField returns a new logger with a single additional field.
func (z *ZapLogger) WithField(key string, value interface{}) Logger {
	return &ZapLogger{
		logger: z.logger.With(zap.Any(key, value)),
	}
}

// Close flushes any buffered log entries.
func (z *ZapLogger) Close() error {
	if err := z.logger.Sync(); err != nil {
		return fmt.Errorf("failed to sync logger: %w", err)
	}

	return nil
}

// NewDevelopmentLogger creates a new logger with development configuration.
// Uses zap.NewDevelopmentEncoderConfig for human-readable console output.
func NewDevelopmentLogger() (Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build development logger: %w", err)
	}

	return &ZapLogger{logger: logger}, nil
}
