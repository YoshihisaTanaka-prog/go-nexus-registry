package main

import (
	"context"
	"fmt"
	"os"
	"web_app/api"
	"web_app/cryption"
	"web_app/customError"
	"web_app/dbClient"
	"web_app/ldap"
)

var envKeys = []string{
	"GO_MANAGER_COOKIE_SECRET",
	"GO_MANAGER_HOST_NAME",
	"GO_MANAGER_KEY",
	"LDAP_ADMIN_PASSWORD",
	"LDAP_ADMIN_USERNAME",
	"LDAP_BIND_CN_NEXUS",
	"LDAP_BIND_PASS_NEXUS",
	"LDAP_DOMAIN",
	"LDAP_PORT",
	"NEXUS_EXPOSED_URL",
	"ROOT_DIR_PATH",
}

func main() {
	if isEnvKeysEmpty() {
		customError.Exit1("以上の環境変数が指定されていないので終了します。")
	}

	cryption.InitCription()
	
	ctx := context.Background()
	dbClient.InitDb(&ctx)
	ldap.InitLdap()

	api.Start()
}


func isEnvKeysEmpty() bool {
	result := false

	for _, key := range envKeys {
		if (isEnvKeyEmpty(key)) {
			result = true
		}
	}

	return result
}

func isEnvKeyEmpty(key string) bool {
	if len(os.Getenv(key)) == 0 {
		fmt.Fprintln(os.Stderr, "環境変数 \"" + key + "\" を指定してください。")
		return true
	} else {
		return false
	}
}
