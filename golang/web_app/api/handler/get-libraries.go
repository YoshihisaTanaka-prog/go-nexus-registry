package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"web_app/nexus"
)

func GetLibraries(c *gin.Context) {
	kind := "npm"
	name := c.Query("name")
	var err error
	v1 := 0
	v1Str := c.Query("v1")
	if v1Str != "" {
		v1, err = strconv.Atoi(v1Str)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			c.JSON(500, gin.H{})
			return
		}
	}
	v2 := 0
	v2Str := c.Query("v2")
	if v2Str != "" {
		v2, err = strconv.Atoi(v2Str)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			c.JSON(500, gin.H{})
			return
		}
	}
	v3 := 0
	v3Str := c.Query("v3")
	if v3Str != "" {
		v3, err = strconv.Atoi(v3Str)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			c.JSON(500, gin.H{})
			return
		}
	}
	limit := 10
	limitStr := c.Query("v3")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			c.JSON(500, gin.H{})
			return
		}
	}
	nexus.Npm.GetLibraries(c, kind, name, v1, v2, v3, limit)
}