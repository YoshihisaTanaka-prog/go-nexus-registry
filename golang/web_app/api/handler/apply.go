package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"os"
	"web_app/nexus"
)

func SetNpmrc() {
	if err := os.MkdirAll("/app/tmp/npm", 0775); err != nil {
		fmt.Fprintln(os.Stderr, "出力ディレクトリ作成エラー:", err)
		os.Exit(1)
	}
}

func Apply(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	var body struct {
		Kind   string `json:"kind"  binding:"required"`
		Name   string `json:"name"  binding:"required"`
		Index *int    `json:"index" binding:"required"`
		V1    *int    `json:"v1"`
		V2    *int    `json:"v2"`
		V3    *int    `json:"v3"`
	}

	if err := c.BindJSON(&body); err != nil {
		fmt.Fprintln(os.Stderr, "sign-up\nリクエストJSONの解析に失敗しました\n", err)
		c.JSON(400, gin.H{"error": "リクエストJSONの解析に失敗しました"})
		return
	}

	uuidV4, err := uuid.NewRandom()
	if err != nil {
		c.JSON(500, gin.H{"error": "UUIDの生成に失敗しました"})
		return
	}

	nexus.Npm.Apply(c, body, userId, uuidV4)
}