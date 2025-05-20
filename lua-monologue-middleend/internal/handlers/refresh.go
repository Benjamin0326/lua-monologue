package handlers

import (
	"lua-monologue-middleend/internal/middleware"
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token 없음"})
		return
	}

	claims := &middleware.Claims{}

	tkn, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
		return middleware.JwtKey, nil
	})

	fmt.Println("✅ Refresh claims.Username:", claims.Username)

	if err != nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token 유효하지 않음"})
		c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
		return
	}

	newAccessToken, err := middleware.GenerateJWT(claims.Username, 15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Access Token 생성 실패"})
		c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
