package logger

import (
	"eme/pkg/config"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	lastDay       int
	logPath       string
	errorFileName string
	errFile       *os.File
)

func init() {
	logPath = config.DefaultConfig.Section("log").Key("log_folder").MustString("logFile")
	errorFileName = config.DefaultConfig.Section("log").Key("log_error").MustString("error.log")
	// 创建日志目录
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

	// 每天零点增加日志文件
	lastDay = nowDay
	now := time.Now()
	tmp := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dayend := tmp.AddDate(0, 0, 1)
	time.AfterFunc(dayend.Sub(now), updateWriter)

}
