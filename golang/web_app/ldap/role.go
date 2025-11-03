package ldap

import (
	"fmt"
	"strings"
	"web_app/customError"
	"web_app/dbClient"
)

var viewersId string

func initRole() {
	roles, err := dbClient.Role.GetAllRoles()
	if err != nil {
		customError.Exit1("Initializing LDAP Role Error", err)
	}

	for _, role := range roles {
		id := role.ID
		mode := role.Mode
		txt, code :=  runLdap("Searching Group error", "ldapsearch", []string{"-b", getBaseDn(mode), fmt.Sprintf("(cn=%s)", id)})
		if code == 0 {
			if role.Mode == "viewers" {
				viewersId = id
			}
			foundNum := len(strings.Split(txt, "\ndn:")) - 1
			if foundNum == 0 {
				txt, code = createRole(mode, id, role.Name)
				if code != 0 {
					dbClient.Role.DeleteRole(id)
					customError.Exit1("Initializing LDAP Role Error for ", id, txt)
				}
				if role.Name == "admins" && role.Mode == "apis" {
					AssignApi(id, "admin")
				}
			}
		}
	}
}

func getBaseDn(mode string) string {
	switch mode {
	case "admins", "editors", "viewers", "custom":
		return fmt.Sprintf("ou=nexus_groups,%s", baseDn)
	case "neplus":
		return fmt.Sprintf("ou=neplus_groups,%s", baseDn)
	case "apis":
		return fmt.Sprintf("ou=api_groups,%s", baseDn)
	default:
		return ""
	}
}

func createRole(mode string, id string, name string) (txt string, exitCode int) {
	ldif := strings.Join([]string{
		fmt.Sprintf("dn: cn=%s,%s", id, getBaseDn(mode)),
		fmt.Sprintf("changetype: add"),
		fmt.Sprintf("objectclass: groupOfNames"),
		fmt.Sprintf("cn: %s", id),
		fmt.Sprintf("description: Groups for Nexus named %s", name),
		fmt.Sprintf("member:"),
	}, "\n")
	return runLdap("Creating Group error", "ldapmodify", []string{"-c"}, ldif)
}

func CreateNexusRole(id string, name string) (txt string, exitCode int) {
	return createRole("custom", id, name)
}

func searchRole(mode string, id string) (txt string, exitCode int) {
	return runLdap("Searching Group error", "ldapsearch", []string{"-b", getBaseDn(mode), fmt.Sprintf("(cn=%s)", id)})
}

func SearchNexusRole(id string) (txt string, exitCode int) {
	return searchRole("custom", id)
}

func SearchNePlusRole(id string) (txt string, exitCode int) {
	return searchRole("neplus", id)
}

func SearchApisRole(id string) (txt string, exitCode int) {
	return searchRole("apis", id)
}

func assignBase(roleId string, clientId string, isUser bool) (txt string, exitCode int) {
	role, err := dbClient.Role.FindById(roleId)
	if err != nil {
		return "A DB Error Occurred While Searching Role Record", 400
	}
	if role == nil {
		return "Role Record Not Found.", 404
	}
	mode := role.Mode
	if mode == "apis" {
		if isUser {
			customError.Exit1("Assigning LDAP Role Error", "Trying to add an API Account to a User Role")
		}
	} else {
		if !isUser {
			customError.Exit1("Assigning LDAP Role Error", "Trying to add an User Account to a API Role")
		}
	}
	roleDn := getBaseDn(role.Mode)
	var userDn string
	if strings.HasPrefix(roleDn, "ou=api_groups") {
		userDn = fmt.Sprintf("uid=_api_%s,ou=apis,%s", clientId, baseDn)
	} else {
		userDn = fmt.Sprintf("uid=%s,ou=users,%s", clientId, baseDn)
	}
	roleDn = fmt.Sprintf("cn=%s,%s", roleId, roleDn)

	ldif := strings.Join([]string{
		fmt.Sprintf("dn: %s", roleDn),
		"changetype: modify",
		"add: member",
		fmt.Sprintf("member: %s", userDn),
	}, "\n")
	
	return runLdap("Adding User to Role error", "ldapmodify", []string{"-c"}, ldif)
}

func AssignUser(roleId string, userId string) (txt string, exitCode int) {
	return assignBase(roleId, userId, true)
}

func AssignApi(roleId string, userId string) (txt string, exitCode int) {
	return assignBase(roleId, userId, false)
}

func DeleteRole(id string) (txt string, exitCode int)  {
	return runLdap("Deleting Group error", "ldapdelete", []string{fmt.Sprintf("cn=%s,%s", id, baseDn)})
}
