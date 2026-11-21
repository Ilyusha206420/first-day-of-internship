package http

import (
	"net/http"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/services"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ss *services.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie(accessCookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		publicID, err := ss.ValidateAccessToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		c.Set("userPublicID", publicID)
		c.Next()
	}
}
