package main

import (
	"crypto/tls"
	"fmt"
	"net"

	"golang.org/x/crypto/acme/autocert"

	"github.com/oklog/run"
	"github.com/spf13/afero"
	"github.com/tsingson/fastx/utils"
	"github.com/tsingson/fastx/zaplogger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"go.uber.org/zap"
)

// FasthttpServ
// web 服务器配置与启动, 注意, 这里会尝试重用 http 监听端口, 以便可以运行多个实例
// 如果需要绑定 cpu 核心, 需要修改 main 以限制为单核运行
func StaticHttpServ(addr, wwwroot string, log *zap.Logger, zaplog *zaplogger.ZapLogger) {
	var (
		ln, lnTls net.Listener
		err, err1 error
		g         run.Group
		// 	fastLogger FastLogger

	)

	// 	var filesHandler = fasthttp.FSHandler("/var/www/files", 0)
	// router.Get("/", filesHandler )

	// reuse port
	ln, err = reuseport.Listen("tcp4", addr)
	if err != nil {
		log.Info("working in Microsoft Windows", zap.String("addr", addr))
		// for windows
		ln, err1 = net.Listen("tcp4", addr)
		if err1 != nil {
			log.Fatal("Error", zap.Error(err1))
			panic("tcp connect error")
		}
		log.Info("Listener success ")
	}
	ln.Close()

	r := myRouter(wwwroot, zaplog)

	s := &fasthttp.Server{
		Handler:            r.ServeFastHTTP, //    fsHandler(wwwroot),
		Name:               ServerName,
		ReadBufferSize:     BufferSize,
		MaxConnsPerIP:      10,
		MaxRequestsPerConn: 100,
		// 	MaxRequestBodySize: 100<<20, // 100MB
		MaxRequestBodySize: 100 << 20, // 1024 * 1024 * 4, // MaxRequestBodySize: 100<<20, // 100MB
		Concurrency:        MaxFttpConnect,
		DisableKeepalive:   false,
		Logger:             zaplog,
		// 	MaxKeepaliveDuration: time.Minute, // 新增限制
	}

	/**
	g.Add(func() error {
		return s.Serve(ln)
	}, func(error) {
		ln.Close()
	})
*/
	//
	//
	lnTls, err = reuseport.Listen("tcp4", ":https")
	if err != nil {
		log.Info("working in Microsoft Windows", zap.String("addr", ":https"))
		// for windows
		ln, err1 = net.Listen("tcp", ":https")
		if err1 != nil {
			log.Fatal("Error", zap.Error(err1))
			panic("tcp connect error")
		}
	}

	path, _ := utils.GetCurrentExecDir()
	afs := afero.NewOsFs()

	cachePath := path + "/cache"
	check, _ := afero.DirExists(afs, cachePath)
	if !check {
		e := afs.MkdirAll(cachePath, 0755)
		if e != nil {
			// TODO: handl err
		}
	}

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("g2cn.cn", "www.g2cn.cn"),
		Cache:      autocert.DirCache(cachePath),
	}

	cfg := m.TLSConfig()

	lnTlsCfg := tls.NewListener(lnTls, cfg)

	g.Add(func() error {
		return s.Serve(lnTlsCfg)
	}, func(error) {
		lnTlsCfg.Close()
	})

	if err = g.Run(); err != nil {
		log.Fatal("Error", zap.Error(err))
		panic("fasthttp running error")
	}

	log.Info("fasthttpgx server start success  ")
	fmt.Println("fasthttpgx server start success  ", addr)

}
