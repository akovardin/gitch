package tasks

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"

	"gohome.4gophers.ru/kovardin/gitch/app/settings"
	"gohome.4gophers.ru/kovardin/gitch/pkg/syncer"
)

type Sync struct {
	app      *pocketbase.PocketBase
	settings *settings.Settings
}

func NewSync(app *pocketbase.PocketBase, settings *settings.Settings) *Sync {
	w := &Sync{
		app:      app,
		settings: settings,
	}

	return w
}

func (t *Sync) Do(service *models.Record) {
	projects, err := t.app.Dao().FindRecordsByFilter(
		"projects",
		"enabled = true && service = {:service}",
		"-created",
		1000,
		0,
		dbx.Params{"service": service.Id},
	)

	if err != nil {
		t.app.Logger().Warn("error on get projects", "err", err)
	}

	for _, p := range projects {
		if !p.GetBool("enabled") {
			continue
		}

		s := syncer.New(
			p.GetString("from"),
			p.GetString("to"),
			service.GetString("key"),
		)

		if err := s.Sync(); err != nil {
			t.app.Logger().Error("error on sync project", "name", p.GetString("name"), "err", err)
		}
	}
}
