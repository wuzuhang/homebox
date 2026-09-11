package middleware

import (
	"net/http"
	"strings"

	"home-box/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuthJWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "请求未携带 Token，无访问权限",
			})
			c.Abort() // 阻止后续 Handler 执行
			return
		}

		// 按空格分割，标准格式为: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Token 格式错误，应为 Bearer <token>",
			})
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Token 无效或已过期: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 将解析出来的用户信息存入 Gin 上下文 (c)，方便后续 Handler 随时获取
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next() // 校验通过，放行继续执行后续接口
	}
}
