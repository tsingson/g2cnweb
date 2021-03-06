package main

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/tsingson/fastx/utils"
	"github.com/valyala/fasthttp"
)

// StaticFsHandler
func StaticFsHandler(webRoot string) func(ctx *fasthttp.RequestCtx) {
	fs := &fasthttp.FS{
		// Path to directory to serve.
		Root: webRoot, // "/var/www/static-site",

		// Generate index pages if client requests directory contents.
		GenerateIndexPages: false,
		IndexNames:         []string{"index.html"},
		// Enable transparent compression to save network traffic.
		Compress:        true,
		AcceptByteRange: true,
		CacheDuration:   90 * time.Second,
	}

	// func New(fs *fasthttp.FS) func(ctx *fasthttp.RequestCtx, next func(error)) {
	staticHandler := fs.NewRequestHandler()
	return func(ctx *fasthttp.RequestCtx) {
		//
		var (
			m, path string
			afs     afero.Fs
		)
		//
		m = string(ctx.Method())

		if m != "GET" && m != "HEAD" {
			ctx.Error("not Allow Method", 500)
			return
		}

		path = string(ctx.Path())
		afs = afero.NewOsFs()

		fileInfo, err := afs.Stat(utils.StrBuilder(fs.Root, path))
		// if err != nil && os.IsNotExist(err) {
		if err != nil {
			ctx.Error("not Found", 500)
			return
		}

		// An exist file
		// fasthttp.FS handle it
		if !fileInfo.IsDir() {
			staticHandler(ctx)
			return
		}
		staticHandler(ctx)
	}

}

// fsHandler
func fsHandler(wwwroot string) fasthttp.RequestHandler {
	fs := &fasthttp.FS{
		// Path to directory to serve.
		Root:       wwwroot,
		IndexNames: []string{"index.html"},
		// Generate index pages if client requests directory contents.
		GenerateIndexPages: false,
		PathNotFound:       pathNotFound,

		// Enable transparent compression to save network traffic.
		Compress:        true,
		AcceptByteRange: true,
		CacheDuration:   90 * time.Second,
	}

	// Create request handler for serving static files.
	return fs.NewRequestHandler()
}

// pathNotFound
func pathNotFound(ctx *fasthttp.RequestCtx) {
	ctx.Error("", 500)
	return
}

// filterPath
func filterPath(path string, index []string) error {
	for _, v := range index {
		_, err := os.Stat(path + v)
		if err == nil {
			return nil
		}
	}
	return errors.New("Not found index page in directory: " + path)
}

// design and code by tsingson
