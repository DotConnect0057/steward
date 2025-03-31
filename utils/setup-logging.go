package utils

import (
    "go.uber.org/zap"
)

// Logging initializes and returns a Zap SugaredLogger
func SetupLogging() *zap.SugaredLogger {
    // Initialize Zap logger
    zapLogger, err := zap.NewProduction() // Use zap.NewDevelopment() for a more human-readable format
    if err != nil {
        panic("Failed to initialize logger: " + err.Error())
    }
    return zapLogger.Sugar()
}