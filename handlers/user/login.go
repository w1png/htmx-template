package user_handlers

import (
	"html/template"
	"net/http"
	"reflect"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/w1png/htmx-template/config"
	"github.com/w1png/htmx-template/errors"
	"github.com/w1png/htmx-template/storage"
	"github.com/w1png/htmx-template/utils"
)

func LoginPageApiHandler(c echo.Context) error {
	if c.Request().Context().Value("user") != nil {
		c.Response().Header().Set("HX-Redirect", "/admin")
		c.Response().Header().Set("HX-Replace-Url", "/admin")
		return c.Redirect(http.StatusFound, "/admin")
	}

	tmpl, err := template.ParseFiles("templates/user/login.html")
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(c.Response(), "content", utils.MarshalResponse(c, nil))
}

func LoginPageHandler(c echo.Context) error {
	if c.Request().Context().Value("user") != nil {
		return c.Redirect(http.StatusFound, "/admin")
	}

	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/navbar.html",
		"templates/user/login.html",
	)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(c.Response(), "base", utils.MarshalResponse(c, nil))
}

func PostLoginHandler(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" {
		return c.String(http.StatusBadRequest, "Имя пользователя не может быть пустым")
	}

	if password == "" {
		return c.String(http.StatusBadRequest, "Пароль не может быть пустым")
	}

	user, err := storage.StorageInstance.GetUserByUsername(username)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusBadRequest, "Неправильный логин или пароль")
		}
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	if !user.ComparePassword(password) {
		return c.String(http.StatusBadRequest, "Неправильный логин или пароль")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
	})

	tokenString, err := token.SignedString([]byte(config.ConfigInstance.JWTSecret))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	http.SetCookie(c.Response(), &http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
		Path:  "/",
	})

	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/admin/navbar.html",
		"templates/admin/index.html",
	)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(c.Response().Writer, "body", utils.MarshalResponse(c, nil))
}
