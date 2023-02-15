package workers

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type FileDeleter interface {
	GetLeastUsedFilesIds(date *time.Time) ([]int, error)
}

type worker struct {
	log        *logrus.Logger
	cron       *cron.Cron
	repository FileDeleter
}

func (w *worker) removeLeastUsedFiles() {
	log := w.log

	log.Info("worker.removeLeastUsedFiles.started")
	now := time.Now()
	ids, err := w.repository.GetLeastUsedFilesIds(&now)

	if err != nil {
		log.WithField("message", err.Error()).Error("worker.removeLeastUsedFiles.error")
		return
	}

	log.WithField("ids", ids).Error("worker.removeLeastUsedFiles.success")
}

func (w *worker) Start() {
	w.cron.Start()
	w.log.Info("worker.started")
}

func (w *worker) Stop() {
	w.cron.Stop()
	w.log.Info("worker.stopped")
}

func NewWorker(log *logrus.Logger, repository FileDeleter) *worker {
	cron := cron.New()
	worker := &worker{
		log,
		cron,
		repository,
	}

	cron.AddFunc("* * * * *", worker.removeLeastUsedFiles)

	return worker
}
