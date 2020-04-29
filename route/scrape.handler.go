package route

import (
	"net/http"
	"github.com/labstack/echo/v4"
)


type Task struct{
	Name string `json:name xml:name`
}

func ScrapeHandler(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, &Task{
		Name: "Scrape one",
	})
}