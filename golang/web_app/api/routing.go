package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	// "net/http"
	"os"
	"time"
	"web_app/api/handler"
)

var pagePaths = []string{
	"sign-up",
	"sign-in",
	"apply",
	"manage",
	"role",
}

func Start() {
	handler.InitJwt()
	handler.SetNpmrc()

	time.Sleep(time.Second * 1)
	fmt.Fprintln(os.Stdout, "Webサーバを起動します。")

	r := gin.Default()

	r.Use(handler.AuthProxy)

	// /api グループを作成
	api := r.Group("/api/v1")
	{
		api.POST("/sign-up", handler.SignUp)
		api.POST("/sign-in", handler.SignIn)
		api.POST("/apply", handler.Apply)
		api.GET("/get-libraries", handler.GetLibraries)
		api.POST("/update-is-published", handler.UpdateIsPublished)
		api.GET("/get-roles", handler.GetRoles)
		api.GET("/get-role-details", handler.GetRoleDetails)
		api.POST("/create-role", handler.CreateRole)
		api.PUT("/update-role", handler.UpdateRole)
	}

	// r.GET("/sse", handler.SSE)

	// 必要に応じて静的ファイルのルートを設定
	r.Static("/assets", "/app/public/assets")
	for _, path := range pagePaths {
		r.GET("/" + path, func(c *gin.Context) {
			c.File("/app/public/htmls/" + path + ".html")
		})
	}

	// /icon.png -> icon.png を返す
	r.GET("/icon.png", func(c *gin.Context) {
		c.File("/app/public/icon.png")
	})
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