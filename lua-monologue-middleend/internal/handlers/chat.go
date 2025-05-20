package handlers

import (
	"fmt"
	"net/http"

	"lua-monologue-middleend/internal/db"
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

	userName := c.GetString("username")

	fmt.Println(userName)

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
	response, err := llmclient.CallLLM(userName, contentStr)
	if err != nil {
		fmt.Println("❌ 오류:", err)
		return
	}

	fmt.Println("🧠 LLM 응답:", response)

	fmt.Println("📌 받은 일기 내용:", contentStr)
	c.JSON(http.StatusOK, gin.H{"message": "채팅 수신 완료", "data": response})
	grpcclient.SendChatMessage(contentStr, "user", userName) // user or assistant and user-test

	grpcclient.SendChatMessage(response, "assistant", userName)
}

func HandleChatMessageGet(c *gin.Context) {
	userName := c.GetString("username")
	query_rows, err := db.DB.Query("SELECT user_id, role, content, created_at from messages WHERE user_id=$1 ORDER BY created_at ASC", userName)

	if err != nil {
		c.JSON(500, gin.H{"error": "Fail to query from DB"})
		return
	}

	defer query_rows.Close()

	type Message struct {
		UserID    string `json:"user_id"`
		Role      string `json:"role"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	}

	var messages []Message

	for query_rows.Next() {
		var msg Message
		if err := query_rows.Scan(&msg.UserID, &msg.Role, &msg.Content, &msg.CreatedAt); err != nil {
			c.JSON(500, gin.H{"error": "스캔 실패"})
			return
		}
		messages = append(messages, msg)
	}

	c.JSON(200, messages)

}
