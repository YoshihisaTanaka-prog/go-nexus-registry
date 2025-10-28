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
	"web_app/customError"
)

func runLDAP(errorMessage string, command string, args []string, inputs ...string) (mean string, responseCode int) {
	localArgs := []string{}
	localArgs = append(localArgs, connectionArgs...)
	localArgs = append(localArgs, args...)
	cmd := exec.Command(command, localArgs...)
	joinedInputs := ""
	if len(inputs) > 0 {
		joinedInputs = strings.Join(inputs, "\n")
		cmd.Stdin = bytes.NewBufferString(joinedInputs)
	}
	fmt.Printf("Executing: %s %s\n%s\n", command, strings.Join(localArgs, " "), joinedInputs)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return customError.GetLdapResult(err, errorMessage)
	}
	return "", 0
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
		fmt.Fprintln(os.Stderr, "Error has ocured while encoding hash.")
		return "", err
	}

	h := sha1.New()
	h.Write([]byte(password))
	h.Write(salt)
	hash := h.Sum(nil)

	ssha := append(hash, salt...)
	return "{SSHA}" + base64.StdEncoding.EncodeToString(ssha), nil
}
