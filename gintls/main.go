package main

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/allegro/bigcache"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
	"github.com/tsingson/fastx/ginzap"
	"github.com/tsingson/fastx/utils"
	"github.com/tsingson/fastx/zaplogger"
	"github.com/valyala/fasthttp/reuseport"
	"go.uber.org/zap"

	"golang.org/x/crypto/acme/autocert"
)

var (
	// 	output = log.New(os.Stdout, "", 0)
	log    *zap.Logger
	zaplog *zaplogger.ZapLogger
	cache  *bigcache.BigCache

	path, currentPath string

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

		log = zaplogger.NewZapLog(logPath, "gintls", true)

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

//

func main() {
	var (
		ln, lnTls net.Listener
		err, err1 error
	)

	wwwroot := path + "/g2cncn/public"
	log.Info("current exec path", zap.String("path", wwwroot))

	r := gin.New()
	// Zap logger

	// Add middleware to Gin, requires sync duration & zap pointer
	r.Use(ginzap.Logger(log))
	// Other gin configs
	r.Use(gin.Recovery())

	r.Use(static.Serve("/", static.LocalFile(wwwroot, false)))

	// Ping handler
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

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

	//
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
	//
	//
	lnTls, err = reuseport.Listen("tcp4", ":https")
	if err != nil {
		log.Info("working in Microsoft Windows", zap.String("addr", ":https"))
		// for windows
		lnTls, err1 = net.Listen("tcp", ":https")
		if err1 != nil {
			log.Fatal("Error", zap.Error(err1))
			panic("tcp connect error")
		}
	}

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("g2cn.cn", "www.g2cn.cn"),
		Cache:      autocert.DirCache(cachePath),
	}
	myTls := &tls.Config{GetCertificate: m.GetCertificate}

	// 	cfg := m.TLSConfig()

	// lnTlsCfg := tls.NewListener(lnTls, cfg)

	s := &http.Server{
		Addr:      ":https",
		TLSConfig: myTls,
		Handler:   r,
	}

	go http.Serve(ln, m.HTTPHandler(nil))

	go s.ServeTLS(lnTls, "", "")

	select {}
	// log.Fatal(autotls.RunWithManager(r, &m))
}
