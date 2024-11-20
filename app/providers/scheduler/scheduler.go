package scheduler

import (
	"github.com/mrrizkin/omniscan/app/providers/logger"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	*cron.Cron
	log *logger.Logger
}

func (*Scheduler) Construct() interface{} {
	return func(log *logger.Logger) (*Scheduler, error) {
		return &Scheduler{
			Cron: cron.New(),
			log:  log,
		}, nil
	}
}

func (s *Scheduler) Add(spec string, cmd func()) {
	s.AddFunc(spec, cmd)
}

func (s *Scheduler) Start() {
	s.log.Info("starting cron", "entries", len(s.Entries()))
	s.Cron.Start()
}
