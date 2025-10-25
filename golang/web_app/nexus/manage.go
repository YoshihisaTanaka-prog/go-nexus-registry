package nexus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"web_app/dbClient"
)

func GetLibraries(c *gin.Context, kind string, name string, v1 int, v2 int, v3 int, limit int) {
	fmt.Fprintln(os.Stdout, kind, limit, name, v1, v2, v3)
	libraries, nextCursor, err := dbClient.SavedLibrary.FindLibraries(kind, name, v1, v2, v3, limit)
	if err == nil {
		c.JSON(200, gin.H{
			"data":   libraries,
			"cursor": nextCursor,
			"limit":  len(libraries),
		})
		return
	}
	fmt.Fprintln(os.Stderr, err)
	c.JSON(500, gin.H{})
}

func UpdateIsPublished(c *gin.Context, id string, isPublished bool) {
	fmt.Fprintln(os.Stdout, id, isPublished)
	library, err := dbClient.SavedLibrary.UpdateIsPublished(id, isPublished)
	if err == nil {
		c.IndentedJSON(200, library)
		return
	}
	fmt.Fprintln(os.Stderr, err)
	c.JSON(500, gin.H{})
}