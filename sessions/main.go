package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boj/redistore"

	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo"
)

func get(c echo.Context) error {
	sess, err := session.Get("name", c)
	if err != nil || sess.Values["name"] == nil {
		return c.HTML(http.StatusOK, `
		<button onclick="rr()">click</button>
<script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
<script>
    function rr() {
        $.ajax({
            url: '/',
            type: 'post',
            async: false,
            contentType: "application/json",
            data: JSON.stringify({
                "name": "wotmshuaisi"
            }),
            success: function (arg) {
                document.write(arg)
            }
        });
    }
</script>
		`)
	}
	c.SetCookie(&http.Cookie{})
	return c.HTML(200, sess.Values["name"].(string))
}

func post(c echo.Context) error {
	var payload map[string]interface{}
	err := c.Bind(&payload)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	sess, _ := session.Get("name", c)
	sess.Values["name"] = payload["name"]
	err = sess.Save(c.Request(), c.Response())
	fmt.Println(err)
	return c.JSON(http.StatusOK, payload)
}

func main() {
	r, err := redistore.NewRediStore(10, "tcp", "localhost:6379", "", []byte("sessions"))

	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(session.Middleware(
		r,
	))
	e.GET("/", get)
	e.POST("/", post)
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
