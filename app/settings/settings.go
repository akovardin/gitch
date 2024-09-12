package settings

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

type Settings struct {
	app *pocketbase.PocketBase
}

func New(app *pocketbase.PocketBase) *Settings {
	return &Settings{
		app: app,
	}
}

func (s *Settings) SSHKey() string {
	record, err := s.app.Dao().FindFirstRecordByFilter("settings", "key = {:key}", dbx.Params{"key": "ssh_key"})
	if err != nil {
		s.app.Logger().Warn("error on search settings", "err", err)
		return ""
	}

	return record.GetString("value")
}

func (s *Settings) Period(def string) string {
	record, err := s.app.Dao().FindFirstRecordByFilter("settings", "key = {:key}", dbx.Params{"key": "period"})
	if err != nil {
		s.app.Logger().Warn("error on search settings", "err", err)
		return def
	}

	if record.GetString("value") == "" {
		return def
	}

	return record.GetString("value")
}
