package main

import (
	"time"
)

const (
	ServerName            = "EPG-cache-service"
	Version               = "0.1.1-20180418"
	LogFileNameTimeFormat = "2006-01-02-15"
	MaxFttpConnect        = 30000
	BufferSize            = 1024 * 2
	LiveApk               = "com.cts.live"
	LiveEpg               = "50.7.118.186"
	YuiApk                = "com.smalls.wonderfulyueplus"
	YuiEpg                = "23.237.56.42"
	ViewApk               = "com.smalls.newvideotwo"
	ViewEpg               = "50.7.118.186"
	SportsTV              = "com.smalls.sports"
	// for Adult APK only
	ActiveUrl  = "http://23.239.118.10/rpc/active"
	AuthUrl    = "http://23.239.118.10/rpc/auth"
	AdultEpgIp = "23.239.118.10"
	ServerPort = ":8095"
	VodGslb    = "192.154.108.2:45000"
	LiveGslb   = "192.154.108.2:43000"
	// aaaSession setting
	DefaultExpiration = time.Hour * 48
	ActiveExpiration = time.Minute * 5
	AuthExpiration = time.Minute * 15
	CleanupInterval   = time.Hour * 1
	//

	ExpiredDate = "2019-01-01 00:00:00"
)

/**
// config for adult system
直播 gslb：192.154.108.2:43000
点播 gslb：192.154.108.2:45000
3a 地址：23.239.118.10:8095
epg 地址:23.239.118.10


// apkType config
粤+精彩   com.smalls.wonderfulyueplus
缤纷视界   com.smalls.newvideotwo
adult apk    com.smalls.redshoes
*/
