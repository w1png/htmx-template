package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/w1png/htmx-template/config"
	admin_handlers "github.com/w1png/htmx-template/handlers/admin"
	user_handlers "github.com/w1png/htmx-template/handlers/user"
	"github.com/w1png/htmx-template/middleware"
)

type HTTPServer struct {
	echo *echo.Echo
}

func NewHTTPServer() *HTTPServer {
	server := &HTTPServer{
		echo: echo.New(),
	}

	user_page_group := server.echo
	user_page_group.Use(middleware.UseAuth)
	user_api_group := server.echo.Group("/api")

	admin_page_group := server.echo.Group("/admin", middleware.UseAdmin)
	admin_api_group := admin_page_group.Group("/api")

	server.echo.Use(echoMiddleware.Logger())
	server.echo.Use(echoMiddleware.Recover())

	server.echo.Static("/static", "static")

	gather_funcs := []func(*echo.Echo, *echo.Group, *echo.Group, *echo.Group){
		user_handlers.GatherLoginRoutes,
		admin_handlers.GatherUsersRoutes,
		admin_handlers.GatherIndexUsers,
	}

	for _, f := range gather_funcs {
		f(user_page_group, user_api_group, admin_page_group, admin_api_group)
	}

	return server
}

func (s *HTTPServer) Run() error {
	return s.echo.Start(fmt.Sprintf(":%s", config.ConfigInstance.Port))
}
