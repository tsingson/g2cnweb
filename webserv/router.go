package main

import (
	"github.com/tsingson/fastx/middle"
	"github.com/tsingson/fastx/zaplogger"
	"github.com/tsingson/phi"
	"github.com/valyala/fasthttp"
)

// proxy router
func myRouter(webroot string, log *zaplogger.ZapLogger) *phi.Mux {
	r := phi.NewMux()
	//
	r.Use(log.FastHttpZapLogHandler)
	r.Use(middle.Recoverer)
	r.Get("/", redirectStatic)
	r.Get("/hello", helloHandler)
	r.Get("/phpmyadmin", redirectStatic)
	r.ServeFiles("/cn/*filepath", webroot)
	// r.NotFound( nullHandler)
	// r.MethodNotAllowed(nullHandler)

	// 	r.NotFound(StaticFsHandler(webroot))
	// 	r.Mount("/cn", staticRouter(webroot, log))

	return r
}

// redirectStatic
func redirectStatic(ctx *fasthttp.RequestCtx) {
	ctx.Redirect("/cn/index.html", fasthttp.StatusMovedPermanently)
	return
}