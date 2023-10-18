package worker

import (
	"context"
	"gitch/app/config"
	"gitch/pkg/logger"
	"gitch/pkg/syncer"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

type Worker struct {
	log *logger.Logger
	cfg config.Application
}

func New(lc fx.Lifecycle, cfg config.Application, log *logger.Logger) *Worker {
	w := &Worker{
		log: log,
		cfg: cfg,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			w.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			w.Stop()
			return nil
		},
	})

	return w
}

func (w *Worker) Start() {
	w.log.Info("start worker")
	ticker := time.NewTicker(time.Minute)
	go func() {
		for ; true; <-ticker.C {

			w.process()
		}

		//for range ticker.C {
		//	u.process()
		//}
	}()
}

func (w *Worker) Stop() {

}

func (w *Worker) process() {
	for name, p := range w.cfg.Projects {
		if !p.Enable {
			continue
		}

		s := syncer.New(
			p.From,
			p.To,
			w.cfg.Key,
		)

		if err := s.Sync(); err != nil {
			w.log.Error("error on sync project", zap.String("name", name))
		}
	}
}
