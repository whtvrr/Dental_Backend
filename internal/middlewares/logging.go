package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/whtvrr/Dental_Backend/pkg/logger"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := logger.GerLogger()
		l.Info()
		c.Next()
	}
}
