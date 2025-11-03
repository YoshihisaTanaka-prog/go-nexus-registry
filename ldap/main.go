package main

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// 🔸 追加で注入したい環境変数名リスト（平文）
var additionalEnvKeys = []string{
	"LDAP_ADMIN_PASSWORD",
	"LDAP_ADMIN_USERNAME",
	"DOMAIN_NAME",
	"LDAP_PORT",
}

func main()  {
	domainLabels := strings.Split(os.Getenv("DOMAIN_NAME"), ".");
	if len(domainLabels) == 0 {
		fmt.Fprintln(os.Stderr, "'DOMAIN_NAME'の値を設定してください。")
		os.Exit(1)
	}
	baseDnSlice := []string{}
	for _, domainLabel := range domainLabels {
		baseDnSlice = append(baseDnSlice, "dc=" + domainLabel)
	}
	baseDn := strings.Join(baseDnSlice, ",")

	envVars := convertEnvVars(baseDn, domainLabels[0])
	
	fmt.Fprintln(os.Stdout, "設定ファイルを更新します。") 
	substituteSlapdConfFile(envVars)
	fmt.Fprintln(os.Stdout, "openldapサーバを起動します。") 
	cmd := exec.Command("slapd", "-u", "root", "-g", "root", "-h", "ldap://0.0.0.0:" + os.Getenv("LDAP_PORT"), "-d", "320")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
	fmt.Fprintf(os.Stderr, "OpenLDAP 実行エラー: %v\n", err)
		os.Exit(1)
	}
}

func convertEnvVars(baseDn string, rootDc string) map[string]string {
	adminPass := os.Getenv("LDAP_ADMIN_PASSWORD")

	if adminPass == "" {
		fmt.Fprintln(os.Stderr, "環境変数 LDAP_ADMIN_PASSWORD が必要です。")
		os.Exit(1)
	}

	vars := map[string]string{
		"LDAP_ROOT": baseDn,
		"LDAP_ROOT_DC": rootDc,
		"LDAP_ADMIN_PASS_HASH": generateSSHA(adminPass),
	}

	for _, key := range additionalEnvKeys {
		vars[key] = os.Getenv(key)
	}

	return vars
}

func generateSSHA(password string) string {
	salt := make([]byte, 4)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}

	h := sha1.New()
	h.Write([]byte(password))
	h.Write(salt)
	hash := h.Sum(nil)

	ssha := append(hash, salt...)
	return "{SSHA}" + base64.StdEncoding.EncodeToString(ssha)
}

func substituteSlapdConfFile(envVars map[string]string) {
	keptFilePath := "/customized/saved-data/slapd.conf"
	var readFilePath string
	if _, err := os.Stat(keptFilePath); err == nil {
		readFilePath = keptFilePath
	} else {
		readFilePath = "/etc/openldap/slapd.conf"
	}
	writeFilePath := "/etc/openldap/slapd.conf"

	fmt.Fprintln(os.Stdout, readFilePath + " を読み込みます。") 
	fileData, err := os.ReadFile(readFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, readFilePath + " 読み込みエラー:", err)
		os.Exit(1)
	}

	if readFilePath == writeFilePath {
		// ファイルに書き込み
		if err := os.WriteFile(keptFilePath, fileData, 0644); err != nil {
			fmt.Fprintln(os.Stderr, keptFilePath + " 書き込みエラー:", err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stdout, keptFilePath + " を更新しました。")
	}

	fmt.Fprintln(os.Stdout, readFilePath + " を読み込みました。\n環境変数を注入します。")

	newContentsLines := []string{}
	
	contentsLines := strings.Split(string(fileData), "\n")

	for _, contentsLine := range contentsLines {
		if (strings.HasPrefix(contentsLine, "suffix")) {
			newContentsLines = append(newContentsLines, "suffix		\"" + envVars["LDAP_ROOT"] + "\"")
		} else if (strings.HasPrefix(contentsLine, "rootdn")) {
			newContentsLines = append(newContentsLines, "rootdn		\"cn=" + envVars["LDAP_ADMIN_USERNAME"] + "," + envVars["LDAP_ROOT"] + "\"")
		} else if (strings.HasPrefix(contentsLine, "rootpw")) {
			newContentsLines = append(newContentsLines, "rootpw		" + envVars["LDAP_ADMIN_PASS_HASH"])
		} else {
			newContentsLines = append(newContentsLines, contentsLine)
		}
	}

	newContents := []byte(strings.Join(newContentsLines, "\n"))
	// ファイルに書き込み
	if err := os.WriteFile(writeFilePath, []byte(newContents), 0644); err != nil {
		fmt.Fprintln(os.Stderr, writeFilePath + " 書き込みエラー:", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, writeFilePath + " を更新しました。")
}