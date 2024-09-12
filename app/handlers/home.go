package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/tools/template"

	"gohome.4gophers.ru/kovardin/gitch/views"
)

type Home struct {
	app      *pocketbase.PocketBase
	registry *template.Registry
}

func NewHome(app *pocketbase.PocketBase, registry *template.Registry) *Home {
	return &Home{
		app:      app,
		registry: registry,
	}
}

func (h *Home) Home(c echo.Context) error {
	html, err := h.registry.LoadFS(views.FS,
		"layout.html",
		"home/home.html",
	).Render(map[string]any{})

	if err != nil {
		return apis.NewNotFoundError("", err)
	}

	return c.HTML(http.StatusOK, html)
}
