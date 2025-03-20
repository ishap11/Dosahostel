package middleware

import (
	"net/http"

	"github.com/adityjoshi/Dosahostel/utils"
	"github.com/gin-gonic/gin"
)

func AuthorizeComplaint() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
			c.Abort()
			return
		}

		claims, err := utils.DecodeStudentJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userType, ok := claims["user"].(map[string]interface{})["type"].(string)
		if !ok || userType != "student" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized for student"})
			c.Abort()
			return
		}

		c.Next()
	}
}
