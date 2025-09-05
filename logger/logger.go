package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	
	// 設置日誌格式
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})
	
	// 設置日誌級別
	Log.SetLevel(logrus.InfoLevel)
	
	// 設置輸出到文件和控制台
	setupLogOutput()
}

func setupLogOutput() {
	// 創建logs目錄
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		Log.WithError(err).Error("無法創建日誌目錄")
		return
	}
	
	// 按日期分割日誌文件
	logFile := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	
	// 打開日誌文件
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.WithError(err).Error("無法打開日誌文件")
		return
	}
	
	// 同時輸出到文件和控制台
	Log.SetOutput(io.MultiWriter(os.Stdout, file))
}

// 設置日誌級別
func SetLevel(level string) {
	switch level {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}
}

// 便捷的日誌方法
func Debug(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Log.WithFields(fields[0]).Debug(msg)
	} else {
		Log.Debug(msg)
	}
}

func Info(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Log.WithFields(fields[0]).Info(msg)
	} else {
		Log.Info(msg)
	}
}

func Warn(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Log.WithFields(fields[0]).Warn(msg)
	} else {
		Log.Warn(msg)
	}
}

func Error(msg string, err error, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Log.WithError(err).WithFields(fields[0]).Error(msg)
	} else {
		Log.WithError(err).Error(msg)
	}
}

func Fatal(msg string, err error, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Log.WithError(err).WithFields(fields[0]).Fatal(msg)
	} else {
		Log.WithError(err).Fatal(msg)
	}
}

// 帶上下文的日誌記錄器
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return Log.WithField(key, value)
}

func WithError(err error) *logrus.Entry {
	return Log.WithError(err)
}
