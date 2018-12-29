package main

import (
	"fmt"

	"github.com/allegro/bigcache"
	"github.com/mholt/certmagic"
	"github.com/spf13/afero"
	"github.com/tsingson/fastx/middle"
	"github.com/tsingson/fastx/utils"
	"github.com/tsingson/fastx/zaplogger"
	"github.com/tsingson/phi"
	"github.com/valyala/fasthttp"
	"github.com/valyala/tcplisten"
	"go.uber.org/zap"

	"github.com/sanity-io/litter"
)

const (
	LogFileNamePrefix = "g2cn-cn"
	PidFileName       = "pid-webserv"
	WebRoot           = "/home/www/www"
	WebPort           = ":80"
)

var (
	// 	output = log.New(os.Stdout, "", 0)
	log    *zap.Logger
	zaplog *zaplogger.ZapLogger
	cache  *bigcache.BigCache

	path, currentPath string
	domainName        string = "www.g2cn.cn"
	// apkEpgMap         sync.Map
)

func init() {
	afs := afero.NewOsFs()
	// get run path
	{
		path, _ = utils.GetCurrentExecDir()
		currentPath, _ = utils.GetCurrentPath()

	}
	{
		logPath := path + "/log"
		check, _ := afero.DirExists(afs, logPath)
		if !check {
			e := afs.MkdirAll(logPath, 0755)
			if e != nil {
				// TODO: handl err
			}
		}

		// log setup

		log = zaplogger.NewZapLog(logPath, LogFileNamePrefix, true)

		atom := zap.NewAtomicLevel()
		atom.SetLevel(zap.InfoLevel)

		zaplog = zaplogger.InitZapLogger(log)
		log.Info("- - - - - - - - - - - - - - -")
		log.Info("log init success")
		log.Info("path", zap.String("current", currentPath))

	}
	/**
		var err error
		cache, err = bigcache.NewBigCache(bigcache.Config{
			Shards:             128,              // number of shards (must be a power of 2)
			LifeWindow:         10 * time.Minute, // time after which entry can be evicted
			CleanWindow:        30 * time.Second,
			MaxEntriesInWindow: 100 * 60 * 10,   // rps * lifeWindow
			MaxEntrySize:       1024 * 1024 * 3, // max entry size in bytes, used only in initial memory allocation
			Verbose:            false,            // prints information about additional memory allocation
			HardMaxCacheSize:   8192,            // Mb
		})
		if err != nil {
			log.Fatal("cache Init Error", zap.Error(err))
			os.Exit(1)
		}
	*/

}

func main() {

	// 	var (
	// ln net.Listener
	// 	g run.Group
	// )
	path, _ := utils.GetCurrentExecDir()
	cachePath := path + "/tls"
	litter.Dump(cachePath)
	afs := afero.NewOsFs()
	check, _ := afero.DirExists(afs, cachePath)
	if !check {
		e := afs.MkdirAll(cachePath, 0755)
		if e != nil {
			// TODO: handl err
		}
	}
	//
	var tlsStorage = certmagic.FileStorage{Path: cachePath}
	cfg := certmagic.Config{
		CA:     certmagic.LetsEncryptProductionCA,
		Email:  "tsingson@me.com",
		Agreed: true,
		// plus any other customization you want
	}

	magic := certmagic.NewWithCache(certmagic.NewCache(&tlsStorage), cfg)

	// this obtains certificates or renews them if necessary

	err := magic.Manage([]string{domainName})
	if err != nil {

	}

	key := certmagic.KeyBuilder{}

	certFile := key.SiteCert(certmagic.LetsEncryptProductionCA, domainName, )
	keyFile := key.SitePrivateKey(certmagic.LetsEncryptProductionCA, domainName)
	litter.Dump(certFile)
	litter.Dump(keyFile)

	wwwroot := currentPath + "/g2cncn/public"

	r := phi.NewMux()
	//
	r.Use(zaplog.FastHttpZapLogHandler)
	r.Use(middle.Recoverer)
	r.Get("/hello", helloHandler)
	r.Get("/", redirectStatic)
	r.ServeFiles("/cn/*filepath", wwwroot)

	s := &fasthttp.Server{
		Handler: r.ServeFastHTTP, //    fsHandler(wwwroot),
		Name:    "tls",
		// ReadBufferSize:     1024 * 2,
		MaxConnsPerIP:      10,
		MaxRequestsPerConn: 100,
		// 	MaxRequestBodySize: 100<<20, // 100MB
		MaxRequestBodySize: 100 << 20, // 1024 * 1024 * 4, // MaxRequestBodySize: 100<<20, // 100MB
		// 	Concurrency:        MaxFttpConnect,
		DisableKeepalive: false,
		Logger:           zaplog,
		// 	MaxKeepaliveDuration: time.Minute, // 新增限制
	}

	lncfg := tcplisten.Config{
		ReusePort: true,
	}
	ln, err := lncfg.NewListener("tcp4", ":443")
	if err != nil {
		log.Panic("cannot listen to :443", zap.Error(err))
	}

	go func() {
		err = s.ServeTLS(ln, cachePath+"/"+certFile, cachePath+"/"+keyFile)
		if err != nil {
			log.Fatal("Error", zap.Error(err))
			log.Panic("fasthttp running error", zap.Error(err))
		}
	}()

	log.Info("fasthttpgx server start success  ")
	fmt.Println("fasthttpgx server start success  ", ":443")
	select {}
}

func helloHandler(ctx *fasthttp.RequestCtx) {
	// 	fmt.Fprintf(ctx, "Hello, world!\n\n")
	/**
	fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
	fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
	fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
	fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
	fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
	*/
	ctx.SetContentType("text/plain; charset=utf8")

	// fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	// fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
	// fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
	// fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

	// log.Info().Str("ip", ctx.RemoteIP().String()).Str("serv", ctx.LocalIP().String()).Msg("log in helloHandler")

	// Set arbitrary headers
	// 	ctx.Response.Header.Set("X-My-Header", "my-header-value")
	// Set cookies
	/**
	var c fasthttp.Cookie
	c.SetKey("cookie-name")
	c.SetValue("cookie-value")
	ctx.Response.Header.SetCookie(&c)
	*/
	return
}

// redirectStatic
func redirectStatic(ctx *fasthttp.RequestCtx) {
	ctx.Redirect("/cn/index.html", fasthttp.StatusMovedPermanently)
	return
}
