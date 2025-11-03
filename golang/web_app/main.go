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
	"NEXUS_ADMIN_PASSWORD",
	"NEXUS_ADMIN_USERNAME",
	"NEXUS_EXPOSED_HOST_NAME",
	"POSTGRES_PASSWORD",
	"POSTGRES_USER",
}


func main() {
	if isEnvKeysEmpty() {
		customError.Exit1("以上の環境変数が指定されていないので終了します。")
	}

	if os.Getenv("NEPLUS_ADMIN_USERNAME") == os.Getenv("NEXUS_ADMIN_USERNAME") && os.Getenv("NEPLUS_ADMIN_PASSWORD") != os.Getenv("NEXUS_ADMIN_PASSWORD") {
		customError.Exit1("管理者権限の環境変数エラー:\nNexusとNePlusの管理者アカウント名は同一ですが、パスワードが異なります。\n")
	}

	cryption.InitCription()

	apiPasswords := map[string]string{
		"admin": rand.Text(),
	}

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
