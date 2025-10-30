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

func InitLdap() {
	domainLabels := strings.Split(os.Getenv("LDAP_DOMAIN"), ".");
	baseDnSlice := []string{}
	for _, domainLabel := range domainLabels {
		baseDnSlice = append(baseDnSlice, "dc=" + domainLabel)
	}
	baseDn = strings.Join(baseDnSlice, ",")
	bindDn = fmt.Sprintf("cn=%s,%s", os.Getenv("LDAP_ADMIN_USERNAME"), baseDn)
	connectionArgs = []string{
		"-x",
		"-H", "ldap://ldap:" + os.Getenv("LDAP_PORT"),
		"-D", bindDn,
		"-w", os.Getenv("LDAP_ADMIN_PASSWORD"),
	}
	go func() {
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
				fmt.Println("Passed authentication")
				go setupLdap()
				break
			} else {
				customError.GetLdapResult(err, "Authenticating error")
			}
		}
	}()
}

func AddUser(email string, password string) (mean string, responseCode int) {
	if (!strings.HasSuffix(email, "@" + os.Getenv("LDAP_DOMAIN"))) {
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
	txt, exitCode := runLdap("Adding Ldap error", "ldapadd", []string{}, ldif)
	if exitCode == 0 {
		return userName, 0
	}
	return txt, exitCode
}

func SearchUser(email string) (mean string, responseCode int) {
	userName := getUserName(email)
	fmt.Fprintln(os.Stdout, "Searching user:", userName)
	txt, exitCode := runLdap("Searching User error", "ldapsearch", []string{"-b", baseDn, fmt.Sprintf("(uid=%s)", userName)})
	if exitCode == 0 {
		return userName, 0
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
	fmt.Println("Deleting user:", userName)
	txt, exitCode := runLdap("Deleting user error", "ldapdelete", []string{userDN})
	if exitCode == 0 {
		return userName, 0
	}
	return txt, exitCode
}

func Authenticate(email string, password string) (mean string, responseCode int) {
	if (!strings.HasSuffix(email, "@" + os.Getenv("LDAP_DOMAIN"))) {
		return "Invalid Email Domain error", 422
	}
	userName := getUserName(email)
	userDN := getUserDn(userName)
	fmt.Println("Authenticating user:", userName)
	txt, exitCode := runLdapAsUser("Authenticating error", "ldapwhoami", userDN, password, []string{})
	if exitCode == 0 {
		return userName, 0
	}
	return txt, exitCode
}