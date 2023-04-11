package domain

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type IpConfContext struct {
	Ctx       *context.Context
	AppCtx    *app.RequestContext
	ClinetCtx *ClientContext
}

type ClientContext struct {
	IP string `json:"ip"`
}

func BuildIPConfContext(c *context.Context, ctx *app.RequestContext) *IpConfContext {
	ipConfContext := &IpConfContext{
		Ctx:       c,
		AppCtx:    ctx,
		ClinetCtx: &ClientContext{},
	}
	return ipConfContext
}
