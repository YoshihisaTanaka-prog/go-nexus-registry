package ldap

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	bindDn = fmt.Sprintf("cn=%s,ou=service_accounts,%s", os.Getenv("LDAP_BIND_CN"), baseDn)
	connectionArgs = []string{
		"-x",
		"-H", "ldap://ldap",
		"-D", bindDn,
		"-w", os.Getenv("LDAP_BIND_PASS"),
	}
}

func AddUser(email string, password string) (mean string, responseCode int) {
	userName := getUserName(email)
	hashedPassword, err := generateSSHA(password)
	if err != nil {
		return "Encoding hash error", 500
	}

	userDN := getUserDn(userName)
	ldif := fmt.Sprintf("dn: %s\nobjectClass: inetOrgPerson\nuid: %s\ncn: %s\nsn: %s\nmail: %s\nuserPassword: %s\n", userDN, userName, "New User", "New User", email, hashedPassword)
	fmt.Fprintln(os.Stdout, "Adding user:", userName)
	return runLDAP("Adding Ldap error", "ldapadd", []string{}, ldif)
}

func SearchUser(email string) (mean string, responseCode int) {
	userName := getUserName(email)
	fmt.Fprintln(os.Stdout, "Searching user:", userName)
	return runLDAP("Searching User error", "ldapsearch", []string{"-b", baseDn, fmt.Sprintf("\"(uid=%s\")", userName)})
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
	return runLDAP("Changing password error", "ldapmodify", []string{}, ldif)
}

func DeleteUser(email string) (mean string, responseCode int) {
	userName := getUserName(email)
	userDN := getUserDn(userName)
	fmt.Println("Deleting user:", userName)
	return runLDAP("Deleting user error", "ldapdelete", []string{userDN})
}

func Authenticate(email string, password string) (mean string, responseCode int) {
	userName := getUserName(email)
	userDN := getUserDn(userName)
	fmt.Println("Authenticating user:", userName)
	cmd := exec.Command("ldapwhoami", "-x", "-H", "ldap://ldap", "-D", userDN ,"-w", password)
	if err := cmd.Run(); err != nil {
		return customError.GetLdapResult(err, "Authenticating error")
	}
	return "", 0
}