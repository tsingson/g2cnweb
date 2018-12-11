package main

import (
	"fmt"
	"net"

	"github.com/tsingson/fastx/zaplogger"
	"github.com/tsingson/phi"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"go.uber.org/zap"
)

// FasthttpServ
// web 服务器配置与启动, 注意, 这里会尝试重用 http 监听端口, 以便可以运行多个实例
// 如果需要绑定 cpu 核心, 需要修改 main 以限制为单核运行
func StaticHttpServ(addr, wwwroot string, log *zap.Logger, zaplog *zaplogger.ZapLogger) {
	var (
		listener  net.Listener
		err, err1 error
		// 	fastLogger FastLogger

	)



	// 	var filesHandler = fasthttp.FSHandler("/var/www/files", 0)
	// router.Get("/", filesHandler )

	// reuse port
	listener, err = reuseport.Listen("tcp4", addr)
	if err != nil {
		log.Info("working in Microsoft Windows", zap.String("addr", addr))
		// for windows
		listener, err1 = net.Listen("tcp", addr)
		if err1 != nil {
			log.Fatal("Error", zap.Error(err1))
			panic("tcp connect error")
		}
	}

	// run fasthttp serv
	go func() {
		// fasthttp server setting here
		s := &fasthttp.Server{
			Handler:           fsHandler(wwwroot),
			Name:               ServerName,
			ReadBufferSize:     BufferSize,
			MaxConnsPerIP:      10,
			MaxRequestsPerConn: 100,
			// 	MaxRequestBodySize: 100<<20, // 100MB
			MaxRequestBodySize: 100<<20 , //1024 * 1024 * 4, // MaxRequestBodySize: 100<<20, // 100MB
			Concurrency:        MaxFttpConnect,
			DisableKeepalive:   false,
			Logger:             zaplog,
			// 	MaxKeepaliveDuration: time.Minute, // 新增限制
		}
		if err = s.Serve(listener); err != nil {
			log.Fatal("Error", zap.Error(err))
			panic("fasthttp running error")
		}
	}()
	log.Info("fasthttpgx server start success  ")
	fmt.Println("fasthttpgx server start success  ", addr)

}
//
func FasthttpServ(addr, webRoot string, log *zap.Logger, zaplog *zaplogger.ZapLogger) {
	var (
		listener  net.Listener
		err, err1 error
		// 	fastLogger FastLogger
		router *phi.Mux
	)

	// fastLogger.Logger = log
	router = StaticRouter(zaplog)
	router.NotFound(staticFsHandler(webRoot))

	// 	var filesHandler = fasthttp.FSHandler("/var/www/files", 0)
	// router.Get("/", filesHandler )

	// reuse port
	listener, err = reuseport.Listen("tcp4", addr)
	if err != nil {
		log.Info("working in Microsoft Windows", zap.String("addr", addr))
		// for windows
		listener, err1 = net.Listen("tcp", addr)
		if err1 != nil {
			log.Fatal("Error", zap.Error(err1))
			panic("tcp connect error")
		}
	}

	// run fasthttp serv
	go func() {
		// fasthttp server setting here
		s := &fasthttp.Server{
			Handler:            router.ServeFastHTTP,
			Name:               ServerName,
			ReadBufferSize:     BufferSize,
			MaxConnsPerIP:      10,
			MaxRequestsPerConn: 10,
			// 	MaxRequestBodySize: 100<<20, // 100MB
			MaxRequestBodySize: 1024 * 1024 * 4, // MaxRequestBodySize: 100<<20, // 100MB
			Concurrency:        MaxFttpConnect,
			DisableKeepalive:   false,
			Logger:             zaplog,
			// 	MaxKeepaliveDuration: time.Minute, // 新增限制
		}
		if err = s.Serve(listener); err != nil {
			log.Fatal("Error", zap.Error(err))
			panic("fasthttp running error")
		}
	}()
	log.Info("fasthttpgx server start success  ")
	fmt.Println("fasthttpgx server start success  ", addr)

}
