package user_handlers

import (
	"html/template"

	"github.com/labstack/echo"
	"github.com/w1png/htmx-template/utils"
)

func IndexApiHandler(c echo.Context) error {
	tmpl, err := template.ParseFiles("templates/user/index.html")
	if err != nil {
		return err
	}

	return tmpl.Execute(c.Response(), utils.MarshalResponse(c, nil))
}

func IndexHandler(c echo.Context) error {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/navbar.html",
		"templates/user/index.html",
	)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(c.Response(), "base", utils.MarshalResponse(c, nil))
}
