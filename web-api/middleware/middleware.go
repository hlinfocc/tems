package middleware

import (
	"log"
	"net/http"
	"strings"
	"tems-web-api/utils"

	"github.com/gin-gonic/gin"
)

// JWT认证中间件
func JWTAuthMiddleware(appType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 排除登录接口和信息绑定接口，不进行JWT认证
		path := c.Request.URL.Path
		// 定义白名单路径集合
		skipPaths := map[string]bool{
			"/api/v1/banners":        true,
			"/api/v1/classes/list":   true,
			"/api/v1/students/login": true,
			"/api/v1/students/bind":  true,
			"/manager/api/login":     true,
		}

		if skipPaths[path] {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.GetHeader("token")
		}
		log.Println(token)
		customClaims, e := utils.JWTParse(token, appType)
		if e != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "鉴权失败",
			})
			// 中止请求处理链！！！
			c.Abort()
			return
		}
		c.Set("customClaims", customClaims)
		c.Next()
	}
}

// 中间件扩展：按域名后缀分组
func DomainGroupMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := strings.Split(c.Request.Host, ":")[0]
		// 根据域名处理相应内容
		c.Set("domain", domain)
		c.Next()
	}
}

// 中间件扩展：拦截静态资源目录访问
func NoDirectoryListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		staticUriPrefixes := []string{
			"/uploads",
		}
		// 检查请求路径是否以 / 结尾（目录请求）
		for _, item := range staticUriPrefixes {
			if utils.HasPrefixIgnoreCase(c.Request.URL.Path, item) && strings.HasSuffix(c.Request.URL.Path, "/") {
				c.AbortWithStatus(404) // 返回 404
				return
			}
		}

		c.Next()
	}
}
