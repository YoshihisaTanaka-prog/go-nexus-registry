package main

import (
	"fmt"
	"os"
	"web_app/api"
	"web_app/customError"
	"web_app/dbClient"
	"web_app/ldap"
)

var envKeys = []string{
	"GO_MANAGER_HOST_NAME",
	"GO_MANAGER_JWT_SECRET",
	"LDAP_BIND_CN_GO",
	"LDAP_DOMAIN",
}

func main() {
	if isEnvKeysEmpty() {
		customError.Exit1("以上の環境変数が指定されていないので終了します。")
	}

	ldap.InitLdap()
	
	dbClient.InitDb()
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
