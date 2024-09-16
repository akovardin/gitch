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

func (s *Settings) Value(key string) string {
	record, err := s.app.Dao().FindFirstRecordByFilter("settings", "key = {:key}", dbx.Params{"key": "key"})
	if err != nil {
		s.app.Logger().Warn("error on search settings", "err", err)
		return ""
	}

	return record.GetString("value")
}
