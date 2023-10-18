package main

import (
	"fmt"
	"gitch/pkg/syncer"
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"gitch/app/config"
	"gitch/app/server"
	"gitch/app/worker"
	"gitch/pkg/logger"
)

func main() {
	app := &cli.App{
		Name:  "getapp",
		Usage: "make an explosive entrance",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "config",
				Usage: "environment",
			},
			&cli.StringFlag{
				Name:  "configs",
				Value: "./",
				Usage: "configs path",
			},
		},
		Commands: cli.Commands{
			&cli.Command{
				Name: "server",
				Action: func(c *cli.Context) error {
					fmt.Println(os.Getenv("HOME"))
					app(c,
						fx.Invoke(func(server *server.Server) {}),
						fx.Invoke(func(worker *worker.Worker) {}),
						fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
							return &fxevent.ZapLogger{Logger: log}
						}),
					).Run()

					return nil
				},
			},
			&cli.Command{
				Name: "sync",
				Action: func(c *cli.Context) error {
					app(c,
						fx.Invoke(func(log *logger.Logger, cfg config.Application) {
							snc := syncer.New(
								"git@gitflic.ru:getapp/example.git",
								"git@github.com:kovardin/example.git",
								os.Getenv("HOME")+"/.ssh/id_rsa",
							)
							if err := snc.Sync(); err != nil {
								log.Error("error on sync", zap.Error(err))
							}

							os.Exit(0)
						}),
						fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
							return &fxevent.ZapLogger{Logger: log}
						}),
					).Run()

					return nil
				},
			},
		},
	}

	app.Run(os.Args)
}

func app(c *cli.Context, opts ...fx.Option) *fx.App {
	env := c.String("env")
	cfg := c.String("configs")

	opts = append(opts,
		fx.Provide(
			func() config.Config {
				return config.New(env, cfg)
			},
			logger.New,
			server.New,
			worker.New,
		),
	)

	return fx.New(
		opts...,
	)
}
