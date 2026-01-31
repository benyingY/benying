package shared

import (
	"os"

	"github.com/gin-gonic/gin"
)

func NewServer(serviceName string) *gin.Engine {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestIDMiddleware())
	router.Use(func(c *gin.Context) {
		c.Set("service_name", serviceName)
		c.Next()
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": serviceName,
		})
	})
	return router
}
