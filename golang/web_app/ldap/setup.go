package ldap

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// 追加で注入したい環境変数名リスト（平文）
var additionalEnvKeys = []string{
	"LDAP_ADMIN_PASSWORD",
	"LDAP_ADMIN_USERNAME",
	"LDAP_BIND_CN_NEXUS",
	"LDAP_BIND_PASS_NEXUS",
	"LDAP_DOMAIN",
	"LDAP_ORGANISATION",
	"LDAP_PORT",
}

var savedEnvVarsJsonFilePath = "/app/saved-data/saved-env-vars.json"

func exit1(args ...any)  {
	if len(args) > 0 {
		fmt.Fprintln(os.Stderr, args...)
	}
	os.Exit(1)
}

func setupLdap(apiPasswords map[string]string) {
	envVars := getEnvVars(apiPasswords)
	adminDn := fmt.Sprintf("cn=%s,%s", envVars["LDAP_ADMIN_USERNAME"], baseDn)
	adminPassword := envVars["LDAP_ADMIN_PASSWORD"]
	_, exitCode := runLdapAsUser("OpenLDAP 設定エラー", "ldapsearch", adminDn, adminPassword, []string{"-b", baseDn, fmt.Sprintf("(cn=%s)", envVars["LDAP_BIND_CN_NEXUS"])})
	if exitCode == 0 {
		updateSchema("/app/templates/update.ldif.template", adminDn, adminPassword, envVars)
	} else {
		updateSchema("/app/templates/init.ldif.template", adminDn, adminPassword, envVars)
	}
}

func getEnvVars(apiPasswords map[string]string) map[string]string {
	domainLabels := strings.Split(os.Getenv("LDAP_DOMAIN"), ".");
	if len(domainLabels) == 0 {
		exit1("'LDAP_DOMAIN'の値を設定してください。")
	}
	baseDnSlice := []string{}
	for _, domainLabel := range domainLabels {
		baseDnSlice = append(baseDnSlice, "dc=" + domainLabel)
	}
	baseDn := strings.Join(baseDnSlice, ",")
	rootDc := domainLabels[0]

	adminPass := os.Getenv("LDAP_ADMIN_PASSWORD")
	bindPassNexus := os.Getenv("LDAP_BIND_PASS_NEXUS")

	if adminPass == "" || bindPassNexus == "" {
		exit1("環境変数 LDAP_ADMIN_PASSWORD, LDAP_BIND_PASS_NEXUS が必要です。")
	}

	vars := map[string]string{
		"LDAP_ROOT": baseDn,
		"LDAP_ROOT_DC": rootDc,
		"LDAP_ADMIN_PASS_HASH": generateSSHAForLdap(adminPass),
		"LDAP_BIND_PASS_NEXUS_HASH": generateSSHAForLdap(bindPassNexus),
		"LDAP_BIND_PASS_ADMIN_API_HASH": generateSSHAForLdap(apiPasswords["admin"]),
	}

	for _, key := range additionalEnvKeys {
		vars[key] = os.Getenv(key)
	}

	return vars
}

func updateSchema(templateFilePath string, dn string, password string, envVars map[string]string) {
	templateFileData, err := os.ReadFile(templateFilePath)
	if err != nil {
		exit1("テンプレート読み込みエラー:", err)
	}

	fmt.Fprintln(os.Stdout, templateFilePath + " を読み込みました。\n環境変数を注入します。")
	content := substituteEnv(string(templateFileData), envVars)
	fmt.Fprintln(os.Stdout, "環境変数を注入しました。\nLDAPのデータを更新します。")

	_, exitCode := runLdapAsUser("OpenLDAP 設定エラー", "ldapmodify", dn, password, []string{"-c"}, content)
	if exitCode != 0 {
		os.Exit(1)
	}
}

func toJsonBytes(mapData map[string]string) []byte {
	filteredMapData := make(map[string]string)

	for k, v := range mapData {
		if !strings.HasSuffix(k, "_HASH") {
			filteredMapData[k] = v
		}
	}

	bytes, err := json.Marshal(filteredMapData)
	if err != nil {
		exit1("JSON marshal error:", err)
	}
	return bytes
}

func loadJsonBytes() []byte {
	bytes, err := os.ReadFile(savedEnvVarsJsonFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		exit1("failed to read file:", err)
	}
	return bytes
}

func compareEnvVarsAndUpdateSaveData(currentEnvVarsBytes []byte, savedEnvVarsBytes []byte) bool {
	defer fmt.Fprintln(os.Stdout, "セーブデータと現在のデータを比較しました。")
	fmt.Fprintln(os.Stdout, "セーブデータと現在のデータを比較中...") 
	if bytes.Equal(currentEnvVarsBytes, savedEnvVarsBytes) {
		return true
	} else {
		go func() {
			fmt.Fprintln(os.Stdout, "セーブデータを更新しています...")
			if err := os.WriteFile(savedEnvVarsJsonFilePath, currentEnvVarsBytes, 0644); err != nil {
				fmt.Fprintln(os.Stderr, "セーブデータの更新に失敗しました。\n", err)
				os.Exit(1)
			}
			fmt.Fprintln(os.Stdout, "セーブデータを更新しました。")
		}()
		return false
	}
}

func generateSSHAForLdap(password string) string {
	salt := make([]byte, 4)
	_, err := rand.Read(salt)
	if err != nil {
		exit1(err)
	}

	h := sha1.New()
	h.Write([]byte(password))
	h.Write(salt)
	hash := h.Sum(nil)

	ssha := append(hash, salt...)
	return "{SSHA}" + base64.StdEncoding.EncodeToString(ssha)
}

func substituteEnv(input string, vars map[string]string) string {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	return re.ReplaceAllStringFunc(input, func(s string) string {
		key := re.FindStringSubmatch(s)[1]
		if val, ok := vars[key]; ok {
			return val
		}
		if val := os.Getenv(key); val != "" {
			return val
		}
		return s
	})
}