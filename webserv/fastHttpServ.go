package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"golang.org/x/crypto/acme/autocert"

	"github.com/oklog/run"
	"github.com/spf13/afero"
	"github.com/tsingson/fastx/utils"
	"github.com/tsingson/fastx/zaplogger"
	"github.com/tsingson/phi"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"go.uber.org/zap"
)

// TlsServ
func TlsServ(addr, wwwroot string, log *zap.Logger, zaplog *zaplogger.ZapLogger) {

	keyFile := "./cert/cn.key"
	certFile := "./cert/cn.pem"

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

	if err := s.ListenAndServeTLS(":443", certFile, keyFile); err != nil {
		panic(err)
	}

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
		ln, lnTls net.Listener
		err, err1 error
		// 	fastLogger FastLogger
		router *phi.Mux
	)

	// fastLogger.Logger = log
	router = myRouter(webRoot, zaplog)
	router.NotFound(StaticFsHandler(webRoot))

	// 	var filesHandler = fasthttp.FSHandler("/var/www/files", 0)
	// router.Get("/", filesHandler )

	// reuse port
	ln, err = reuseport.Listen("tcp4", ":http")
	if err != nil {
		log.Info("working in Microsoft Windows", zap.String("addr", ":http"))
		// for windows
		ln, err1 = net.Listen("tcp", ":http")
		if err1 != nil {
			log.Fatal("Error", zap.Error(err1))
			panic("tcp connect error")
		}
	}

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

	//  s.ServeTLS()

	var g run.Group
	g.Add(func() error {
		return s.Serve(ln)
	}, func(error) {
		ln.Close()
	})

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
