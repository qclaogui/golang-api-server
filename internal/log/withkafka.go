// +build kafka

package log

import (
	"bytes"
	"net"
	"os"
	"sync"
	"time"

	"github.com/qclaogui/golang-api-server/internal"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log is global logger
	Log *zap.Logger
	// timeFormat is custom Time format
	timeFormat = "2006-01-02T15:04:05.999999999Z07:00"
	// kafka address
	brokers = "172.16.10.156:9092,172.16.10.157"
	// onceInit guarantee initialize logger only once
	onceInit sync.Once
)

type ProjectInfo struct {
	Project   string `json:"project"`
	ServerMac string `json:"server_mac"`
}

func (p *ProjectInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("project", p.Project)
	enc.AddString("server_mac", p.ServerMac)
	return nil
}
func init() {
	_ = newLogger(0, timeFormat)
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
		// Assume that we have clients for two Kafka topics. The clients implement
		// zapcore.WriteSyncer and are safe for concurrent use. (If they only
		// implement io.Writer, we can use zapcore.AddSync to add a no-op Sync
		// method. If they're not safe for concurrent use, we can add a protecting
		// mutex with zapcore.Lock.)
		kaf := internal.NewLog2Kafka(&internal.KafkaConfig{Hosts: brokers, Topic: "log-demo"})
		topicInfo := zapcore.AddSync(kaf)
		topicError := zapcore.AddSync(kaf)
		// High-priority output should also go to standard error, and low-priority
		// output should also go to standard out.
		consoleInfo := zapcore.Lock(os.Stdout)
		consoleError := zapcore.Lock(os.Stderr)
		// Configure console output.
		var useCustomTimeFormat bool
		ecfg := zap.NewProductionEncoderConfig()
		if len(timeFormat) > 0 {
			//ecfg.TimeKey = "tsp"
			ecfg.EncodeTime = customTimeEncoder
			useCustomTimeFormat = true
		}
		// Optimize the Kafka output for machine consumption and the console output
		// for human operators.
		kafkaEncoder := zapcore.NewJSONEncoder(ecfg)
		consoleEncoder := zapcore.NewJSONEncoder(ecfg)
		// Join the outputs, encoders, and level-handling functions into
		// zapcore.Cores, then tee the four cores together.
		core := zapcore.NewTee(
			zapcore.NewCore(kafkaEncoder, topicInfo, lowPriority),
			zapcore.NewCore(kafkaEncoder, topicError, highPriority),
			//
			zapcore.NewCore(consoleEncoder, consoleInfo, lowPriority),
			zapcore.NewCore(consoleEncoder, consoleError, highPriority),
		)
		pj := ProjectInfo{
			Project:   "golang-api-server",
			ServerMac: GetMacAddr(),
		}
		// From a zapcore.Core, it's easy to construct a Logger.
		Log = zap.New(core).With(zap.Object("info", &pj))
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

// GetMacAddr gets the MAC hardware address of the hosts machine
func GetMacAddr() (addr string) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				// Don't use random as we have a real address
				addr = i.HardwareAddr.String()
				break
			}
		}
	}
	return
}
