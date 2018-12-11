//  2018-02-17 02:49
package main

import (
	"github.com/allegro/bigcache"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/afero"
	"github.com/tsingson/fastx/utils"
	"github.com/tsingson/fastx/zaplogger"
	"go.uber.org/zap"
	"runtime"
)

const (
	LogFileNamePrefix = "g2cn-cn"
	PidFileName = "pid-webserv"
	WebRoot     = "/home/www/www"
	WebPort     = ":80"
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
		currentPath = path

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

		log = zaplogger.NewZapLog(logPath, LogFileNamePrefix, true )

		atom := zap.NewAtomicLevel()
		atom.SetLevel(zap.InfoLevel)

		zaplog = zaplogger.InitZapLogger(log)
		log.Info("- - - - - - - - - - - - - - -")
		log.Info("log init success")

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
func localWeb() {
	//

	// gops tracing
	// 	if err := agent.Listen(agent.Options{ConfigDir: currentPath}); err != nil {
	// 	log.Fatal("google gops Init Fail")
	// 	}

	log.Info("- - - - - - - - - - - - - - -")
	log.Info("daemon started")

	// 	middle.Log = log
	// FasthttpServ(config.AaaConfig.ServerPort, log)
	// 	FasthttpServ(":8000", "/Users/qinshen/git/linksmart/bin",  log, zaplog)
	FasthttpServ(WebPort, WebRoot, log, zaplog)
	// InitHttpProxy()

	// Wait forever.
	select {}

}

//
func main() {
	//
	runtime.GOMAXPROCS(128)

	// gops tracing
	// 	if err := agent.Listen(agent.Options{ConfigDir: currentPath}); err != nil {
	// 	log.Fatal("google gops Init Fail")
	// 	}

	// daemon
	cntxt := &daemon.Context{
		PidFileName: PidFileName,
		PidFilePerm: 0644,
		LogFileName: path + "/log/g2cn-daemon.log",
		LogFilePerm: 0640,
		WorkDir:     path,
		Umask:       027,
		Args:        []string{"webserv"},
	}

	d, err1 := cntxt.Reborn()
	if err1 != nil {
		log.Fatal("cat's reborn ", zap.Error(err1))
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Info("- - - - - - - - - - - - - - -")
	log.Info("daemon started")

	// 	middle.Log = log
	// FasthttpServ(config.AaaConfig.ServerPort, log)
	// 	FasthttpServ(":8000", "/Users/qinshen/git/linksmart/bin",  log, zaplog)
	StaticHttpServ(WebPort, WebRoot, log, zaplog)
	// InitHttpProxy()

	// Wait forever.
	select {}

}

//
