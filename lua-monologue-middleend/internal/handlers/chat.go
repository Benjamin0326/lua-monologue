package handlers

import (
	"fmt"
	"net/http"

	grpcclient "lua-monologue-middleend/internal/grpc"
	"lua-monologue-middleend/internal/llmclient"

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

	// "content" 키가 있는지 확인
	contentRaw, ok := req["content"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content 필드가 없습니다"})
		return
	}

	// 타입이 string인지 확인
	contentStr, ok := contentRaw.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content 필드는 문자열이어야 합니다"})
		return
	}

	// LLM 통신
	response, err := llmclient.CallLLM(contentStr)
	if err != nil {
		fmt.Println("❌ 오류:", err)
		return
	}

	fmt.Println("🧠 LLM 응답:", response)

	fmt.Println("📌 받은 일기 내용:", contentStr)
	c.JSON(http.StatusOK, gin.H{"message": "채팅 수신 완료", "data": response})
	grpcclient.SendChatMessage(contentStr)
}
