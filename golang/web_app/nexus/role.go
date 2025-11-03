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

func (requestBody requestBodyType)toBuffer() *bytes.Buffer {
	bodyBytes, err:= json.Marshal(requestBody)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Stringifying Model Error:", err)
		return bytes.NewBuffer([]byte{})
	}
	return bytes.NewBuffer(bodyBytes)
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

	req, err := http.NewRequest("POST", url, requestBody.toBuffer())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Creating Role in Nexus Error", err)
		return false
	}

	req.SetBasicAuth(nexusAdminConfig.Username, apiPasswords["admin"])
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Creating Role in Nexus Error", err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode < 400
}

func FindRole(id string) (requestBody *requestBodyType, ok bool) {
	url := fmt.Sprintf("%s/service/rest/v1/security/roles/%s", nexusAdminConfig.URL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Searching Role in Nexus Error", err)
		return nil, false
	}

	req.SetBasicAuth(nexusAdminConfig.Username, apiPasswords["admin"])
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Searching Role in Nexus Error", err)
		return nil, false
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read response body:", err)
		return nil, false
	}

	if resp.StatusCode == 404 {
		return nil, true
	}

	if resp.StatusCode > 399 {
		fmt.Fprintln(os.Stderr, "Searching Role in Nexus Error:", string(bodyBytes))
		return nil, false
	}

	requestBody = &requestBodyType{}
	if err := json.Unmarshal(bodyBytes, requestBody); err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse response body:", err)
		return nil, false
	}
	return requestBody, true
}

func updateRoleUnit(requestBody requestBodyType) bool {
	currentRequestBody, ok := FindRole(requestBody.Id)

	if !ok {
		return false
	}

	if currentRequestBody == nil {
		if !CreateRole(requestBody.Id, requestBody.Name) {
			return false
		}
	} else {
		if requestBody.Name == "" {
			requestBody.Name = (*currentRequestBody).Name
			requestBody.Description = (*currentRequestBody).Name
		}

		if len(requestBody.Privileges) == 0 {
			requestBody.Privileges = (*currentRequestBody).Privileges
		}
	}

	url := fmt.Sprintf("%s/service/rest/v1/security/roles/%s", nexusAdminConfig.URL, requestBody.Id)

	req, err := http.NewRequest("PUT", url, requestBody.toBuffer())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Updating Role in Nexus Error", err)
		return false
	}

	req.SetBasicAuth(nexusAdminConfig.Username, apiPasswords["admin"])
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Updating Role in Nexus Error", err)
		return false
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read response body:", err)
		return false
	}

	if resp.StatusCode > 399 {
		fmt.Fprintln(os.Stderr, "Updating Role in Nexus Error:", string(bodyBytes))
		return false
	}

	return true
}

func UpdateRoleName(id string, name string) bool {
	requestBody := requestBodyType{
		Id:          id,
		Name:        name,
		Description: name,
		Privileges:  []string{},
		Roles:       []string{},
	}
	return updateRoleUnit(requestBody)
}

func UpdateRolePrivileges(id string, privileges []string) bool {
	requestBody := requestBodyType{
		Id:          id,
		Name:        "",
		Description: "",
		Privileges:  privileges,
		Roles:       []string{},
	}
	return updateRoleUnit(requestBody)
} 