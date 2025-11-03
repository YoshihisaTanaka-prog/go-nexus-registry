package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"time"
	"web_app/api/handler"
)

var pagePathData =  map[string][]string{
	"all": []string{
		"sign-up",
		"sign-in",
		"apply",
	},
	"editors": []string{
		"manage",
	},
	"admins": []string{
		"role",
	},
}

func Start() {
	handler.Init()

	time.Sleep(time.Second * 1)
	fmt.Fprintln(os.Stdout, "Webサーバを起動します。")

	r := gin.Default()

	r.Use(handler.AuthProxy)

	// /api グループを作成
	api := r.Group("/api/v1")
	{
		api.POST("/sign-up", handler.SignUp)
		api.POST("/sign-in", handler.SignIn)
		api.GET("/my-profile", handler.GetMyProfile)
		api.POST("/apply", handler.Apply)
		api.GET("/get-libraries", handler.EditorAuthProxy, handler.GetLibraries)
		api.POST("/update-is-published", handler.EditorAuthProxy, handler.UpdateIsPublished)
		api.GET("/get-roles", handler.AdminAuthProxy, handler.GetRoles)
		api.GET("/get-role-details", handler.AdminAuthProxy, handler.GetRoleDetails)
		api.POST("/create-role", handler.AdminAuthProxy, handler.CreateRole)
		api.PUT("/update-role", handler.AdminAuthProxy, handler.UpdateRole)
	}

	// r.GET("/sse", handler.SSE)

	for key, paths := range pagePathData {
		switch key {
		case "admins":
			for _, path := range paths{
				r.GET("/" + path, handler.AdminAuthProxy, func(c *gin.Context) {
					c.File("/app/public/htmls/" + path + ".html")
				})
			}
		case "editors":
			for _, path := range paths{
				r.GET("/" + path, handler.EditorAuthProxy, func(c *gin.Context) {
					c.File("/app/public/htmls/" + path + ".html")
				})
			}
		default:
			for _, path := range paths{
				r.GET("/" + path, func(c *gin.Context) {
					c.File("/app/public/htmls/" + path + ".html")
				})
			}
		}
	}

	// 必要に応じて静的ファイルのルートを設定
	r.Static("/assets", "/app/public/assets")
	
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