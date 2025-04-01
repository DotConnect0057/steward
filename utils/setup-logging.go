package utils

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "os"
)

// SetupLogging initializes and returns a Zap SugaredLogger.
// If verbose is true, it enables detailed logging to both file and console.
func SetupLogging(verbose bool) *zap.SugaredLogger {
    var zapConfig zap.Config
    if verbose {
        zapConfig = zap.NewDevelopmentConfig()
    } else {
        zapConfig = zap.NewProductionConfig()
    }

    // Redirect logs to a file
    logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        panic("Failed to open log file: " + err.Error())
    }

    // Create a core for logging to a file
    fileCore := zapcore.NewCore(
        zapcore.NewJSONEncoder(zapConfig.EncoderConfig),
        zapcore.AddSync(logFile),
        zapConfig.Level,
    )

    // Create a core for logging to the console
    consoleCore := zapcore.NewCore(
        zapcore.NewConsoleEncoder(zapConfig.EncoderConfig),
        zapcore.AddSync(os.Stdout),
        zapConfig.Level,
    )

    var logger *zap.Logger
    if verbose {
        // Combine file and console cores
        logger = zap.New(zapcore.NewTee(fileCore, consoleCore))
    } else {
        // Use only the file core
        logger = zap.New(fileCore)
    }

    return logger.Sugar()
}