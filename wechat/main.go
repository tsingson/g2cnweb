package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
	"unsafe"

	"github.com/spf13/afero"
	"github.com/tsingson/fastx/fasturl"
	"github.com/tsingson/fastx/utils"
	"github.com/tsingson/fastx/zaplogger"
	"github.com/valyala/fastjson"

	"github.com/allegro/bigcache"
	"github.com/sanity-io/litter"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
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

func main() {
	var (
		accessToken string
		expiresIn   int
	)
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wxff5d2a2a34edb5d4&secret=21ade2a4e5ae100864b463dcfe2d9f93"

	client := &fasthttp.Client{}
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	// defer fasthttp.ReleaseRequest(request)
	// 	defer fasthttp.ReleaseResponse(response)

	// 	request.SetConnectionClose()
	request.SetRequestURI(url)
	request.Header.Add("Accept", "application/json")
	request.MultipartForm()
	request.RemoveMultipartFormFiles()

	err := client.DoTimeout(request, response, 15*time.Second)

	if err != nil {
		log.Info("get token err", zap.Error(err))
	}
	if response.StatusCode() == 200 {
		var p fastjson.Parser
		v, err := p.ParseBytes(response.Body())
		if err != nil {
			log.Fatal("json parse err", zap.Error(err))
		}
		accessToken = b2s(v.GetStringBytes("access_token"))
		expiresIn = v.GetInt("expires_in")
		litter.Dump(accessToken)
		litter.Dump(expiresIn)
	}
	fmt.Println(" ")
	litter.Dump(response.StatusCode())
	litter.Dump(string(response.Body()))

	getIp := "https://api.weixin.qq.com/cgi-bin/getcallbackip?access_token="
	utl2 := utils.StrBuilder(getIp, accessToken)
	resp, err1 := fasturl.FastGet(utl2, 5*time.Second)
	if err1 == nil {
		litter.Dump(b2s(resp.Body()))
	}
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}


func postFile(filename string, targetUrl string) error {

	// multipart.WriteField
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	// iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}

func postFileUpload(url, name string ) error {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		part, err := m.CreateFormFile("myFile", "foo.txt")
		if err != nil {
			return
		}
		file, err := os.Open(name)
		if err != nil {
			return
		}
		defer file.Close()
		if _, err = io.Copy(part, file); err != nil {
			return
		}
	}()
	http.Post(url, m.FormDataContentType(), r)
}