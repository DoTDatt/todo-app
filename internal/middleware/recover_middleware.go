package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Hệ thống gặp sự cố nghiêm trọng : %v", err)
				c.JSON(500, gin.H{"error": "Đã xảy ra lỗi nội bộ"})
				c.Abort()
			}
		}()
		c.Next()
	}
}
