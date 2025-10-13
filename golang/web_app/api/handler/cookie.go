package handler

import (
	"github.com/gin-gonic/gin"
	"os"
)

func SetCookie(c *gin.Context, key string, value string) {
	secure := true
	hostName := os.Getenv("GO_MANAGER_HOST_NAME")
	if hostName == "localhost" {
		secure = false
	}
	c.SetCookie(key, value, 0, "/", os.Getenv("GO_MANAGER_HOST_NAME"), secure, true)
}