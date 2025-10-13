package handler

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"web_app/ldap"
)

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
		jwt, err := CreateJWT(body.Email, []string{})
		if err == nil {
			SetCookie(c, "test", jwt )
			c.JSON(200, gin.H{})
			return
		}
	}
	DeleteCookie(c, "test")
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
		jwt, err := CreateJWT(body.Email, []string{})
		if err == nil {
			SetCookie(c, "test", jwt )
			c.JSON(200, gin.H{})
			return
		}
	}
	DeleteCookie(c, "test")
	c.JSON(code, gin.H{"message": message})
}