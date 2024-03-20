package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger
var lmux sync.RWMutex

// Set will configure the logger based on inputs and return the created logger
// This is thread safe
func Set(beVerbose bool, useJson bool) *zap.SugaredLogger {
	lmux.Lock()
	defer lmux.Unlock()
	logCfg := zap.NewProductionConfig()
	logCfg.DisableStacktrace = true

	if !useJson {
		logCfg.Encoding = "console" // "json" is the default in production configs
	}

	if beVerbose {
		logCfg.DisableStacktrace = false // We probably want this data for debugging
	}

	// Set the time format to iso8601 timestamp format
	logCfg.EncoderConfig.TimeKey = "timestamp"
	logCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	log = zap.Must(logCfg.Build()).Sugar()

	if beVerbose {
		// Set the log level so we get debug logs
		// this is dynamic and if we set it before creating the logger, it does not stick
		logCfg.Level.SetLevel(zapcore.DebugLevel)
	}

	return log
}

// Get will return the logger if it already exists and generate the default if it does not
// This is thread safe
func Get() *zap.SugaredLogger {
	lmux.RLock()
	if log == nil {
		lmux.RUnlock()
		return Set(false, false) // Not verbose, no json
	}

	defer lmux.RUnlock()

	return log
}
