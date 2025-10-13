package handler

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"net/url"
	"strings"
	"web_app/ldap"
)

var cookieSessionKey = "_session"

var publicPaths = []string{
	"/assets/api.js",
	"/assets/api.css",
	"/assets/BaseBase.js",
	"/assets/BaseBase.css",
	"/.well-known",
}

var guestOnlyPaths = []string{
	"/sign-in",
	"/sign-up",
	"/api/v1/sign-in",
	"/api/v1/sign-up",
	"/assets/sign-in.js",
	"/assets/sign-up.js",
	"/favicon.ico",
}

func SignUp(c *gin.Context) {
	// ctx = c.Request.Context()

	var body struct {
		Email  string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		fmt.Fprintln(os.Stderr, "sign-up\nリクエストJSONの解析に失敗しました\n", err)
		c.JSON(400, gin.H{"error": "リクエストJSONの解析に失敗しました"})
		return
	}

	fmt.Fprintln(os.Stdout, "sign-up:", body)

	message, code := ldap.AddUser(body.Email, body.Password)
	fmt.Println(message, code)

	if code == 0 {
		jwt, err := CreateJwt(body.Email, []string{})
		if err == nil {
			SetCookie(c, cookieSessionKey, jwt )
			c.JSON(200, gin.H{})
			return
		}
	}
	DeleteCookie(c, cookieSessionKey)
	c.JSON(code, gin.H{"message": message})
}

func SignIn(c *gin.Context) {
	// ctx = c.Request.Context()

	var body struct {
		Email  string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		fmt.Fprintln(os.Stderr, "sign-up\nリクエストJSONの解析に失敗しました\n", err)
		c.JSON(400, gin.H{"error": "リクエストJSONの解析に失敗しました"})
		return
	}

	fmt.Fprintln(os.Stdout, "sign-in:", body)

	message, code := ldap.Authenticate(body.Email, body.Password)
	fmt.Println(message, code)

	if code == 0 {
		jwt, err := CreateJwt(body.Email, []string{})
		if err == nil {
			SetCookie(c, cookieSessionKey, jwt )
			c.JSON(200, gin.H{})
			return
		}
	}
	DeleteCookie(c, cookieSessionKey)
	c.JSON(code, gin.H{"message": message})
}

func AuthProxy(c *gin.Context) {
	path := c.Request.URL.Path

	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			c.Next()
			return
		}
	}

	isNeedAuth := true
	for _, guestOnlyPath := range guestOnlyPaths {
		if strings.HasPrefix(path, guestOnlyPath) {
			isNeedAuth = false
			break
		}
	}

	redirectPath := "/sign-in?redirect=" + url.PathEscape(path)
	cookie, err := c.Cookie(cookieSessionKey)
	if err == nil {
		if isNeedAuth {
			roles := ParseJwt(cookie)
			if slices.Contains(roles, "_401") {
				DeleteCookie(c, cookieSessionKey)
				c.Redirect(302, redirectPath)
				c.Abort()
			} else {
				c.Set("roles", roles)
				c.Next()
			}
			return
		} else {
			DeleteCookie(c, cookieSessionKey)
			c.Redirect(302, "/")
			c.Abort()
			return
		}
	} else {
		if isNeedAuth {
			DeleteCookie(c, cookieSessionKey)
			c.Redirect(302, redirectPath)
			c.Abort()
			return
		} else {
			c.Next()
			return
		}
	}
}