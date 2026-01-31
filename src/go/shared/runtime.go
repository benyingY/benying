package shared

import (
	"net"

	"github.com/gin-gonic/gin"
)

func RunServer(router *gin.Engine, settings ServiceSettings) error {
	addr := net.JoinHostPort(settings.Host, settings.Port)
	return router.Run(addr)
}
