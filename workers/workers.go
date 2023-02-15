package workers

import (
	"context"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type FileDeleter interface {
	GetLeastUsedFilesIds(date *time.Time) ([]int, error)
}

type worker struct {
	context    context.Context
	log        *logrus.Logger
	cron       *cron.Cron
	repository FileDeleter
}

func (w *worker) removeLeastUsedFiles() {
	log := w.log

	log.Info("worker.removeLeastUsedFiles.started")
	duration, err := time.ParseDuration(os.Getenv("FILE_RETENTION_PERIOD"))

	if err != nil {
		log.WithField("message", err.Error()).Error("worker.removeLeastUsedFiles.failedToParseDurationError")
		return
	}

	now := time.Now().Add(duration)
	ids, err := w.repository.GetLeastUsedFilesIds(&now)

	if err != nil {
		log.WithField("message", err.Error()).Error("worker.removeLeastUsedFiles.queryError")
		return
	}

	log.WithField("ids", ids).Info("worker.removeLeastUsedFiles.success")
}

func (w *worker) Start() {
	w.cron.Start()
	w.log.Info("worker.started")
}

func (w *worker) Stop() {
	w.cron.Stop()
	w.log.Info("worker.stopped")
}

func NewWorker(context context.Context, log *logrus.Logger, repository FileDeleter) *worker {
	cron := cron.New()
	worker := &worker{
		context,
		log,
		cron,
		repository,
	}

	cron.AddFunc("* * * * *", worker.removeLeastUsedFiles)

	return worker
}
