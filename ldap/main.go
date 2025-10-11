package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// 🔸 追加で注入したい環境変数名リスト（平文）
var additionalEnvKeys = []string{
	"LDAP_ADMIN_PASSWORD",
	"LDAP_ADMIN_USERNAME",
	"LDAP_BIND_CN_GO",
	"LDAP_BIND_CN_NEXUS",
	"LDAP_BIND_PASS_GO",
	"LDAP_BIND_PASS_NEXUS",
	"LDAP_DOMAIN",
	"LDAP_ORGANISATION",
	"LDAP_PORT",
}

var savedEnvVarsJsonFilePath = "/customized/saved-data/saved-env-vars.json"

func main()  {
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

	envVars := convertEnvVars(baseDn, domainLabels[0])

	currentEnvVarsBytes := toJsonBytes(envVars)
	savedEnvVarsBytes := loadJsonBytes()

	if (!compareEnvVarsAndUpdateSaveData(currentEnvVarsBytes, savedEnvVarsBytes)) {
		fmt.Fprintln(os.Stdout, "環境変数の変更が確認されたので、設定ファイルを更新します。") 
		substituteSlapdConfFile(envVars)
		makeSubstitutedFile(envVars, "init")
		fmt.Fprintln(os.Stdout, "設定ファイルの更新が完了しました。") 
		go func() {
			defer fmt.Fprintln(os.Stdout, "設定ファイルの更新をシステムに反映させました。") 
			time.Sleep(time.Second * 3)
			fmt.Fprintln(os.Stdout, "設定ファイルの更新をシステムに反映させます。")
			cmd := exec.Command(
				"ldapadd",
				"-x",
				"-H", "ldap://localhost",
				"-D", "cn=" + envVars["LDAP_ADMIN_USERNAME"] + "," + envVars["LDAP_ROOT"],
				"-w", envVars["LDAP_ADMIN_PASSWORD"],
				"-f", "/customized/saved-data/init/init.ldif",
				"-d", "320",
			)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "OpenLDAP 設定反映エラー: %v\n", err)
				os.Exit(1)
			}
		}()
	}
	
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
	bindPassNexus := os.Getenv("LDAP_BIND_PASS_NEXUS")
	bindPassGo := os.Getenv("LDAP_BIND_PASS_GO")

	if adminPass == "" || bindPassNexus == "" || bindPassGo == "" {
		fmt.Fprintln(os.Stderr, "環境変数 LDAP_ADMIN_PASSWORD, LDAP_BIND_PASS_NEXUS, LDAP_BIND_PASS_GO が必要です。")
		os.Exit(1)
	}

	vars := map[string]string{
		"LDAP_ROOT": baseDn,
		"LDAP_ROOT_DC": rootDc,
		"LDAP_ADMIN_PASS_HASH": generateSSHA(adminPass),
		"LDAP_BIND_PASS_NEXUS_HASH": generateSSHA(bindPassNexus),
		"LDAP_BIND_PASS_GO_HASH": generateSSHA(bindPassGo),
	}

	for _, key := range additionalEnvKeys {
		vars[key] = os.Getenv(key)
	}

	return vars
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
		fmt.Fprintln(os.Stderr, "JSON marshal error:", err)
		os.Exit(1)
	}
	return bytes
}

func loadJsonBytes() []byte {
	bytes, err := os.ReadFile(savedEnvVarsJsonFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		fmt.Fprintln(os.Stderr, "failed to read file:", err)
		os.Exit(1)
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

func makeSubstitutedFile(envVars map[string]string, fileNames ...string) {
	var templatePath, targetPath string

	switch len(fileNames) {
	case 0:
		fmt.Fprintln(os.Stderr, "引数エラー: 操作するファイル名を出力してください。")
		os.Exit(1)
	case 1:
		templatePath = "/customized/templates/" + fileNames[0] + ".ldif.template"
		targetPath = "/customized/saved-data/init/" + fileNames[0] + ".ldif"
	case 2:
		templatePath = fileNames[0]
		targetPath = fileNames[1]
	default:
		fmt.Fprintln(os.Stderr, "引数エラー: 指定するファイル数は1つまたは2つです。")
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, templatePath + " を読み込みます。") 
	templateFileData, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "テンプレート読み込みエラー:", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, templatePath + " を読み込みました。\n環境変数を注入します。") 
	content := string(templateFileData)
	content = substituteEnv(content, envVars)

	fmt.Fprintln(os.Stdout, "環境変数を注入しました。\n" + targetPath + " を更新します。") 
	targetDirName, targetBaseName := filepath.Split(targetPath)
	// 出力先ディレクトリを自動生成（存在しない場合）
	if err := os.MkdirAll(targetDirName, 0775); err != nil {
		fmt.Fprintln(os.Stderr, "出力ディレクトリ作成エラー:", err)
		os.Exit(1)
	}
	// ファイルに書き込み
	if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
		fmt.Fprintln(os.Stderr, targetBaseName + " 書き込みエラー:", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, targetPath + " を更新しました。")
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

func substituteSlapdConfFile(envVars map[string]string) {
	filePath := "/etc/openldap/slapd.conf"

	fmt.Fprintln(os.Stdout, filePath + " を読み込みます。") 
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, filePath + " 読み込みエラー:", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, filePath + " を読み込みました。\n環境変数を注入します。")
	contentsLines := strings.Split(string(fileData), "\n")
	newContentsLines := make([]string, 0)

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

	fmt.Fprintln(os.Stdout, "環境変数を注入しました。\n" + filePath + " を更新します。")

	newContents := []byte(strings.Join(newContentsLines, "\n"))
	// ファイルに書き込み
	if err := os.WriteFile(filePath, []byte(newContents), 0644); err != nil {
		fmt.Fprintln(os.Stderr, filePath + " 書き込みエラー:", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, filePath + " を更新しました。")
}