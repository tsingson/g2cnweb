package main

import (
	"github.com/tsingson/fastx/zaplogger"
	"github.com/tsingson/phi"
)

// proxy router
func StaticRouter(log *zaplogger.ZapLogger) *phi.Mux {
	router := phi.NewMux()

	router.Use(log.FastHttpZapLogHandler)

	//
	// router.Use(middle.Recoverer)



	return router
}
