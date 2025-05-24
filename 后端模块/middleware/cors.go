package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"taskmanager/config"
)

// CORSMiddleware 处理跨域请求的中间件
func CORSMiddleware() gin.HandlerFunc {
	// 获取配置
	appConfig := config.GetConfig()
	allowedOrigins := appConfig.CORSAllowedOrigins

	return func(c *gin.Context) {
		// 获取请求的Origin
		origin := c.Request.Header.Get("Origin")

		// 检查Origin是否在允许列表中
		allowOrigin := "*"
		if origin != "" {
			for _, allowed := range allowedOrigins {
				if allowed == "*" || strings.EqualFold(allowed, origin) {
					allowOrigin = origin
					break
				}
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
