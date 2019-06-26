package main

import (
	"eme/middleware/auth"
	_ "eme/models"
	"eme/pkg/config"
	_ "eme/pkg/logger"
	"eme/routers"
	"log"
	"syscall"

	"fmt"
	"time"

	"github.com/fvbock/endless"
)

// @title 电子物证
// @version 1.0
// @BasePath /v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// 编译启动后，可以通过发送 SIGTERM 信号，可以实现服务优雅的重启
	// 执行 kill -1 pid 服务会等待旧的请求处理完成，同时开启一个新的进程去处理新来的请求, 可以在客户毫不知情的情况下，实现优雅的热更新服务。
	// 执行 kill -9 pid 会真正的杀死这个服务
	endless.DefaultReadTimeOut = time.Duration(config.ReadTimeout) * time.Second
	endless.DefaultWriteTimeOut = time.Duration(config.WriteTimeout) * time.Second
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", config.HTTPPort)
	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	fmt.Printf("listening and serve on :%d\n ", config.HTTPPort)
	defer func() {
		//在系统宕机或者重启时，把权限黑名单保存到本地，系统启动时，从该文件中重新加载权限黑名单
		auth.BlackList.SaveFile("blackList.db")
	}()
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
