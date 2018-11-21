// Package handlers ...
package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

type user struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

var (
	tempUser = &user{}
)

// @Title getuser
// @Description get user info
// @Produce  json
// @Param   X-Token     header    string     true        "token"
// @Success 200 {object}  handlers.user
// @Resource /hello
// @Router /user/get [get]
func getuser(c echo.Context) error {
	if c.Request().Header.Get("X-Token") != "qwerty" {
		return c.NoContent(http.StatusForbidden)
	}
	return c.JSON(http.StatusOK, tempUser)
}

// @Title postuser
// @Description modify user info
// @Accept  json
// @Param   X-Token     header    string     true        "token"
// @Param	user	body    handlers.user     true        "user info"
// @Success 200 string null ok
// @Resource /hello
// @Router /user/post [post]
func postuser(c echo.Context) error {
	if c.Request().Header.Get("X-Token") != "qwerty" {
		return c.NoContent(http.StatusForbidden)
	}
	err := c.Bind(tempUser)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
