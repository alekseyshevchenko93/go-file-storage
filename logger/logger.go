package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type Logger struct {
	logrus *log.Logger
}

var logger *Logger

type CustomerFormatter struct {
	log.Formatter
}

func (u CustomerFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	data := e.Data

	e.Data = logrus.Fields{
		"service": os.Getenv("SERVICE_NAME"),
		"payload": data,
	}

	return u.Formatter.Format(e)
}

func Init() {
	logrus := log.New()
	logrus.SetFormatter(CustomerFormatter{&log.JSONFormatter{}})

	logger = &Logger{
		logrus,
	}
}

func GetLogger() *log.Logger {
	return logger.logrus
}
