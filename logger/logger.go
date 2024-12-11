package logger

import (
	"bytes"
	"fmt"
	"path"

	"github.com/lanthora/cacao/argp"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	logger = logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&logFormatter{})

	switch argp.Get("loglevel", "info") {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	}

	if argp.Get("logfile", "disable") == "enable" {
		logger.SetOutput(&lumberjack.Logger{
			Filename:   path.Join(argp.Get("storage", "."), "logs/cacao.log"),
			MaxSize:    10,    // 每个日志文件最大尺寸（M）
			MaxBackups: 5,     // 保留旧文件的最大个数
			MaxAge:     3,     // 保留旧文件的最大天数
			Compress:   true,  // 是否压缩/归档旧文件
		})
	}

	Info("loglevel=[%v]", logger.GetLevel().String())
}

var logger *logrus.Logger

type logFormatter struct{}

func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	b.WriteString(msg)
	return b.Bytes(), nil
}

func Fatal(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Info(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Debug(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}
