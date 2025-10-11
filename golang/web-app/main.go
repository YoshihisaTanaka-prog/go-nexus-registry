package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	domainLabels := strings.Split(os.Getenv("LDAP_DOMAIN"), ".");
	if len(domainLabels) == 0 {
		fmt.Fprintln(os.Stderr, "'LDAP_DOMAIN'の値を設定してください。")
		os.Exit(1)
	}
	baseDnSlice := []string{}
	for _, domainLabel := range domainLabels {
		baseDnSlice = append(baseDnSlice, "dc=" + domainLabel)
	}
	baseDn := strings.Join(baseDnSlice, ",")
	bindDn := "cn=" + os.Getenv("LDAP_BIND_CN") + ",ou=service_accounts," + baseDn

	time.Sleep(time.Second * 10)
	cmd := exec.Command(
		"ldapsearch",
		"-x",
		"-H", "ldap://ldap:" + os.Getenv("LDAP_PORT"),
		"-D", bindDn,
		"-w", os.Getenv("LDAP_BIND_PASS"),
		"-b", baseDn,
		"(objectClass=*)",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
	fmt.Fprintf(os.Stderr, "OpenLDAP 実行エラー: %v\n", err)
		os.Exit(1)
	}
}
