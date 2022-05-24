package middleware

import (
	token "GoEcommerceApp/tokens"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authetication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.get("token")
		if ClientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No auth header"})
			c.Abort()
			return
		}
		claims, err := token.ValidateToken(ClientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("userId", claims.UserId)
		c.Next()
	}
}
