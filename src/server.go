package main

import (	
	"os"
	"scraper/backend/route"
	_ "scraper/backend/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// "scraper/backend/env"
	// echoSwagger "github.com/swaggo/echo-swagger"
	// "scraper/backend/fuzzy"
)

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
	e.Use(middleware.CORS())
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"*"},
	// 	AllowHeaders: []string{echo.HeaderAccessControlAllowOrigin, echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowCredentials, echo.HeaderXRequestedWith},
	// 	AllowMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
	// 	AllowCredentials: true,
	// }))
	
	/* Uncomment to load env variables and attach doc routes*/
	// env.LoadEnv()
	// e.GET("/swagger/*", echoSwagger.WrapHandler)

	/* routes */
	route.Register(e)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
