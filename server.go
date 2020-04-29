package main

import (
	"scraper/backend/env"
	"scraper/backend/route"
	_ "scraper/backend/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	// "scraper/backend/fuzzy"
)

var environ *env.Var = env.LoadEnv()

// @title sample API docs
// @version 0.1.0
// @description Backend server for the sample app
// @tag.name Auth
// @tag.description user authentication operations
// @termsOfService http://swagger.io/terms/
// @BasePath /
func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	if environ.Env == "development" {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// routes
	// route.Auth(e, environ.Secret)
	// route.Scrape(e
	route.Register(e, environ.Secret)

	e.Logger.Fatal(e.Start(":" + environ.Port))
}
