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
	w.log.Info("worker.removeLeastUsedFiles")
}

func (w *worker) Start() {
	w.cron.Start()
	w.log.Info("worker.started")
}

func (w *worker) Stop() {
	w.cron.Stop()
	w.log.Info("worker.stopped")
}

func NewWorker(log *logrus.Logger) *worker {
	cron := cron.New()
	worker := &worker{
		log,
		cron,
	}

	cron.AddFunc("* * * * *", worker.removeLeastUsedFiles)

	return worker
}
