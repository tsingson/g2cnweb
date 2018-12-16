package main

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/tsingson/fastx/utils"
	"github.com/tsingson/fastx/zaplogger"
	"github.com/tsingson/phi"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"net"
	"net/http"
)

// FasthttpServ
// web 服务器配置与启动, 注意, 这里会尝试重用 http 监听端口, 以便可以运行多个实例
// 如果需要绑定 cpu 核心, 需要修改 main 以限制为单核运行
func StaticHttpTlsServ(addr, wwwroot string, log *zap.Logger, zaplog *zaplogger.ZapLogger) {
	var (
		listener, tls net.Listener
		err, err1     error
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
	/**
		magic := certmagic.New(certmagic.Config{
			CA:     certmagic.LetsEncryptStagingCA,
			Email:  "tsingson@gmail.com",
			Agreed: true,
			// plus any other customization you want
		})

		// this obtains certificates or renews them if necessary
		err = magic.Manage([]string{"g2cn.cn"})
		if err != nil {
			log.Fatal("Error", zap.Error(err))
			panic("tcp connect error")
		}

		// to use its certificates and solve the TLS-ALPN challenge,
		// you can get a TLS config to use in a TLS listener!
		// tlsConfig := magic.TLSConfig()

		tls, err = certmagic.Listen([]string{"g2cn.cn"})

		if err != nil {
			log.Fatal("Error", zap.Error(err))
			panic("tcp connect error")
		}
	*/

	// run fasthttp serv
	go func() {
		// fasthttp server setting here
		s := &fasthttp.Server{
			Handler:            fsHandler(wwwroot),
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
		if err = s.Serve(listener); err != nil {
			log.Fatal("Error", zap.Error(err))
			panic("fasthttp running error")
		}
	}()
	//
	//

	path, _ := utils.GetCurrentExecDir()
	secretDir := utils.StrBuilder(path, "/cert")

	afs := afero.NewOsFs()
	check, _ := afero.DirExists(afs, secretDir)
	if !check {
		e := afs.MkdirAll(secretDir, 0755)
		if e != nil {
			// TODO: handl err
		}
	}
	m := autocert.Manager{
		Cache:      autocert.DirCache(secretDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("g2cn.cn", "www.g2cn.cn"),
	}

	tls = m.Listener()

	// run fasthttp serv
	go func() {
		// fasthttp server setting here
		s := &fasthttp.Server{
			Handler:            fsHandler(wwwroot),
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
		if err = s.Serve(tls); err != nil {
			log.Fatal("Error", zap.Error(err))
			panic("fasthttp running error")
		}
	}()
	log.Info("fasthttpgx server start success  ")
	fmt.Println("fasthttpgx server start success  ", addr)

}

func httpRedirectHandler(ctx *fasthttp.RequestCtx) {
	toURL := "https://"

	// since we redirect to the standard HTTPS port, we
	// do not need to include it in the redirect URL
	requestHost, _, err := net.SplitHostPort(string(ctx.Host()))
	if err != nil {
		requestHost = string(ctx.Host()) // host probably did not contain a port
	}

	toURL += requestHost
	toURL += string(ctx.RequestURI())

	// get rid of this disgusting unencrypted HTTP connection 🤢
	ctx.Response.Header.Set("Connection", "close")
	ctx.Redirect(toURL, http.StatusMovedPermanently)
	// 	http.Redirect(w, r, toURL, http.StatusMovedPermanently)
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
