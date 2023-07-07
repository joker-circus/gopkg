package ginutil

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 跨域请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		var headerKeys []string
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}

		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			// 可将将* 替换为指定的域名
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			// c.AbortWithStatus(http.StatusNoContent)
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}

// 严格版的跨域请求
func StrictCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}

// 检测 Params 参数是否合法
func CheckParamsIsValid() gin.HandlerFunc {
	return func(c *gin.Context) {
		paramsMap := make(map[string]string)
		for _, param := range c.Params {
			paramsMap[param.Key] = param.Value
			if param.Value == "" {
				c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("params【%s】 is required!", param.Key)})
				c.Abort()
				return
			}
			// 如果参数包含空格
			if strings.Contains(param.Value, " ") {
				c.JSON(http.StatusBadRequest, fmt.Sprintf("Params[%s] cannot be empty or contain spaces!", param.Key))
				c.Abort()
				return
			}
		}
	}
}

// 若 Params 参数以 uid 结尾，进行 uint64 转换
func CheckParamsUintId() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, value := range c.Params {
			// 如果参数末尾是uid，校验uint64
			if strings.HasSuffix(value.Key, "uid") {
				_, err := strconv.ParseUint(value.Value, 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, fmt.Sprintf("Params[%s] not valid", value.Key))
					c.Abort()
					return
				}
			}
		}
	}
}
