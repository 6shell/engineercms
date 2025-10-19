package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"net/http/pprof"
)

type ProfController struct {
	web.Controller
}

func (c *ProfController) Get() {
	switch c.Ctx.Input.Param(":app") {
	default:
		pprof.Index(c.Ctx.ResponseWriter, c.Ctx.Request)
	case "":
		pprof.Index(c.Ctx.ResponseWriter, c.Ctx.Request)
	case "cmdline":
		pprof.Cmdline(c.Ctx.ResponseWriter, c.Ctx.Request)
	case "profile":
		pprof.Profile(c.Ctx.ResponseWriter, c.Ctx.Request)
	case "symbol":
		pprof.Symbol(c.Ctx.ResponseWriter, c.Ctx.Request)
	}
	c.Ctx.ResponseWriter.WriteHeader(200)
}
