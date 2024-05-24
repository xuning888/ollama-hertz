package main

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xuning888/yoyoyo/config"
	"github.com/xuning888/yoyoyo/internal/controller"
	"github.com/xuning888/yoyoyo/internal/dal"
	mredis "github.com/xuning888/yoyoyo/internal/dal/redis"
	"github.com/xuning888/yoyoyo/internal/service"
	"github.com/xuning888/yoyoyo/pkg/http"
	"github.com/xuning888/yoyoyo/pkg/lock"
	"github.com/xuning888/yoyoyo/pkg/logger"
	"time"
)

func init() {
	config.Init()
	dal.Init()
	logger.InitLogger()
	client, ok := mredis.Client.(redis.UniversalClient)
	if !ok {
		panic("init redis lock error")
	}
	lock.Init(client)
}

func Register(router *gin.Engine) {
	chatService := service.NewChatService()
	chatController := controller.NewChatController(chatService)
	chatController.Register(router)
}

func main() {
	router := gin.Default()
	// register html
	router.Static("/", "./static")

	// register router
	Register(router)

	server := http.NewServer(
		router,
		config.DefaultConfig.ServerPort,
		http.WithShutdownTimout(time.Second*time.Duration(25)),
	)
	server.Serve()
}
