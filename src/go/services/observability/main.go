package main

import (
	"enterprise-llm-agent-platform/shared"
	"github.com/gin-gonic/gin"
)

const serviceName = "observability"

func buildServer() *gin.Engine {
	return shared.NewServer(serviceName)
}

func main() {
	settings := shared.LoadSettings(serviceName)
	router := buildServer()
	if err := shared.RunServer(router, settings); err != nil {
		panic(err)
	}
}
