package admin_handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GatherIndexUsers(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("", IndexHandler)
	admin_api_group.GET("", IndexApiHandler)
}

func IndexApiHandler(c echo.Context) error {
	c.Response().Header().Set("HX-Redirect", "/admin/users")
	c.Response().Header().Set("HX-Replace-Url", "/admin/users")
	return c.Redirect(http.StatusFound, "/admin/users")
}

func IndexHandler(c echo.Context) error {
	c.Response().Header().Set("HX-Redirect", "/admin/users")
	c.Response().Header().Set("HX-Replace-Url", "/admin/users")
	return c.Redirect(http.StatusFound, "/admin/users")
}
