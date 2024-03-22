package logger

import (
	"io"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger
var stderrw io.Writer
var stdoutw io.Writer
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

	log = zap.Must(logCfg.Build())

	if beVerbose {
		// Set the log level so we get debug logs
		// this is dynamic and if we set it before creating the logger, it does not stick
		logCfg.Level.SetLevel(zapcore.DebugLevel)
	}

	// Create the stderr and stdout writers for os/exec
	errl, _ := zap.NewStdLogAt(log, zapcore.ErrorLevel)
	stderrw = errl.Writer()
	stdoutw = zap.NewStdLog(log).Writer()

	return log.Sugar()
}

func createIfNeeded() {
	lmux.RLock()
	defer lmux.RUnlock()

	// Do we need to create the logger?
	if log == nil {
		lmux.RUnlock()
		Set(false, false)
		lmux.RLock()
	}
}

// Get will return the logger if it already exists and generate the default if it does not
// This is thread safe
func Get() *zap.SugaredLogger {
	createIfNeeded()

	lmux.RLock()
	defer lmux.RUnlock()
	return log.Sugar()
}

func GetRaw() *zap.Logger {
	createIfNeeded()
	lmux.RLock()
	defer lmux.RUnlock()

	return log
}

func GetStdErrWriter() io.Writer {
	createIfNeeded()
	lmux.RLock()
	defer lmux.RUnlock()

	return stderrw
}

func GetStdOutWriter() io.Writer {
	createIfNeeded()
	lmux.RLock()
	defer lmux.RUnlock()

	return stdoutw
}
