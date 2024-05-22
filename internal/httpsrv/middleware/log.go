package middleware

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(fmt.Sprintf("[%s]: %s", c.Request.Method, c.Request.URL.Path))

		c.Next()
	}
}
