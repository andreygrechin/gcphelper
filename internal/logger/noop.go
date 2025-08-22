package logger

import "go.uber.org/zap"

// NewNoOpLogger creates a new no-op logger that implements the Logger interface.
// It wraps zap's nop logger inside our ZapLogger to satisfy the interface,
// keeping behavior close to the vanilla zap no-op logger.
func NewNoOpLogger() Logger {
	return &ZapLogger{logger: zap.NewNop()}
}
