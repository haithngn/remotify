package logging

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func InitLogger(lvl log.Level, logDir string) error {
	log.SetLevel(lvl)

	lumberjackLogger := &lumberjack.Logger{
		Filename:   logDir + "/logs/remotify.log",
		MaxSize:    50, // MB
		MaxBackups: 20,
		MaxAge:     14, // days
		Compress:   true,
	}
	log.SetOutput(io.MultiWriter(os.Stdout, lumberjackLogger))

	log.SetReportCaller(true)

	log.SetFormatter(&log.TextFormatter{
		PadLevelText:    true,
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return nil
}
