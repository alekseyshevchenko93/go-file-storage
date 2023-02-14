package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func New() *log.Logger {
	logrus := log.New()
	logrus.SetFormatter(CustomerFormatter{&log.JSONFormatter{}})

	return logrus
}

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
