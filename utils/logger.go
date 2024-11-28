package utils

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

// InitLogger initializes the logger instance for the library
func InitLogger(log *zap.Logger) {
	logger = log
}

// GetLogger returns the library's logger instance
func GetLogger() *zap.Logger {
	return logger
}

// DefaultInit initializes a default logger if none is provided
func DefaultInit() {
	if logger == nil {
		log, _ := zap.NewProduction()
		logger = log
	}
}
