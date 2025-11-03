package ldap

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"web_app/customError"
)

var (
	connectionArgs []string
	baseDn string
	bindDn string
)

func InitLdap(apiPasswords map[string]string) {
	domainLabels := strings.Split(os.Getenv("DOMAIN_NAME"), ".");
	baseDnSlice := []string{}
	for _, domainLabel := range domainLabels {
		baseDnSlice = append(baseDnSlice, "dc=" + domainLabel)
	}
	baseDn = strings.Join(baseDnSlice, ",")
	bindDn = fmt.Sprintf("cn=%s,%s", os.Getenv("LDAP_ADMIN_USERNAME"), baseDn)
	connectionArgs = []string{
		"-x",
		"-H", "ldap://ldap:" + os.Getenv("LDAP_PORT"),
	}
	for i := 0; i <= 60; i++ {
		if i == 60 {
			exit1("Failed to Connect to LDAP Server")
		}
		time.Sleep(time.Second)
		cmd := exec.Command(
			"ldapwhoami",
			connectionArgs...
		)
		if err := cmd.Run(); err == nil {
			fmt.Fprintln(os.Stdout, "Passed authentication")
			setupLdap(apiPasswords)
			initRole()
			break
		} else {
			customError.GetLdapResult(err, "Authenticating error")
		}
	}
}

func AddUser(email string, password string) (mean string, responseCode int) {
	if (!strings.HasSuffix(email, "@" + os.Getenv("DOMAIN_NAME"))) {
		return "Invalid Email Domain error", 422
	}
	userName := getUserName(email)
	hashedPassword, err := generateSSHA(password)
	if err != nil {
		return "Encoding hash error", 500
	}

	userDN := getUserDn(userName)
	ldif := fmt.Sprintf("dn: %s\nobjectClass: inetOrgPerson\nuid: %s\ncn: %s\nsn: %s\nmail: %s\nuserPassword: %s\n", userDN, userName, "New User", "New User", email, hashedPassword)
	fmt.Fprintln(os.Stdout, "Adding user:", userName)
	txt, exitCode := runLdap("Adding User error", "ldapadd", []string{}, ldif)
	if exitCode == 0 {
		txt, exitCode = AssignUser(viewersId, userName)
		if exitCode == 0 {
			return userName, 0
		}
	}
	return txt, exitCode
}

func SearchAllUser() []string {
	users := []string{}
	txt, exitCode := runLdap("Searching User error", "ldapsearch", []string{"-b", fmt.Sprintf("ou=users,%s",baseDn), "(objectClass=inetOrgPerson)"})
	if exitCode != 0 {
		return users
	}
	lines := strings.Split(txt, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "dn: ") {
			splitedTexts := strings.Split(line, ",")
			users = append(users, strings.Split(splitedTexts[0], "=")[1])
		}
	}
	return users;
}

func SearchUser(email string) (mean string, responseCode int) {
	userName := getUserName(email)
	fmt.Fprintln(os.Stdout, "Searching user:", userName)
	txt, exitCode := runLdap("Searching User error", "ldapsearch", []string{"-b", fmt.Sprintf("ou=users,%s",baseDn), fmt.Sprintf("(uid=%s)", userName)})
	if exitCode == 0 {
		splitedTexts := strings.Split(txt, "\ndn: ")
		if len(splitedTexts) > 1 {
			return userName, 0
		}
		return "", 404
	}
	return txt, exitCode
}

func ChangePassword(email string, password string) (mean string, responseCode int) {
	userName := getUserName(email)
	hashedPassword, err := generateSSHA(password)
	if err != nil {
		return "Encoding hash error", 500
	}

	userDN := getUserDn(userName)
	ldif := fmt.Sprintf("dn: %s\nchangetype: modify\nreplace: userPassword\nuserPassword: %s\n", userDN, hashedPassword)
	fmt.Fprintln(os.Stdout, "Changing password for:", userName)
	txt, exitCode := runLdap("Changing password error", "ldapmodify", []string{}, ldif)
	if exitCode == 0 {
		return userName, 0
	}
	return txt, exitCode
}

func DeleteUser(email string) (mean string, responseCode int) {
	userName := getUserName(email)
	userDN := getUserDn(userName)
	fmt.Fprintln(os.Stdout, "Deleting user:", userName)
	txt, exitCode := runLdap("Deleting user error", "ldapdelete", []string{userDN})
	if exitCode == 0 {
		return userName, 0
	}
	return txt, exitCode
}

func Authenticate(email string, password string) (mean string, responseCode int, groups []string) {
	if (!strings.HasSuffix(email, "@" + os.Getenv("DOMAIN_NAME"))) {
		return "Invalid Email Domain error", 422, []string{}
	}
	userName := getUserName(email)
	userDN := getUserDn(userName)
	fmt.Fprintln(os.Stdout, "Authenticating user:", userName)
	txt, exitCode := runLdapAsUser("Authenticating error", "ldapwhoami", userDN, password, []string{})
	if exitCode == 0 {
		txt, exitCode := runLdap(
			"Searching Assigned Group error",
			"ldapsearch",
			[]string{
				"-b", baseDn,
				fmt.Sprintf("(member=uid=%s,ou=users,%s)", userName, baseDn),
			},
		)
		if exitCode == 0 {
			groups := []string{}
			lines := strings.Split(txt, "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "cn: ") {
					splitedTexts := strings.Split(line, "cn: ")
					groups = append(groups, splitedTexts[1])
				}
			}
			return userName, 0, groups
		} else {
			return "Searching Assigned Group error", 500, []string{}
		}
	}
	return txt, exitCode, []string{}
}