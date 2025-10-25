package npm

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"web_app/dbClient"
)

func GetLibraries(c *gin.Context, kind string, name string, v1 int, v2 int, v3 int, limit int) {
	libraries, nextCursor, err := dbClient.SavedLibrary.FindLibraries(kind, name, v1, v2,v3)
	if err != nil {
		fmt.Fprintln(os.StdErr, err)
		c.JSON(500, gin.H{
			data: libraries,
			cursor: nextCursor,
			limit: limit,
		})
		return
	}

	c.JSON(200, gin.H{})
}