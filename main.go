package main

import (
	"github.com/gin-gonic/gin"
	"webhook/handler"
)

func main() {
	r := gin.Default()
	r.Any("/", handler.Index)
	r.POST("/webhook", handler.GitHubEvent)
	r.Run(":3009")
}