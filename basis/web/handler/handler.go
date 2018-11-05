package handler

import "github.com/labstack/echo"

// SetRouter ...
func SetRouter() *echo.Echo {
	e := echo.New()
	g := e.Group("/v1")
	// basis
	gg := g.Group("/basis")
	setBasisRouter(gg)
	return e
}
