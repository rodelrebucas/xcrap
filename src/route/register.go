package route

import (
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, secret string) {
	e.POST("/auth", AuthenticateHandler(secret))
	e.GET("/scrape", ScrapeHandler)
}