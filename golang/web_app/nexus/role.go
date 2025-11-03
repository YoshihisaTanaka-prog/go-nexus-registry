package nexus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)
var nexusAdminConfig = struct {
	URL        string
	Username   string
} {
	"http://nexus:8081",
	"_api_admin",
}

type requestBodyType struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Privileges  []string `json:"privileges"`
	Roles       []string `json:"roles"`
}

func (requestBody requestBodyType)toByteSlice() []byte {
	bytes, err:= json.Marshal(requestBody)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Stringifying Model Error:", err)
		return []byte{}
	}
	return bytes
}

func CreateRole(id string, name string) bool {
	url := fmt.Sprintf("%s/service/rest/v1/security/roles", nexusAdminConfig.URL)

	requestBody := requestBodyType{
		Id:            id,
		Name:          name,
		Description:   name,
		Privileges:  []string{},
		Roles:       []string{},
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody.toByteSlice()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Creating Role in LDAP Error", err)
		return false
	}

	req.SetBasicAuth(nexusAdminConfig.Username, apiPasswords["admin"])
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Creating Role in LDAP Error", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return false
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read response body:", err)
		return false
	}
	fmt.Fprintln(os.Stdout, resp.StatusCode, string(bodyBytes))
	return true
}

func SearchRole() {}

func updateRoleUnit(requestBody requestBodyType) {
}

// func UpdateRoleName(id string, name string) {
// 	url := fmt.Sprintf("%s/service/rest/v1/security/roles", nexusAdminConfig.URL)
// }

// func UpdateRolePrivileges(id string, privileges []string) {
// } 