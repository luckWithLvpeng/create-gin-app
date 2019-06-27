package middlelayer

/*
#include "middlelayerc.h"
*/
import "C"

import (
	"bufio"
	"eme/models"
	"eme/tools"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron"
)

// ChannelParam  开启通道保存视频的数据格式
// ID  通道ID
// Width  保存视频的宽度
// Height  保存视频的高度
type ChannelParam struct {
	ID     int
	Name   string
	Width  int
	Height int
	size   int
	fps    int
	format int
	path   string //录像保存路径
}

// UsbInfo ...
type UsbInfo struct {
	UsbMnt  bool   //是否插入U盘
	UsbCap  uint64 //U盘剩余空间
	UsbFull bool   //空间是否达到限制容量
	DelDate string //本次删除录像日期
	Err     string //写入录像错误信息
}

var (
	timeLocation      = beego.AppConfig.String("location")
	usbCapTopic       = beego.AppConfig.String("MQTT_TOPIC") + "/usb"
	videoDuration     = beego.AppConfig.String("videoDuration")
	channelsForRecord = make(map[int]*ChannelParam)
	timer             *time.Ticker
	usbCapLimit       int
	o                 sync.Once
	saveBufferSize    = 20
	saveBuffer        = make(chan ChannelParam, saveBufferSize)
	saveLocker        sync.Mutex //saveBuffer操作锁
	saveToUsbLocker   sync.Mutex //saveTousb锁

)

//初始化u盘容量限制值
func init() {
	capLimit, e := beego.AppConfig.Int("usbCapLimit")
	if e != nil {
		usbCapLimit = 2000000 //2G
	} else {
		usbCapLimit = capLimit
	}
}

// StartRecord 开始保存某一个通道的视频，传入ChannelParam类型指针
func StartRecord(data *ChannelParam) {
	data.size = data.Width * data.Height * 3 / 2
	data.fps = 12
	data.format = C.GST_VIDEO_WRITER_FRAME_FORMAT_I420
	if channelsForRecord[data.ID] == nil {
		go startRecordById(data.ID, data)
	}
	channelsForRecord[data.ID] = data
	o.Do(initOnce)

}

// StopRecord 停止保存某一个通道的视频
func StopRecord(channelID int) {
	go stopRecord(channelID)
}

//initOnce  记录视频任务
func initOnce() {
	go saveRecordToUsb()
	c := cron.New()
	c.AddFunc("0 0/"+videoDuration+" * * * *", func() { go startRecord() })
	c.Start()
}

//startRecord 循环channelsForRecord 多协程开启录制任务
func startRecord() {
	// now := time.Now()
	// logs.Debug(fmt.Sprintf("...System start new record at %v...",now))
	for channelID, channelParam := range channelsForRecord {
		go startRecordById(channelID, channelParam)
	}
}

//startRecordById  根据通道Id周期录像，判断有无u盘，存在u盘开始录像，首先保存上次录像文件，再开启新的录像
func startRecordById(channelID int, channelParam *ChannelParam) {
	logs.Debug(fmt.Sprintf("...ID: %v...start record...", channelID))
	if channelParam == nil {
		return
	}
	if !hasUsb() {
		usb := UsbInfo{}
		if data, jerr := json.Marshal(usb); jerr == nil {
			tools.PublishTopic(usbCapTopic, string(data))
		}
		return
	}

	//新建路径
	filePath := "/tmp/video/" + strconv.Itoa(channelID) + "_" + channelParam.Name + "/"
	local, _ := time.LoadLocation(timeLocation)
	now := time.Now().In(local)
	day := now.Format("20060102")
	t := now.Format("20060102150405")
	filePath += day + "/"
	videoPath := filePath + t + ".mp4"
	if !isExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return
		}
	}
	usbPath := strings.Replace(filePath, "tmp", "mnt/usb", 1)
	if !isExist(usbPath) {
		err := os.MkdirAll(usbPath, os.ModePerm)
		if err != nil {
			return
		}
	}
	//保存上次录像
	//logs.Debug(fmt.Sprintf("...ID: %v...video's this path is %v...",channelID,videoPath))
	if channelParam.path != "" {
		//logs.Debug(fmt.Sprintf("...ID: %v...video's last path is %v...",channelID,channelParam.path))
		saveRecord(channelID, *channelParam)
	}
	// 更新临时路径
	channelParam.path = videoPath
	var param C.GST_VIDEO_WRITER_Parameters
	param.id = C.int(channelParam.ID)
	param.width = C.int(channelParam.Width)
	param.height = C.int(channelParam.Height)
	param.size = C.int(channelParam.size)
	param.fps = C.int(channelParam.fps)
	param.format = C.int(channelParam.format)
	param.p_video_path = C.CString(videoPath)
	// 开始创建文件视频
	C.MyAddVideowriterChannel(&param)
}

//isExist  pathExists 判断路径是否存在
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//stopRecord 停止录像，同时删除录像map对应数据
func stopRecord(channelID int) {
	if channelsForRecord[channelID] != nil {
		saveRecord(channelID, *channelsForRecord[channelID])
	}
	delete(channelsForRecord, channelID)
}

//saveRecord 根据通道ID保存视频同时将params通过chan触发tmp>u盘move
func saveRecord(channelID int, channelParam ChannelParam) {
	if channelParam.path != "" {
		C.MyDelVideowriterChannel(C.int(channelID))
		go func() {
			saveLocker.Lock()
			defer saveLocker.Unlock()
			l := len(saveBuffer)
			//logs.Debug(fmt.Sprintf("...ID:.%v ... Put new chan to buffer buffer len is %v...",channelID,l))
			if l >= saveBufferSize {
				usb := GetUsbCap(true)
				usb.Err = "U盘写入数据太慢，请更换！"
				if data, jerr := json.Marshal(usb); jerr == nil {
					tools.PublishTopic(usbCapTopic, string(data))
				}
			} else {
				saveBuffer <- channelParam
			}
		}()
	}
}

//saveRecordToUsb 根据chan消息将tmp文件夹中录像移动到u盘
func saveRecordToUsb() {
	for {
		params := <-saveBuffer
		saveToUsbLocker.Lock()
		//logs.Debug(fmt.Sprintf("...%s...Moving record from tmp to usb...",params.path))
		usb := GetUsbCap(true)
		if usb.UsbMnt {
			tmpPath := params.path
			usbPath := strings.Replace(tmpPath, "tmp", "mnt/usb", 1)
			err := moveFile(tmpPath, usbPath)
			if err != nil {
				usb.Err = fmt.Sprintf("通道：%v,U盘保存视频时发生错误：%v", params.ID, err)
				logs.Error(fmt.Sprintf("......Move record from tmp to usb has err: %s id is %v\n", err, params.ID))
			}
		}
		filePath := "/tmp/video/" + strconv.Itoa(params.ID) + "_" + params.Name + "/"
		delTmpLastDir(filePath)
		if data, jerr := json.Marshal(usb); jerr == nil {
			tools.PublishTopic(usbCapTopic, string(data))
		}
		saveToUsbLocker.Unlock()
	}
}

//stopAllRecord 停止全部录像任务
func stopAllRecord() {
	for channelID, _ := range channelsForRecord {
		models.ToggleChannelRecorded(channelID)
		stopRecord(channelID)
	}
	filePath := "/tmp/video/"
	err := os.RemoveAll(filePath)
	if err != nil {
		logs.Error(fmt.Sprintf("When stop all record ,tmp remove all record has err: %s \n", err))
	}
}

//moveFile 从源文件复制到目标文件，然后删除源文件
func moveFile(originalPath string, newPath string) error {
	originalFile, err := os.Open(originalPath)
	if err != nil {
		return err
	}
	defer os.Remove(originalPath)
	defer originalFile.Close()
	newFile, err1 := os.Create(newPath)
	if err1 != nil {
		return err1
	}
	defer newFile.Close()
	_, err2 := io.Copy(newFile, originalFile)
	if err2 != nil {
		return err2
	}
	err = newFile.Sync()
	if err != nil {
		return err
	}
	return nil
}

//获取U盘数据
func GetUsbCap(delDo bool) UsbInfo {
	usb := UsbInfo{}
	if !hasUsb() {
		return usb
	}
	var stat syscall.Statfs_t
	syscall.Statfs("/mnt/usb", &stat)
	d := stat.Bavail * uint64(stat.Bsize) >> 10
	usb.UsbCap = d
	if d > 0 {
		usb.UsbCap = d
		if d < uint64(usbCapLimit) {
			usb.UsbFull = true
			if delDo {
				usb.DelDate = delUsbOldData()
			}
		}
		usb.UsbMnt = true
	} else {
		usb.UsbMnt = false
		logs.Error(fmt.Sprintf("USB ERROR: %s\n", err))
	}
	return usb

}

//delUsbOldData 删除u盘中最早日期文件夹或者当日最早那个视频文件，返回删除的日期字符串
func delUsbOldData() string {
	var oldDate = 99999999
	var delDate string
	var oldPath = make(map[int][]string)
	for channelID, channelParam := range channelsForRecord {
		filePath := "/mnt/usb/video/" + strconv.Itoa(channelID) + "_" + channelParam.Name + "/"
		files, err := ioutil.ReadDir(filePath)
		var fileNames []int
		if err == nil {
			for _, v := range files {
				match, _ := regexp.MatchString("^((19|20)[0-9]{2})(0[1-9]|1[012])(0[1-9]|[12][0-9]|3[01])$", v.Name())
				if !match {
					continue
				}
				d, e := strconv.Atoi(v.Name())
				if e == nil {
					fileNames = append(fileNames, d)
				}
			}
			if len(fileNames) > 0 {
				if fileNames[0] < oldDate {
					delete(oldPath, oldDate)
					oldDate = fileNames[0]
					oldStr := strconv.Itoa(fileNames[0])
					oldPath[oldDate] = append(oldPath[oldDate], filePath+oldStr)
				} else if fileNames[0] == oldDate {
					oldStr := strconv.Itoa(fileNames[0])
					oldPath[oldDate] = append(oldPath[oldDate], filePath+oldStr)
				}
			}

		}

	}
	if oldDate != 99999999 {
		local, _ := time.LoadLocation(timeLocation)
		now := time.Now().In(local)
		day := now.Format("20060102")
		delDate = strconv.Itoa(oldDate)
		if delDate != day {
			for _, v := range oldPath[oldDate] {
				e := os.RemoveAll(v)
				if e != nil {
					logs.Error(fmt.Sprintf("Delete record has error: %s\n", e))
				}
			}
		} else {
			for _, v := range oldPath[oldDate] {
				records, _ := ioutil.ReadDir(v)
				if len(records) > 0 && !records[0].IsDir() {
					delDate = records[0].Name()
					e := os.Remove(v + "/" + records[0].Name())
					if e != nil {
						logs.Error(fmt.Sprintf("Delete record has error: %s\n", e))
					}
				}
			}
		}
		oldDate = 99999999
	}
	return delDate
}
func hasUsb() bool {
	file, err := os.Open("/proc/diskstats")
	defer file.Close()
	if err != nil {
		logs.Error(fmt.Sprintf("Read /proc/diskstats has err: %s\n", err))
		return false
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}
		name := fields[2]
		match := strings.HasPrefix(name, "sd")
		if match {
			return true
		}
	}
	if err := scanner.Err(); err != nil {
		logs.Error(fmt.Errorf("scan error for /proc/diskstats: %s", err))
		return false
	}
	return false
}

//delTmpLastDir 删除临时文件夹中上一日未成功移动的视频文件
func delTmpLastDir(chpath string) {
	files, err := ioutil.ReadDir(chpath)
	if err == nil {
		for k, v := range files {
			if k != len(files)-1 && v.IsDir() {
				err := os.RemoveAll(chpath + "/" + v.Name())
				if err != nil {
					logs.Error(fmt.Sprintf("Remove tmp dir has  err: %s\n", err))
				}
			}
		}
	} else {
		logs.Error(fmt.Sprintf("Read tmp dir has  err: %s\n", err))
	}
}

//循环发送U盘状态数据
// func sendUsbCap() {
// 	for {
// 		if len(channelsForRecord) > 0 {
// 			usb := GetUsbCap(true)
// 			logs.Debug(fmt.Sprintf("USB: %v\n", usb))
// 			if data, jerr := json.Marshal(usb); jerr == nil {
// 				tools.PublishTopic(usbCapTopic, string(data))
// 			}
// 		}
// 		time.Sleep(time.Second * videoDuration*5)
// 	}
// }
//test 20190521   录像打开后出现挂掉现象
// func getGID() uint64 {
//     b := make([]byte, 64)
//     b = b[:runtime.Stack(b, false)]
//     b = bytes.TrimPrefix(b, []byte("goroutine "))
//     b = b[:bytes.IndexByte(b, ' ')]
//     n, _ := strconv.ParseUint(string(b), 10, 64)
//     return n
// }

// func DoTest(channelID int) {
// 	for {
// 		// addVideoLock.Lock()
// 		logs.Debug(fmt.Sprintf("...%v...start record...",channelID))
// 		logs.Debug(fmt.Sprintf("...%v...goroutine ID is...%v",channelID,getGID()))

// 		channelParam :=channelsForRecord[channelID]
// 		filePath := "/tmp/video/" + strconv.Itoa(channelID) + "_" + channelParam.Name + "/"
// 		local, _ := time.LoadLocation(timeLocation)
// 		now := time.Now().In(local)
// 		day := now.Format("20060102")
// 		t := now.Format("20060102150405")
// 		filePath += day + "/"
// 		videoPath := filePath + t + ".mp4"
// 		if !isExist(filePath) {
// 			// 递归创建目录 ，假如目录不存在
// 			err := os.MkdirAll(filePath, os.ModePerm)
// 			if err != nil {
// 				// 暂时停止创建新的视频
// 				logs.Debug(fmt.Sprintf("...%v... record file err...%v",channelID,err))
// 				continue
// 			}
// 		}
// 		//write
// 		var param C.GST_VIDEO_WRITER_Parameters
// 			param.id = C.int(channelParam.ID)
// 			param.width = C.int(channelParam.Width)
// 			param.height = C.int(channelParam.Height)
// 			param.size = C.int(channelParam.size)
// 			param.fps = C.int(channelParam.fps)
// 			param.format = C.int(channelParam.format)
// 			param.p_video_path = C.CString(videoPath)
// 		C.MyAddVideowriterChannel(&param)
// 		// addVideoLock.Unlock()
// 		time.Sleep(30*time.Second)
// 		logs.Debug(fmt.Sprintf("...%v...end record...",channelID))

// 		//sleep(1)
// 		// addVideoLock.Lock()
// 		C.MyDelVideowriterChannel(C.int(channelID))
// 		logs.Debug(fmt.Sprintf("...%v...save record...",channelID))

// 		// addVideoLock.Unlock()
// 		//save

// 	}
// }
