package handler

import (
	"web_app/customError"
	"web_app/dbClient"
)

var groupIds = struct{
	Admins  string
	Editors string
	Viewers string
} {"", "", ""}

func Init() {
	strongRoles, err := dbClient.Role.GetNePlusRoles()
	if err != nil {
		customError.Exit1(err)
	}
	if len(strongRoles) == 0 {
		customError.Exit1("No NePlus Role Record was found.")
	}
	for	_, role := range strongRoles {
		if role.Name == "NePlus管理者" {
			groupIds.Admins = role.ID
		}
		if role.Name == "NePlus編集者" {
			groupIds.Editors = role.ID
		}
	}
	viewersRoles, err := dbClient.Role.GetViewersRoles()
	if err != nil {
		customError.Exit1(err)
	}
	if len(viewersRoles) == 0 {
		customError.Exit1("No Viewers Role Record was found.")
	}
	for	_, role := range viewersRoles {
		groupIds.Viewers = role.ID
	}
	setJwt()
	setNpmrc()
}