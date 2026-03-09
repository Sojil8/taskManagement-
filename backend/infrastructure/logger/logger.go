package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
)

// Log is the global logger instance
var Log *zap.Logger

// InitLogger initializes a Zap logger based on the environment
func InitLogger() {
	var err error

	ginMode := os.Getenv("GIN_MODE")
	// If it is release mode, use the production JSON logger
	if ginMode == "release" {
		Log, err = zap.NewProduction()
	} else {
		// Otherwise, use a console-friendly development logger
		Log, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	// Replace the global standard log with zap so third-party logs are routed appropriately
	zap.ReplaceGlobals(Log)
}

// Sync flushes any buffered log entries. Call deferred in main.
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
