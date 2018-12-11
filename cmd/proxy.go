package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tsingson/fastx/utils"
	"github.com/tsingson/fastx/zaplogger"
	"github.com/tsingson/phi"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

const (
	TIMEOUT = 10 * time.Second
)

// Let's assume we have this (simplified) definition for the reverse proxy upstreams configuration
type (
	Upstreams struct {
		Method      string   `json:"method"`
		RelativeURI []string `json:"uri"`
		Host        string   `json:"host"`
		Port        int      `json:"port"`
	}

	// HttpProxy is a struct that stores the information for each proxy
	HttpProxy struct {
		proxy       *fasthttp.HostClient
		redirectURI *fasthttp.URI
	}
)

var (
	upstreams = []Upstreams{
		{
			Method:      "GET",
			RelativeURI: []string{"/category", "/epg", "/tvod", "/tvodchannel", "/hotmovie", "/movie", "/channel", "/today", "/check",},
			Host:        "127.0.0.1", //  "198.245.49.120", // or localhost if not running with Docker
			Port:        3004,
		},
	}
)
// Basic handler for the proxy
func (h *HttpProxy) reverseProxyHandler(ctx *fasthttp.RequestCtx) {
	req := &ctx.Request
	resp := &ctx.Response

	rangeByte := req.Header.Peek("Range")
	xrange := "0"
	var key, keyrange string
	if len(rangeByte) > 0 {
		xrange = string(rangeByte)

	}
	key = utils.StrBuilder(ctx.URI().String(), xrange)

	keyrange = utils.StrBuilder(key, "range")

	// jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIiA6ICJ0ZXJtaW5hbF9zdGIiLCAidXNlcl9pZCIgOiAiNDQwZTQxZjctYTcwNi00ZTQxLWFmYTAtNzBhZDIwOGU2NjM5IiwgImV4cCIgOiAxNTQyOTgxOTk2fQ.kMQ2PGc3jxpzrHDkqtK591MBtwrOF61UO3s5t0vLX8I"
	// req.Header.Set("Authorization", utils.StrBuilder("Bearer ", jwtToken))
	req.Header.Add("Content-Type", "application/json")
	log.Info("URI", zap.String("URI", ctx.URI().String()))

	req.Header.Del("Connection")

	if entry, err := cache.Get(key); err == nil {

		resp.SetBody(entry)

		if cr, er := cache.Get(keyrange); er == nil {
			resp.Header.SetBytesV("Content-Range", cr)
		}
		// Transfer-Encoding →chunked
		resp.Header.Set("Content-Location", string(ctx.Path()))
		resp.Header.Set("Transfer-Encoding", "chunked")
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp.Header.Set("Content-Encoding", "gzip")
		return

	}

	// token handle
	// var jwtToken string
	tokenByte := req.Header.Peek("Authorization")
	if len(tokenByte) > 0 {
		// jwtToken = utils.BytesToStringUnsafe(tokenByte)
		// 	req.Header.SetBytesV("Authorization", tokenByte)
		// jwtToken = strings.Replace(token , "Bearer ", "", -1)

	} else {
		/**
		token, er := getToken(log)
		if er == nil {
			jwtToken = utils.StrBuilder("Bearer ", token)
			req.Header.Set("Authorization", jwtToken)
		} else {
			ctx.Error("NO TOKEN", 500)
			return
		}
		*/
		ctx.Error("NO TOKEN", 500)
	}

	// Here we do some things for Auth (check and add headers basically) [..]

	if err := h.proxy.DoTimeout(req, resp, TIMEOUT); err != nil {
		resp.SetStatusCode(fasthttp.StatusServiceUnavailable)
		fmt.Printf("error when proxying the request: %s", err)
		ctx.Error("NO TOKEN or Access ERROR", 500)
	} else {
		/**
		Content-Range →0-199/*
   Content-Type →application/json; charset=utf-8
		 */
		cr := resp.Header.Peek("Content-Range")
		if len(cr) > 0 {
			er := cache.Set(keyrange, cr)
			if er != nil {
				// TODO:  handler err
			}
		} else {
			er := cache.Set(keyrange, []byte("0-199/*"))
			if er != nil {
				// TODO:  handler err
			}
		}
		cache.Set(key, resp.Body())
	}
	// resp.Header.Del("Connection")
	return
}

// proxy router
func ProxyRouter(log *zaplogger.ZapLogger) *phi.Mux {
	router := phi.NewMux()

	router.Use(log.FastHttpZapLogHandler)

	//
	// router.Use(middle.Recoverer)

	// Let's add all the upstreams to our fasthttp router
	for _, upstream := range upstreams {
		//
		for _, uri := range upstream.RelativeURI {
			newProxy := fasthttp.HostClient{
				IsTLS: false,
				Addr:  upstream.Host + ":" + strconv.Itoa(upstream.Port),
				// ReadTimeout: 60, // 如果在生产环境启用会出现多次请求现象
				MaxConns: 100,

				// Keep-alive connections are closed after this duration.
				//
				// By default connection duration is unlimited.
				MaxConnDuration: 60 * 10 * time.Second,

				// Idle keep-alive connections are closed after this duration.
				//
				// By default idle connections are closed
				// after DefaultMaxIdleConnDuration.
				MaxIdleConnDuration: 60 * 10 * time.Second,
			}
			proxyURI := fasthttp.URI{}
			proxyURI.SetPath(uri)

			httpProxy := &HttpProxy{
				proxy:       &newProxy,
				redirectURI: &proxyURI,
			}
			//
			router.Method(upstream.Method, uri, httpProxy.reverseProxyHandler)
		}

		//
	}

	return router
}

