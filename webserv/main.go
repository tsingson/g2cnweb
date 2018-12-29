//  2018-02-17 02:49
package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/RussellLuo/timingwheel"
	"github.com/allegro/bigcache"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/afero"
	"github.com/tsingson/fastx/utils"
	"github.com/tsingson/fastx/zaplogger"
	"go.uber.org/zap"
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



//
func main() {
	//
	runtime.GOMAXPROCS(128)

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
	// daemon

	log.Info("- - - - - - - - - - - - - - -")
	log.Info("daemon started")

	// 	middle.Log = log
	// FasthttpServ(config.AaaConfig.ServerPort, log)
	// 	FasthttpServ(":8000", "/Users/qinshen/git/linksmart/bin",  log, zaplog)
	wwwroot := currentPath + "/g2cncn/public"
	log.Info("current exec path", zap.String("path", wwwroot))
	StaticHttpServ(WebPort, wwwroot, log, zaplog)
	// InitHttpProxy()

	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	exitC := make(chan time.Time, 1)
	tw.AfterFunc(time.Second, func() {
		fmt.Println("The timer fires")
		exitC <- time.Now()
	})

	<-exitC

	// Wait forever.
	select {}

}

//
