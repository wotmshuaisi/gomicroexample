package handlers

import (
	"github.com/labstack/echo"
)

// SetRouter ...
func SetRouter(e *echo.Echo) {
	g := e.Group("/v1")
	g.GET("/user/get", getuser)
	g.POST("/user/post", postuser)
}
