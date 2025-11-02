package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"sync"
	"web_app/dbClient"
)

var accessMutex sync.Mutex

func GetRoles(c *gin.Context) {
	accessMutex.Lock()
	defer accessMutex.Unlock()

	roles, err := dbClient.Role.GetNexusRoles()
	if err != nil {
		fmt.Fprintln(os.Stderr, "get-roles: DB Error:", err)
		c.JSON(400, gin.H{"message": "Internal Server Error"})
		return
	}
	c.JSON(200, roles)
}

func CreateRole(c *gin.Context) {
	accessMutex.Lock()
	defer accessMutex.Unlock()

	userId := c.MustGet("userId").(string)
	
	var body struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		fmt.Fprintln(os.Stderr, "create-role: リクエストJSONの解析に失敗しました: ", err)
		c.JSON(400, gin.H{"error": "リクエストJSONの解析に失敗しました"})
		return
	}

	role, err := dbClient.Role.CreateNexusRole(body.Name, userId)
	if err != nil {
		statusCode, message := dbClient.Role.JudgeError(err)
		fmt.Fprintln(os.Stderr, "create-role: DB Error:", message)
		if statusCode == 409 {
			c.JSON(409, gin.H{"message": message})
		} else {
			c.JSON(400, gin.H{"message": "Internal Server Error"})
		}
		return
	}

	c.JSON(200, role)
}

func UpdateRole(c *gin.Context) {
	accessMutex.Lock()
	defer accessMutex.Unlock()
	
	userId := c.MustGet("userId").(string)

	var body struct {
		Id   string `json:"id"   binding:"required"`
		Name string `json:"name" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		fmt.Fprintln(os.Stderr, "update-role: リクエストJSONの解析に失敗しました:", err)
		c.JSON(400, gin.H{"error": "リクエストJSONの解析に失敗しました"})
		return
	}

	role, err := dbClient.Role.UpdateNexusRole(body.Id, body.Name, userId)
	if err != nil {
		statusCode, message := dbClient.Role.JudgeError(err)
		fmt.Fprintln(os.Stderr, "update-role: DB Error:", message)
		if statusCode == 409 {
			c.JSON(409, gin.H{"message": message})
		} else {
			c.JSON(400, gin.H{"message": "Internal Server Error"})
		}
		return
	}

	c.JSON(200, role)
}
