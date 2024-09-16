package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/pocketbase/pocketbase/tools/template"
	"go.uber.org/fx"

	"gohome.4gophers.ru/kovardin/gitch/app/handlers"
	"gohome.4gophers.ru/kovardin/gitch/app/settings"
	"gohome.4gophers.ru/kovardin/gitch/app/tasks"
	_ "gohome.4gophers.ru/kovardin/gitch/migrations"
	"gohome.4gophers.ru/kovardin/gitch/static"
)

func main() {
	fx.New(
		handlers.Module,
		tasks.Module,

		fx.Provide(pocketbase.New),
		fx.Provide(template.NewRegistry),
		fx.Provide(settings.New),

		fx.Invoke(
			migration,
		),
		fx.Invoke(
			routing,
		),
		fx.Invoke(
			task,
		),
	).Run()
}

func routing(
	app *pocketbase.PocketBase,
	lc fx.Lifecycle,
	settings *settings.Settings,
	home *handlers.Home,
) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/", home.Home)

		// static
		e.Router.GET("/static/*", func(c echo.Context) error {
			p := c.PathParam("*")

			path, err := url.PathUnescape(p)
			if err != nil {
				return fmt.Errorf("failed to unescape path variable: %w", err)
			}

			err = c.FileFS(path, static.FS)
			if err != nil && errors.Is(err, echo.ErrNotFound) {
				return c.FileFS("index.html", static.FS)
			}

			return err
		})

		return nil

	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.Start()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

// go run cmd/gitch/main.go migrate collections --dir ./data/
func migration(app *pocketbase.PocketBase) {
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})
}

func task(app *pocketbase.PocketBase, settings *settings.Settings, sync *tasks.Sync) {
	schedulers := map[string]*cron.Cron{}

	app.OnModelAfterCreate("services").Add(func(e *core.ModelEvent) error {
		service, err := app.Dao().FindRecordById("services", e.Model.GetId())
		if err != nil {
			return err
		}

		schedulers[service.Id] = cron.New()

		update(app, schedulers, sync)

		return nil
	})

	app.OnModelAfterDelete("services").Add(func(e *core.ModelEvent) error {
		if sch, ok := schedulers[e.Model.GetId()]; ok {
			if sch.HasStarted() {
				sch.Stop()
			}
		}
		delete(schedulers, e.Model.GetId())

		update(app, schedulers, sync)

		return nil
	})

	app.OnModelAfterUpdate("services").Add(func(e *core.ModelEvent) error {
		update(app, schedulers, sync)

		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		services, err := app.Dao().FindRecordsByFilter(
			"services",
			"true = true",
			"-created",
			100,
			0,
		)

		if err != nil {
			panic(err)
		}

		for _, service := range services {
			schedulers[service.Id] = cron.New()
		}

		update(app, schedulers, sync)

		return nil
	})
}

func update(
	app *pocketbase.PocketBase,
	schedulers map[string]*cron.Cron,
	sync *tasks.Sync,
) {

	for id, scheduler := range schedulers {
		service, err := app.Dao().FindRecordById("services", id)
		if err != nil {
			app.Logger().Error("error on process service")

			continue
		}

		if scheduler.HasStarted() {
			scheduler.Stop()
		}

		if !service.GetBool("enabled") {
			continue
		}

		period := service.GetString("period")
		if period == "" {
			period = "0 */1 * * *"
		}

		scheduler.MustAdd("sync", period, func() {
			sync.Do(service)
		})

		scheduler.Start()
		app.Logger().Info("scheduled", "period", period, "id", id)
	}
}
