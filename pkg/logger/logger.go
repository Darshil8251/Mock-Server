package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the interface that wraps the basic logging methods
type Logger interface {
	Debug(msg string, fields map[string]any)
	Info(msg string)
	InfoW(msg string, fields map[string]any)
	Warn(msg string, err error)
	WarnW(msg string, err error, fields map[string]any)
	Error(msg string, err error)
	ErrorW(msg string, err error, fields map[string]any)
	With(fields map[string]any) Logger
	Sync() error
}

// zapLogger implements the Logger interface
type zapLogger struct {
	logger *zap.Logger
}

type loggerLevel string

const (
	productionLevel  loggerLevel = "PRODUCTION"
	developmentLevel loggerLevel = "DEVELOPMENT"
)

var (
	loggerInstance *zapLogger
	_              Logger = (*zapLogger)(nil)
)

// New creates a new logger instance with the specified environment
func CreateNewLogger() (Logger, error) {
	var err error
	var config zap.Config
	level := os.Getenv("LEVEL")
	fmt.Printf("level: %s\n", level)

	if level == string(productionLevel) {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	// Replace the global logger
	zap.ReplaceGlobals(logger)
	logger = logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))

	loggerInstance = &zapLogger{logger}

	return loggerInstance, nil
}

// GetLogger returns the singleton logger instance
func GetLogger() Logger {
	if loggerInstance == nil {
		logger, err := CreateNewLogger()
		if err != nil {
			return &zapLogger{zap.NewNop()}
		}

		return logger
	}
	return loggerInstance
}

// Debug logs a debug message
func (l *zapLogger) Debug(msg string, fields map[string]any) {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, Field(key, value))
	}
	l.logger.Debug(msg, zapFields...)
}

// Info logs an info message
func (l *zapLogger) Info(msg string) {
	l.logger.Info(msg)
}

// InfoW logs an info message with fields
func (l *zapLogger) InfoW(msg string, fields map[string]any) {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, Field(key, value))
	}
	l.logger.Info(msg, zapFields...)
}

// Warn logs a warning message
func (l *zapLogger) Warn(msg string, err error) {
	l.logger.Warn(msg, Field("error", err))
}

// WarnW logs a warning message with fields
func (l *zapLogger) WarnW(msg string, err error, fields map[string]any) {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, Field(key, value))
	}
	l.logger.Warn(msg, append([]zap.Field{Field("error", err)}, zapFields...)...)
}

// Error logs an error message with error
func (l *zapLogger) Error(msg string, err error) {
	l.logger.Error(msg, zap.Error(err))
}

// ErrorW logs an error message with error and additional fields
func (l *zapLogger) ErrorW(msg string, err error, fields map[string]any) {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, Field(key, value))
	}
	l.logger.Error(msg, append([]zap.Field{zap.Error(err)}, zapFields...)...)
}

// Fatal logs a fatal message and then calls os.Exit(1)
func (l *zapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// With creates a child logger with the given fields
func (l *zapLogger) With(fields map[string]any) Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, Field(key, value))
	}
	return &zapLogger{l.logger.With(zapFields...)}
}

// Sync flushes any buffered log entries
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

// Field creates a zap.Field
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}
