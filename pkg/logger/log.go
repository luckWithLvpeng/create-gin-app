package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"create-gin-app/pkg/config"

	"github.com/gin-gonic/gin"
)

var (
	lastDay       int
	logPath       string
	errorFileName string
	days          float64
	errFile       *os.File
)

func init() {
	logPath = config.DefaultConfig.Section("log").Key("log_folder").MustString("logFile")
	errorFileName = config.DefaultConfig.Section("log").Key("log_error").MustString("error.log")
	days = config.DefaultConfig.Section("log").Key("days").MustFloat64(7)
	// 创建日志文件夹
	_, err := os.Stat(logPath)
	if err != nil {
		os.Mkdir(logPath, 0755)
	}
	errFile, _ = os.OpenFile(fmt.Sprintf("%s/%s", logPath, errorFileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// 配置 log 格式
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	updateWriter()
}

// updateWriter 更新 writer
func updateWriter() {
	nowDay := time.Now().Second()
	if nowDay != lastDay {
		var file *os.File
		filename := time.Now().Format("2006-01-02")
		logFile := fmt.Sprintf("%s/%s.log", logPath, filename)

		file, _ = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if file != nil {
			gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
			gin.DefaultErrorWriter = io.MultiWriter(os.Stderr, errFile, file)
			log.SetOutput(gin.DefaultWriter)
		}
	}
	// 删除 n 天前的日志
	deleteLogDaysAgo(days)

	// 每天零点增加日志文件
	lastDay = nowDay
	now := time.Now()
	tmp := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dayend := tmp.AddDate(0, 0, 1)
	time.AfterFunc(dayend.Sub(now), updateWriter)

}

func deleteLogDaysAgo(n float64) {
	files := getLogFromFolder(logPath)
	now := time.Now()
	for _, v := range files {
		file := path.Base(v)
		fileSuffix := path.Ext(file)
		filename := strings.TrimSuffix(file, fileSuffix)
		timeArr := strings.Split(filename, "-")
		year, _ := strconv.Atoi(timeArr[0])
		month, _ := strconv.Atoi(timeArr[1])
		day, _ := strconv.Atoi(timeArr[2])
		tmp := time.Date(year, time.Month(month), day, 0, 0, 0, 0, now.Location())
		if now.Sub(tmp).Hours() > n*24 {
			os.Remove(v)
		}
	}
}

//获取目录下的log文件
func getLogFromFolder(path string) (files []string) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return files
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if !fi.IsDir() { // 递归目录下的日志文件
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".log")
			if ok && fi.Name() != errorFileName {
				files = append(files, path+PthSep+fi.Name())
			}
		}
	}
	return files
}
