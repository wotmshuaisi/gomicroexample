package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/gorilla/sessions"

	"github.com/boj/redistore"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

type handler struct {
}

func (h *handler) SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("_user_session", c)
		fmt.Println(sess)
		c.Set("session", sess.Values)
		next(c)
		return nil
	}
}

func (h *handler) get(c echo.Context) error {
	// sess, err := session.Get("_user_session", c)
	// fmt.Println(sess.Values["id"])
	// fmt.Println(reflect.TypeOf(sess.Values["id"]))

	sess := c.Get("session").(map[interface{}]interface{})

	if sess["name"] == nil {
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
	return c.HTML(200, sess["name"].(string))
}

func (h *handler) post(c echo.Context) error {
	var payload map[string]interface{}
	err := c.Bind(&payload)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	sess, _ := session.Get("_user_session", c)
	sess.Values["name"] = payload["name"]
	sess.Values["id"] = 111
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: false,
	}
	err = sess.Save(c.Request(), c.Response())
	return c.JSON(http.StatusOK, payload)
}

func main() {

	r, err := redistore.NewRediStore(10, "tcp", "127.0.0.1:6379", "", []byte("sessions"))

	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(session.Middleware(
		r,
	))
	h := handler{}
	e.Use(h.SessionMiddleware)
	e.GET("/", h.get)
	e.POST("/", h.post)
	fmt.Println("1")
	go http.ListenAndServe("localhost:6060", nil)
	if err := http.ListenAndServe(":8080", e); err != nil {
		log.Fatal(err)
	}
}
