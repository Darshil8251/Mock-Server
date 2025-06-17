package logger

import (
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the interface that wraps the basic logging methods
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string)
	InfoW(msg string, fields ...zap.Field)
	Warn(msg string)
	WarnW(msg string, fields ...zap.Field)
	Error(msg string, err error)
	ErrorW(msg string, err error, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	Sync() error
}

// zapLogger implements the Logger interface
type zapLogger struct {
	logger *zap.Logger
}

type loggerLevel string

const (
	ProductionLevel  loggerLevel = "PRODUCTION"
	DevelopmentLevel loggerLevel = "DEVELOPMENT"
)

var (
	instance *zapLogger
	once     sync.Once
)

// New creates a new logger instance with the specified environment
func CreateNewLogger(env string) (Logger, error) {
	var err error
	once.Do(func() {
		var config zap.Config

		if strings.Compare(env, string(ProductionLevel)) == 0 {
			config = zap.NewProductionConfig()
		} else {
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		// Create the logger
		var logger *zap.Logger
		logger, err = config.Build()
		if err != nil {
			return
		}

		// Replace the global logger
		zap.ReplaceGlobals(logger)

		instance = &zapLogger{logger}
	})

	if err != nil {
		return nil, err
	}

	return instance, nil
}

// Get returns the singleton logger instance
func Get() Logger {
	if instance == nil {
		// Default to development environment if not initialized
		logger, err := CreateNewLogger("development")
		if err != nil {
			// If we can't create a logger, create a no-op logger
			noopLogger := zap.NewNop()
			return &zapLogger{noopLogger}
		}
		return logger
	}
	return instance
}

// Debug logs a debug message
func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info message
func (l *zapLogger) Info(msg string) {
	l.logger.Info(msg)
}

// InfoW logs an info message with fields
func (l *zapLogger) InfoW(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs a warning message
func (l *zapLogger) Warn(msg string) {
	l.logger.Warn(msg)
}

// WarnW logs a warning message with fields
func (l *zapLogger) WarnW(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Error logs an error message with error
func (l *zapLogger) Error(msg string, err error) {
	l.logger.Error(msg, zap.Error(err))
}

// ErrorW logs an error message with error and additional fields
func (l *zapLogger) ErrorW(msg string, err error, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.Error(err)}, fields...)
	l.logger.Error(msg, allFields...)
}

// Fatal logs a fatal message and then calls os.Exit(1)
func (l *zapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// With creates a child logger with the given fields
func (l *zapLogger) With(fields ...zap.Field) Logger {
	return &zapLogger{l.logger.With(fields...)}
}

// Sync flushes any buffered log entries
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

// Field creates a zap.Field
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// StringField creates a zap.Field with a string value
func StringField(key, value string) zap.Field {
	return zap.String(key, value)
}

// IntField creates a zap.Field with an int value
func IntField(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// ErrorField creates a zap.Field with an error value
func ErrorField(err error) zap.Field {
	return zap.Error(err)
}


