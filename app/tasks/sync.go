package tasks

import (
	"github.com/pocketbase/pocketbase"

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

func (t *Sync) Do() {
	projects, err := t.app.Dao().FindRecordsByFilter(
		"projects",
		"enabled = true",
		"-created",
		1000,
		0,
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
			t.settings.SSHKey(),
		)

		if err := s.Sync(); err != nil {
			t.app.Logger().Error("error on sync project", "name", p.GetString("name"), "err", err)
		}
	}
}
