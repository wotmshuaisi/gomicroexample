package handler

import (
	"github.com/labstack/echo"
	"github.com/micro/go-micro/client"
)

// SetRouter ...
func SetRouter(c client.Client) *echo.Echo {
	e := echo.New()
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Output: os.Stdout}))
	g := e.Group("/v1")
	// basis
	gg := g.Group("/basis")
	setBasisRouter(gg, c)
	return e
}
