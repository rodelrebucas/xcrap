package route

import (
	"github.com/labstack/echo/v4"
	"scraper/backend/env"
)

func Register(e *echo.Echo, environ *env.Var) {
	e.POST("/auth", AuthenticateHandler(environ.Secret))
	e.GET("/scrape", ScrapeHandler(environ.RedisHost, environ.RedisPass))
}