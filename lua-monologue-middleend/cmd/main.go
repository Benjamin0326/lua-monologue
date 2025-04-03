package main

import (
	"fmt"
	"log"

	grpcclient "lua-monologue-middleend/internal/grpc"
	"lua-monologue-middleend/internal/handlers"
	"lua-monologue-middleend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	grpcclient.SendChatMessage("안녕 Lua Rust Server야 ㅎㅎ")

	router := gin.Default()
	router.Use(middleware.SetupCors())
	router.POST("/sendchatmessage", handlers.HandleChatMessagePost)

	port := 8080
	log.Printf("✅ 서버 실행 중: http://localhost:%d", port)
	router.Run(fmt.Sprintf(":%d", port))
}
