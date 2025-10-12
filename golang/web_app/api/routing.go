package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var pagePathes = []string{
	"sign-up",
	"sign-in",
	"apply",
}

func Start() {
	r := gin.Default()

	// // TodoエンティティのCRUDエンドポイントを設定
	// r.POST("/todos", func(c *gin.Context) { /* 作成処理 */ })
	// r.GET("/todos", func(c *gin.Context) { /* 一覧取得処理 */ })
	// r.GET("/todos/:id", func(c *gin.Context) { /* 個別取得処理 */ })
	// r.PUT("/todos/:id", func(c *gin.Context) { /* 更新処理 */ })
	// r.DELETE("/todos/:id", func(c *gin.Context) { /* 削除処理 */ })

	// /api グループを作成
	api := r.Group("/api/v1")
	{
		api.POST("/sign-up", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "POST /api/todos"})
		})
		api.POST("/sign-in", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "POST /api/todos"})
		})
	}

	// 必要に応じて静的ファイルのルートを設定
	r.Static("/assets", "/app/public/assets")
	for _, path := range pagePathes {
		r.GET("/" + path, func(c *gin.Context) {
			c.File("/app/public/htmls/" + path + ".html")
		})
	}

	// /favicon.ico -> favicon.ico を返す
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("/app/public/favicon.ico")
	})

	// / → index.html を返す
	r.GET("/", func(c *gin.Context) {
		c.File("/app/public/htmls/index.html")
	})

	// サーバーを起動
	r.Run("0.0.0.0:8080")
}