package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/wotmshuaisi/gomicroexample/swagger/handlers"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @host localhost:8080
// @BasePath /v1
func main() {
	e := echo.New()
	handlers.SetRouter(e)
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
