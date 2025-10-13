package ldap

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
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

func AddUser(email string, password string) error {
	userName := getUserName(email)
	hashedPassword, err := generateSSHA(password)
	if err != nil {
		return err
	}

	userDN := getUserDn(userName)
	ldif := fmt.Sprintf("dn: %s\nobjectClass: inetOrgPerson\nuid: %s\ncn: %s\nsn: %s\nmail: %s\nuserPassword: %s\n", userDN, userName, "New User", "New User", email, hashedPassword)
	fmt.Fprintln(os.Stdout, "Adding user:", userName)

	if err := runLDAP("ldapadd", []string{}, ldif); err != nil {
		return err
	}
	return nil
}

func SearchUser(email string) error {
	userName := getUserName(email)
	fmt.Fprintln(os.Stdout, "Searching user:", userName)
	if err := runLDAP("ldapsearch", []string{"-b", baseDn, fmt.Sprintf("\"(uid=%s\")", userName)}); err != nil {
		return err
	}
	return nil
}

func ChangePassword(email string, password string) error {
	userName := getUserName(email)
	hashedPassword, err := generateSSHA(password)
	if err != nil {
		return err
	}

	userDN := getUserDn(userName)
	ldif := fmt.Sprintf("dn: %s\nchangetype: modify\nreplace: userPassword\nuserPassword: %s\n", userDN, hashedPassword)
	fmt.Fprintln(os.Stdout, "Changing password for:", userName)
	if err := runLDAP("ldapmodify", []string{}, ldif); err != nil {
		return err
	}
	return nil
}

func DeleteUser(email string) error {
	userName := getUserName(email)
	userDN := getUserDn(userName)
	fmt.Println("Deleting user:", userName)
	if err := runLDAP("ldapdelete", []string{userDN}); err != nil {
		return err
	}
	return nil
}

func Authenticate(email string, password string) {
	userName := getUserName(email)
	userDN := getUserDn(userName)
	fmt.Println("Authenticating user:", userName)
	cmd := exec.Command("ldapwhoami", "-x", "-H", "ldap://ldap", "-D", userDN ,"-w", password)
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// システムコールの情報を取得
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				fmt.Printf("Exit Code: %d\n", status.ExitStatus())
			}
		} else {
			fmt.Printf("コマンド実行エラー: %v\n", err)
		}
	}
}

func runLDAP(command string, args []string, inputs ...string) error {
	localArgs := []string{}
	localArgs = append(localArgs, connectionArgs...)
	localArgs = append(localArgs, args...)
	cmd := exec.Command(command, localArgs...)
	if len(inputs) > 0 {
		joinedInputs := strings.Join(inputs, "\n")
		cmd.Stdin = bytes.NewBufferString(joinedInputs)
		fmt.Printf("Executing: %s %s %s\n", command, strings.Join(localArgs, " "), joinedInputs)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func getUserName(email string) string {
	emailParts := strings.Split(email, "@")
	return emailParts[0]
}

func getUserDn(userName string) string {
	return fmt.Sprintf("uid=%s,ou=users,%s", userName, baseDn)
}

func generateSSHA(password string) (string, error) {
	salt := make([]byte, 4)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	h := sha1.New()
	h.Write([]byte(password))
	h.Write(salt)
	hash := h.Sum(nil)

	ssha := append(hash, salt...)
	return "{SSHA}" + base64.StdEncoding.EncodeToString(ssha), nil
}
