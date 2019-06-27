package middlelayer

/*
#include "middlelayerc.h"
#include <stdlib.h>
*/
import "C"

import (
	"crypto/md5"
	"eme/models"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// var host = "http://play.alleyes.com.cn"
// var path := "/webrtc/" + "1c66.m3u8"
// var key := "allEyes123Go"

var (
	ForceRtmpChan      = make(chan int)
	AddChannleChan     = make(chan int, 10)
	DelChannleChan     = make(chan int, 10)
	ChannelErrorChan   = make(chan map[int]error, 20)
	ChannelUrl         = make(map[int]string)
	ChannelSourceStart = make(chan string)
	ChangeRtmpChan     = make(chan string)
)
var (
	playHost         = beego.AppConfig.String("playHost")
	pushHost         = beego.AppConfig.String("pushHost")
	rtmpkey          = beego.AppConfig.String("rtmpkey")
	channelCount     = make(map[int]int)
	ClientsType      = make(map[int]int) //key 为channel id  value 1 >rtmp
	clientsTypeLock  sync.Mutex
	boxMac           string
	channelUrlLock   sync.Mutex
	channelCountLock sync.Mutex
)

func EnableChannelSourceCtl(mac string) {
	boxMac = mac
	go addChannelSource()
	go delChannelSource()
	go swChannelSource()

}

func addChannelSource() {
	for {
		id := <-AddChannleChan
		if channelCount[id] == 0 {
			playUrl, err := addWebrtcChannelData(id)
			if err != nil {
				tmpError := make(map[int]error)
				tmpError[id] = err
				ChannelErrorChan <- tmpError
				logs.Debug(fmt.Sprintf("... channen id  : %v stream start...", id))

			} else {
				if ClientsType[id] == 1 {
					channelUrlLock.Lock()
					ChannelUrl[id] = playUrl
					channelUrlLock.Unlock()
				}
			}
			channelCountLock.Lock()
			channelCount[id] = 1
			channelCountLock.Unlock()
			ChannelSourceStart <- playUrl

		} else {
			channelCountLock.Lock()
			channelCount[id]++
			channelCountLock.Unlock()
			ChannelSourceStart <- ChannelUrl[id]

		}
		logs.Error(fmt.Sprintf("...addChannelSource channen id is : %v,count is %v...", id, channelCount[id]))
		logs.Error(fmt.Sprintf("...url channen id is : %v,url all is %v...", id, ChannelUrl))

	}
}

func swChannelSource() {
	for {
		id := <-ForceRtmpChan
		playUrl := ""
		clientsTypeLock.Lock()
		ClientsType[id] = 1
		clientsTypeLock.Unlock()
		var err error
		if channelCount[id] > 0 {
			delWebrtcChannelData(id)
			time.Sleep(1 * time.Second)
			playUrl, err = addWebrtcChannelData(id)
			if err != nil {
				tmpError := make(map[int]error)
				tmpError[id] = err
				ChannelErrorChan <- tmpError
			} else {
				channelUrlLock.Lock()
				ChannelUrl[id] = playUrl
				channelUrlLock.Unlock()
			}
		} else {
			AddChannleChan <- id
			playUrl = <-ChannelSourceStart
		}
		logs.Error(fmt.Sprintf("...swChannelSource channen id is : %v,count is %v...", id, channelCount[id]))
		ChangeRtmpChan <- playUrl
	}
}
func delChannelSource() {
	for {
		id := <-DelChannleChan
		if channelCount[id] <= 1 {
			logs.Debug(fmt.Sprintf("... channen id  : %v stream stop...", id))
			delWebrtcChannelData(id)
			channelCountLock.Lock()
			delete(channelCount, id)
			delete(ChannelUrl, id)
			channelCountLock.Unlock()
			clientsTypeLock.Lock()
			ClientsType[id] = 0
			clientsTypeLock.Unlock()
		} else if channelCount[id] > 1 {
			channelCountLock.Lock()
			channelCount[id]--
			channelCountLock.Unlock()

		}
		logs.Error(fmt.Sprintf("...delChannelSource channen id is : %v,count is %v...", id, channelCount[id]))

	}

}

func addWebrtcChannelData(id int) (string, error) {
	channel, err := models.Get_Channel_ById(id)
	if err != nil {
		return "", err
	}
	if channel.Enable {
		size := channel.Image_width * channel.Image_height * 3 / 2
		var param C.GST_WEBRTC_Parameters
		param.id = C.int(channel.Id)
		if ClientsType[id] == 1 {
			param.width = C.int(channel.Image_width / 2)
		} else {
			param.width = C.int(channel.Image_width)
		}
		param.height = C.int(channel.Image_height)
		param.size = C.int(size)
		param.fps = C.int(25)
		if models.EngineConfig.Use_Yuv {
			param.format = C.GST_WEBRTC_FRAME_FORMAT_I420
		}
		// param.format = C.int( C.GST_VIDEO_WRITER_FRAME_FORMAT_I420)
		url := C.CString("")
		urlAdd := ""

		push, rtmp, flv, m3u8 := getUrl(boxMac, strconv.Itoa(id), rtmpkey)
		url = C.CString(push)
		m := map[string]string{
			"ForceRtmp": strconv.Itoa(ClientsType[id]),
			"rtmp":      rtmp,
			"flv":       flv,
			"m3u8":      m3u8,
		}
		mbytes, errjson := json.Marshal(m)
		if errjson != nil {
			return "", errjson
		}
		urlAdd = string(mbytes)
		logs.Error(".....url.........", push)
		logs.Error("...urlAdd...........", urlAdd)
		d := C.MyAddWebrtcChannel(&param, url, C.int(ClientsType[id]))
		fmt.Println("webrtc.......%d", d)
		return urlAdd, nil
	} else {
		return "", errors.New("Channel is no enable")
	}
}
func delWebrtcChannelData(id int) {
	C.MyDelWebrtcChannel(C.int(id))
}

func getUrl(mac, id, key string) (push, rtmp, flv, m3u8 string) {
	rtmpPath := "/" + mac + "/" + id
	flvPath := rtmpPath + ".flv"
	m3u8Path := rtmpPath + ".m3u8"
	rtmpPushHost := "rtmp://" + pushHost
	rtmpPlayHost := "rtmp://" + playHost
	httpHost := "http://" + playHost
	push = authKey(rtmpPushHost, rtmpPath, key)
	rtmp = authKey(rtmpPlayHost, rtmpPath, key)
	flv = authKey(httpHost, flvPath, key)
	m3u8 = authKey(httpHost, m3u8Path, key)
	return

}

func md5sum(src string) string {
	data := []byte(src)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has) //将[]byte转成16进制
}

func authKey(host, path, key string) string {
	rand := "0"
	uid := "0"
	exp := time.Now().Unix() + 1*3600
	sstring := fmt.Sprintf("%s-%v-%s-%s-%s", path, exp, rand, uid, key) //将[]byte转成16进制
	hashvalue := md5sum(sstring)
	auth := fmt.Sprintf("%v-%s-%s-%s", exp, rand, uid, hashvalue)
	return fmt.Sprintf("%s%s?auth_key=%s", host, path, auth)

}
