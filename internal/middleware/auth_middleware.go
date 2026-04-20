package middleware

import (
	"net/http"
	"strings"

	"github.com/DoDtatt/todo-app/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Yêu cầu đăng nhập "})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Định dạng Token không hợp lệ"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token hết hạn hoặc không hợp lệ"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {

			c.Set("user_id", int(claims["user_id"].(float64)))
			c.Set("role_id", int(claims["role_id"].(float64)))
		}

		c.Next()
	}
}
