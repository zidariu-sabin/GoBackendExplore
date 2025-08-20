package crons

import (
	"GoBackendExploreMovieTracker/internal/store"
	"log"
	"time"
)

const (
	SCOPE_CLEANUP = "Cleanup"
)

type Cron struct {
	interval time.Duration
	name     string
	scope    string
	job      func(store.CronJobStore) error
}

type CronJobPipeline struct {
	crons  []Cron
	store  store.CronJobStore
	logger *log.Logger
}

func NewCronJobPipeline(crons []Cron, store store.CronJobStore, logger *log.Logger) *CronJobPipeline {
	return &CronJobPipeline{
		crons:  crons,
		store:  store,
		logger: logger,
	}
}

// defining new crons
var RunningCrons = []Cron{
	{
		interval: 1 * time.Hour,
		name:     "DeleteExpiredTokens",
		scope:    SCOPE_CLEANUP,
		job:      func(store store.CronJobStore) error { return store.DeleteExpiredTokens() },
	},
}

type CronJobPipelineControl interface {
	StartCronWorkers()
}

func (c *CronJobPipeline) StartCronWorkers() {
	for _, cron := range c.crons {
		go func(cron Cron) {
			ticker := time.NewTicker(cron.interval)
			defer ticker.Stop()
			for range ticker.C {
				if err := cron.job(c.store); err != nil {
					panic(err)
				}
				log.Printf("CRON JOB:  %s, SCOPE: %s executed successfully", cron.name, cron.scope)
			}
		}(cron)
	}
}
