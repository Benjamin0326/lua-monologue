package handlers

import (
	"log"
	"net/http"

	"lua-monologue-middleend/internal/db"
	"lua-monologue-middleend/internal/middleware"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleRegister(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청"})
		return
	}

	log.Printf(req.Username)
	log.Printf(req.Password)

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "비밀번호 해싱 실패"})
		return
	}

	_, err = db.DB.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", req.Username, string(hash))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB 저장 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "회원가입 성공"})

}

func HandleLogin(c *gin.Context) {
	var req UserRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var hash string
	err := db.DB.QueryRow("SELECT password_hash FROM users WHERE username = $1", req.Username).Scan(&hash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "사용자 없음"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "사용자 없음 또는 비밀번호 불일치"})
		return
	}

	accessToken, err := middleware.GenerateJWT(req.Username, 15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "엑세스 토큰 생성 실패"})
		return
	}

	refreshToken, err := middleware.GenerateJWT(req.Username, 60*24*30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "리프레시 토큰 생성 실패"})
		return
	}

	c.SetCookie("refresh_token", refreshToken, 3600*24*30, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
	})
}

func HandleLogout(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "로그아웃 완료"})
}
