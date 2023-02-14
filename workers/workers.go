package workers

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type worker struct {
	log  *logrus.Logger
	cron *cron.Cron
}

func (w *worker) removeLeastUsedFiles() {

}

func (w *worker) Start() {
	w.cron.Start()
}

func (w *worker) Stop() {
	w.cron.Stop()
}

func NewWorker(log *logrus.Logger) *worker {
	cron := cron.New()
	worker := &worker{
		log,
		cron,
	}

	cron.AddFunc("*/5 * * * *", worker.removeLeastUsedFiles)

	return worker
}
