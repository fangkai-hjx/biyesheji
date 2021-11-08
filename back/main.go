package main

import (
	"context"
	"fmt"
	"gee/config"
	"gee/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	server := &http.Server{
		Addr:    config.ProjectConfig.ProjectHost + ":" + config.ProjectConfig.ProjectPort,
		Handler: router.Router,
	}
	fmt.Println("启动端口：", config.ProjectConfig.ProjectHost+":"+config.ProjectConfig.ProjectPort)
	// 在另一个 goroutine 中执行
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %s\n", err)
		}
	}()

	// 处理中断
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server")

	// 关闭连接
	//util.DB.Close()
	//log.Println("Close the DB Connection")
	//util.DockerClient.Close()
	//log.Println("Close the DockerClient Connection")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %s\n", err)
	}
	log.Println("Server exited")
}
