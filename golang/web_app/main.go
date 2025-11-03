package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"web_app/api"
	"web_app/cryption"
	"web_app/customError"
	"web_app/dbClient"
	"web_app/ldap"
	"web_app/nexus"
)

var envKeys = []string{
	"LDAP_ADMIN_PASSWORD",
	"LDAP_ADMIN_USERNAME",
	"LDAP_BIND_CN_NEXUS",
	"LDAP_BIND_PASS_NEXUS",
	"DOMAIN_NAME",
	"ORGANISATION_NAME",
	"LDAP_PORT",
	"NEPLUS_ADMIN_USERNAME",
	"NEPLUS_ADMIN_PASSWORD",
	"NEPLUS_COOKIE_SECRET",
	"NEPLUS_HOST_NAME",
	"NEPLUS_KEY",
	"NEXUS_EXPOSED_HOST_NAME",
	"POSTGRES_PASSWORD",
	"POSTGRES_USER",
}

var apiPasswords = map[string]string{
	"admin": rand.Text(),
}

func main() {
	if isEnvKeysEmpty() {
		customError.Exit1("以上の環境変数が指定されていないので終了します。")
	}

	cryption.InitCription()
	
	ctx := context.Background()
	dbClient.InitDb(&ctx)
	ldap.InitLdap(apiPasswords)
	nexus.InitNexus(&apiPasswords)
	
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
