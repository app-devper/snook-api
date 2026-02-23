package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewCors(allowOrigins []string) gin.HandlerFunc {
	config := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}
	if len(allowOrigins) == 1 && allowOrigins[0] == "*" {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = allowOrigins
		config.AllowCredentials = true
	}
	return cors.New(config)
}
