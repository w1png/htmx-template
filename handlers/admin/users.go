package admin_handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/w1png/htmx-template/errors"
	"github.com/w1png/htmx-template/models"
	"github.com/w1png/htmx-template/storage"
	"github.com/w1png/htmx-template/utils"
)

func sendUser(c echo.Context, user *models.User, is_edit bool) error {
	edit := ""
	if is_edit {
		edit = "_edit"
	}
	tmpl, err := template.ParseFiles(
		fmt.Sprintf("templates/admin/users/user%s.html", edit),
		"templates/components/loading.html",
	)

	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(c.Response().Writer, "user", user)
}

func sendUsers(c echo.Context, users []*models.User) error {
	tmpl, err := template.ParseFiles(
		"templates/admin/users/index.html",
		"templates/admin/users/user.html",
		"templates/components/loading.html",
	)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(c.Response().Writer, "users", users)
}

func sendAddUser(c echo.Context) error {
	tmpl, err := template.ParseFiles(
		"templates/admin/users/index.html",
		"templates/components/loading.html",
	)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(c.Response().Writer, "add_user_form", nil)
}

func UserIndexHandler(c echo.Context) error {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/admin/navbar.html",
		"templates/admin/users/index.html",
		"templates/admin/users/user.html",
		"templates/components/loading.html",
	)
	if err != nil {
		return err
	}

	users, err := storage.StorageInstance.GetUsers()
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(c.Response().Writer, "base", utils.MarshalResponse(c, users))
}

func UserIndexApiHandler(c echo.Context) error {
	tmpl, err := template.ParseFiles(
		"templates/admin/users/index.html",
		"templates/admin/users/user.html",
		"templates/components/loading.html",
	)
	if err != nil {
		return err
	}

	users, err := storage.StorageInstance.GetUsers()
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(c.Response().Writer, "content", utils.MarshalResponse(c, users))
}

func GetUserHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	user, err := storage.StorageInstance.GetUserById(uint(id))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Пользователь не найден")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return sendUser(c, user, false)
}

func PostUserHandler(c echo.Context) error {
	c.Response().Header().Set("HX-Reswap", "innerHTML")

	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	username := c.FormValue("username")
	password := c.FormValue("password")
	password_repeat := c.FormValue("password_repeat")
	is_admin := c.FormValue("is_admin") == "true"

	if username == "" || password == "" || password_repeat == "" {
		return c.String(http.StatusBadRequest, "Поля не могут быть пустыми")
	}

	if !models.IsUsernameValid(username) {
		return c.String(http.StatusBadRequest, models.GetUsernameRules())
	}

	if password != password_repeat {
		return c.String(http.StatusBadRequest, "Пароли не совпадают")
	}

	if !models.IsPasswordValid(password) {
		return c.String(http.StatusBadRequest, models.GetPasswordRules())
	}

	if _, err := storage.StorageInstance.GetUserByUsername(username); err == nil {
		return c.String(http.StatusBadRequest, "Пользователь с таким именем уже существует")
	} else {
		if err != nil && reflect.TypeOf(err) != reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
	}

	user, err := models.NewUser(username, password, is_admin)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}
	if err := storage.StorageInstance.CreateUser(user); err != nil {
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	c.Response().Header().Del("HX-Reswap")

	return sendUser(c, user, false)
}

func PutUserHandler(c echo.Context) error {
	c.Response().Header().Set("HX-Reswap", "innerHTML")

	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	username := c.FormValue("username")
	password := c.FormValue("password")
	password_repeat := c.FormValue("password_repeat")
	is_admin := c.FormValue("is_admin") == "true"

	if username == "" {
		return c.String(http.StatusBadRequest, "Имя пользователя не может быть пустым")
	}

	if !models.IsUsernameValid(username) {
		return c.String(http.StatusBadRequest, models.GetUsernameRules())
	}

	user, err := storage.StorageInstance.GetUserById(uint(id))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Пользователь не найден")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	user.Username = username

	if password != "" {
		if password != password_repeat {
			return c.String(http.StatusBadRequest, "Пароли не совпадают")
		}

		if !models.IsPasswordValid(password) {
			return c.String(http.StatusBadRequest, models.GetPasswordRules())
		}

		user.PasswordHash, err = user.HashPassword(password)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
	}
	if user.ID != 1 {
		user.IsAdmin = is_admin
	}

	if err := storage.StorageInstance.UpdateUser(user); err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	c.Response().Header().Del("HX-Reswap")
	c.Response().Header().Set("HX-Trigger", fmt.Sprintf("user_saved_%d", user.ID))

	return sendAddUser(c)
}

func EditUserHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	user, err := storage.StorageInstance.GetUserById(uint(id))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Пользователь не найден")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return sendUser(c, user, true)
}

func GetAddUserHandler(c echo.Context) error {
	return sendAddUser(c)
}

func DeleteUserHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := storage.StorageInstance.DeleteUserById(uint(id)); err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Пользователь не найден")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	// users, err := storage.StorageInstance.GetUsers()
	// if err != nil {
	// 	return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	// }

	return c.HTMLBlob(http.StatusOK, []byte(""))
}
