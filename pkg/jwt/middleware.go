package jwt

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Middleware(h gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "midlleware unauthorized",
			})
			return
		}

		headers := strings.Split(authHeader, " ")
		if len(headers) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "midlleware unauthorized",
			})
			return
		}

		if headers[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "midlleware unauthorized",
			})
			return
		}

		log.Println(headers[1])

		token, err := ParseAccessToken(headers[1])
		log.Println(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "midlleware unauthorized",
			})
			return
		}

		h(ctx)
	}
}
