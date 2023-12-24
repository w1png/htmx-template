package admin_handlers

import (
	"html/template"

	"github.com/labstack/echo"
	"github.com/w1png/htmx-template/utils"
)

func AdminIndexHandler(c echo.Context) error {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/admin/navbar.html",
		"templates/admin/index.html",
	)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(c.Response(), "base", utils.MarshalResponse(c, nil))
}

func AdminApiIndexHandler(c echo.Context) error {
	tmpl, err := template.ParseFiles("templates/admin/index.html")
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(c.Response(), "content", utils.MarshalResponse(c, nil))
}
