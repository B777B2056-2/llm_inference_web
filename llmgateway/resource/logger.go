package resource

import (
	"fmt"
	"io"
	"llm_online_interence/llmgateway/confparser"
	"os"
	"path"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func initLogger() {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	switch confparser.ResourceConfig.Logger.Level {
	case LoggerLevelDebug:
		Logger.SetLevel(logrus.DebugLevel)
	case LoggerLevelInfo:
		Logger.SetLevel(logrus.InfoLevel)
	case LoggerLevelWarn:
		Logger.SetLevel(logrus.WarnLevel)
	case LoggerLevelError:
		Logger.SetLevel(logrus.ErrorLevel)
	case LoggerLevelFatal:
		Logger.SetLevel(logrus.FatalLevel)
	}

	if err := os.MkdirAll(confparser.ResourceConfig.Logger.OutPutPath, 0755); err != nil {
		panic(fmt.Errorf("failed to create log dir: %v", err))
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   path.Join(confparser.ResourceConfig.Logger.OutPutPath, "gateway.log"),
		MaxSize:    confparser.ResourceConfig.Logger.MaxSingleFileSizeMB,
		MaxBackups: confparser.ResourceConfig.Logger.MaxBackups,
		MaxAge:     confparser.ResourceConfig.Logger.MaxStorageAgeInDays,
		Compress:   true,
	}
	Logger.SetOutput(io.MultiWriter(os.Stdout, lumberjackLogger))
}

// TraceIdHook 用于在日志中添加 traceId
type TraceIdHook struct {
	TraceId string
}

func NewTraceIdHook(traceId string) logrus.Hook {
	hook := TraceIdHook{
		TraceId: traceId,
	}
	return &hook
}

func (hook *TraceIdHook) Fire(entry *logrus.Entry) error {
	entry.Data["TraceId"] = hook.TraceId
	return nil
}

func (hook *TraceIdHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
