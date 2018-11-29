// +build !kafka

package log

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log is global logger
	Log *zap.Logger
	// timeFormat is custom Time format
	timeFormat = "2006-01-02T15:04:05.999999999Z07:00"
	// onceInit guarantee initialize logger only once
	onceInit sync.Once
)

func init() {
	newLogger(0, timeFormat)
}

// customTimeEncoder encode Time to our custom format
// This example how we can customize zap default functionality
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(timeFormat))
}

// Init initializes log by input parameters
// lvl - global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
// timeFormat - custom time format for logger of empty string to use default
func newLogger(lvl int, timeFormat string) error {
	var err error

	onceInit.Do(func() {
		// First, define our level-handling logic.
		globalLevel := zapcore.Level(lvl)
		// High-priority output should also go to standard error, and low-priority
		// output should also go to standard out.
		// It is usefull for Kubernetes deployment.
		// Kubernetes interprets os.Stdout log items as INFO and os.Stderr log items
		// as ERROR by default.
		highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= globalLevel && lvl < zapcore.ErrorLevel
		})

		// High-priority output should also go to standard error, and low-priority
		// output should also go to standard out.
		consoleInfos := zapcore.Lock(os.Stdout)
		consoleErrors := zapcore.Lock(os.Stderr)

		// Configure console output.
		var useCustomTimeFormat bool
		ecfg := zap.NewProductionEncoderConfig()
		if len(timeFormat) > 0 {
			ecfg.EncodeTime = customTimeEncoder
			useCustomTimeFormat = true
		}

		// Optimize the Kafka output for machine consumption and the console output
		// for human operators.
		consoleEncoder := zapcore.NewJSONEncoder(ecfg)

		// Join the outputs, encoders, and level-handling functions into
		// zapcore.Cores, then tee the four cores together.
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
			zapcore.NewCore(consoleEncoder, consoleInfos, lowPriority),
		)

		// From a zapcore.Core, it's easy to construct a Logger.
		Log = zap.New(core)

		zap.RedirectStdLog(Log)
		if !useCustomTimeFormat {
			Log.Error("time format for logger is not provided - use zap default")
		}
	})

	return err
}

//Debug log
func Debug(msg string, fields ...zap.Field) { Log.Debug(msg, fields...) }

// Info log
func Info(msg string, fields ...zap.Field) { Log.Info(msg, fields...) }

// Warn log
func Warn(msg string, fields ...zap.Field) { Log.Warn(msg, fields...) }

// Error log
func Error(msg string, fields ...zap.Field) { Log.Error(msg, fields...) }

// Fatal log
func Fatal(msg string, fields ...zap.Field) { Log.Fatal(msg, fields...) }
