package npm

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"web_app/dbClient"
)

func GetLibraries(c *gin.Context, kind string, name string, v1 int, v2 int, v3 int, limit int) {
	fmt.Fprintln(os.Stdout, kind, limit, name, v1, v2, v3)
	libraries, nextCursor, err := dbClient.SavedLibrary.FindLibraries(kind, name, v1, v2, v3, limit)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		c.JSON(500, gin.H{})
		return
	}
	c.JSON(200, gin.H{
		"data":   libraries,
		"cursor": nextCursor,
		"limit":  limit,
	})
}