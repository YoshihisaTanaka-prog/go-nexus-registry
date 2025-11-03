package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"sync"
	"web_app/dbClient"
	"web_app/ldap"
	"web_app/nexus"
)

var accessMutex sync.Mutex

func GetRoles(c *gin.Context) {
	roles, err := dbClient.Role.GetNexusRoles()
	if err != nil {
		fmt.Fprintln(os.Stderr, "get-roles: DB Error:", err)
		c.JSON(400, gin.H{"message": "Internal Server Error"})
		return
	}
	c.JSON(200, roles)
}

func GetRoleDetails(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"message": "Query parameter \"id\" is required"})
		return
	}
	role, err := dbClient.Role.FindById(id)
	if err != nil {
		c.JSON(404, gin.H{"message": "No Such Role id: " + id})
		return
	}

	privileges := []string{}
	users := []string{}
	_, code, uids := ldap.SearchNexusRole(id)

	if code == 404 {
		go ldap.CreateNexusRole(id, role.Name)
	} else if code == 0 {
		users = append(users, uids...)
	} else {
		c.JSON(500, gin.H{"message": "Internal Server Error in Nexus"})
		return
	}
	
	nexusResult, ok := nexus.FindRole(id)
	if !ok {
		c.JSON(500, gin.H{"message": "Internal Server Error in Nexus"})
		return
	}
	if nexusResult == nil {
		go nexus.CreateRole(id, role.Name)
	} else {
		privileges = append(privileges, (*nexusResult).Privileges...)
	}
	
	responseData := struct{
		Id           string `json:"id"`
		Name         string `json:"name"`
		Mode         string `json:"mode"`
		Privileges []string `json:"privileges"`
		Users      []string `json:"users"`
	} {
		Id:         id,
		Name:       role.Name,
		Mode:       role.Mode,
		Privileges: privileges,
		Users:      users,
	}
	c.JSON(200, responseData)
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
			c.JSON(500, gin.H{"message": "Internal Server Error"})
		}
		return
	}

	txt, responseCode := ldap.CreateNexusRole(role.ID, role.Name)

	if responseCode != 0 {
		go ldap.DeleteRole(role.ID)
		c.JSON(responseCode, gin.H{"message": "LDAP error " + txt})
		return
	}

	if !nexus.CreateRole(role.ID, role.Name) {
		go dbClient.Role.DeleteRole(role.ID)
		c.JSON(500, gin.H{"message": "Nexus error "})
		return
	}

	c.JSON(200, role)
}

func UpdateRole(c *gin.Context) {
	accessMutex.Lock()
	defer accessMutex.Unlock()
	
	userId := c.MustGet("userId").(string)

	var body struct {
		Id   string         `json:"id" binding:"required"`
		Name string         `json:"name"`
		Privileges []string `json:"privileges"`
	}

	if err := c.BindJSON(&body); err != nil {
		fmt.Fprintln(os.Stderr, "update-role: リクエストJSONの解析に失敗しました:", err)
		c.JSON(400, gin.H{"error": "リクエストJSONの解析に失敗しました"})
		return
	}

	if body.Name == "" {
		role, err := dbClient.Role.UpdateNexusRole(body.Id, body.Name, userId)
		if err != nil {
			statusCode, message := dbClient.Role.JudgeError(err)
			fmt.Fprintln(os.Stderr, "update-role: DB Error:", message)
			if statusCode == 409 {
				c.JSON(409, gin.H{"message": message})
			} else {
				c.JSON(500, gin.H{"message": "Internal Server Error in DB"})
			}
			return
		}
		if !nexus.UpdateRoleName(role.ID, body.Name) {
			c.JSON(500, gin.H{"message": "Internal Server Error in Nexus"})
			return
		}
	}

	if len(body.Privileges) > 0 {
		role, err := dbClient.Role.FindById(body.Id)
		if err != nil {
			c.JSON(500, gin.H{"message": "Internal Server Error in DB"})
			return
		}
		if !nexus.UpdateRolePrivileges(role.ID, body.Privileges) {
			c.JSON(500, gin.H{"message": "Internal Server Error in Nexus"})
			return
		}
	}

	c.JSON(200, gin.H{})
}
