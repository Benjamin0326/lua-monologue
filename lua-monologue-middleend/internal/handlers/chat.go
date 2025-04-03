package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatMessages struct {
	Content string `json:"content"` // ✅ JSON에서 "content" 필드를 Go의 Content로 매핑
}

func HandleChatMessagePost(c *gin.Context) {
	//var chatMessages ChatMessages
	var req map[string]interface{} // JSON 데이터를 받을 변수

	fmt.Println("HandleChatMessagePost")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청 데이터"})
		return
	}

	fmt.Println("📌 받은 일기 내용:", req)
	c.JSON(http.StatusOK, gin.H{"message": "채팅 수신 완료", "data": req})
}
