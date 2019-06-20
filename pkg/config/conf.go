package config

import (
	"fmt"
	"log"
	"os"

	"github.com/go-ini/ini"
)

// DefaultConfig 默认的访问配置文件入口
var (
	DefaultConfig *ini.File
	err           error
	RunMode       string
	HTTPPort      int
	ReadTimeout   int
	WriteTimeout  int
)

func init() {
	DefaultConfig, err = ini.Load("conf/app.conf")
	if err != nil {
		fmt.Printf("Fail to read file conf/app.conf,maybe format err: %v", err)
		os.Exit(1)
	}
	loadServer()
}

func loadServer() {
	sec := DefaultConfig.Section("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}
	RunMode = DefaultConfig.Section("").Key("mode").MustString("pro")

	HTTPPort = sec.Key("port").MustInt(80)
	ReadTimeout = sec.Key("read_timeout").MustInt(60)
	WriteTimeout = sec.Key("write_timeout").MustInt(60)
}
