package handler

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"

	"github.com/micro/go-micro/client"

	"github.com/micro/go-log"

	"github.com/labstack/echo"
	proto "github.com/wotmshuaisi/gomicroexample/basis/proto"
)

type basisHandler struct {
	C proto.SayService
}

func setBasisRouter(g *echo.Group, c client.Client) {
	h := &basisHandler{
		C: proto.SayServiceClient("go.micro.srv.basis", c),
	}
	// g.Use(middleware.GzipWithConfig(middleware.GzipConfig{
	// 	Level: 5,
	// }))
	g.Use(tracerMiddleware("go.micro.web.basis"))
	g.POST("/hello", h.Hello)
}

func (b *basisHandler) Hello(c echo.Context) error {
	span := c.Get(trackerspan).(opentracing.Span)
	ctx := c.Get(trackerctx).(context.Context)

	a := map[string]interface{}{}
	err := c.Bind(&a)
	if err != nil || a["name"].(string) == "" {
		log.Log(a)
		return c.HTML(400, err.Error())
	}
	span.SetTag("payload", a)

	r, err := b.C.Hello(ctx, &proto.Request{Name: a["name"].(string)})
	if err != nil {
		fmt.Print("err:")
		fmt.Println(err.Error())
		return c.HTML(500, err.Error())
	}
	return c.HTML(200, r.Msg)
}
