package customError

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"
)

func Exit1(args ...any) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

// LdapExitCodeUnit はLDAP結果コードの意味とHTTPレスポンスコードを保持する
type LdapExitCodeUnit struct {
	Mean         string
	ResponseCode int
}

// LDAP結果コードとHTTPレスポンスコードの対応表
var ldapExitCodeMap = map[int]LdapExitCodeUnit{
	1:  {"Operations error", http.StatusInternalServerError},
	2:  {"Protocol error", http.StatusInternalServerError},
	3:  {"Time limit exceeded", http.StatusGatewayTimeout},
	4:  {"Size limit exceeded", http.StatusRequestEntityTooLarge},
	7:  {"Auth method not supported", http.StatusUnauthorized},
	8:  {"Stronger auth required", http.StatusUnauthorized},
	16: {"No such attribute", http.StatusBadRequest},
	17: {"Undefined attribute type", http.StatusBadRequest},
	19: {"Constraint violation", http.StatusBadRequest},
	20: {"Attribute or value exists", http.StatusConflict},
	32: {"No such object", http.StatusNotFound},
	34: {"Invalid DN syntax", http.StatusBadRequest},
	48: {"Inappropriate authentication", http.StatusUnauthorized},
	49: {"Invalid credentials", http.StatusUnauthorized},
	50: {"Insufficient access rights", http.StatusForbidden},
	51: {"Server busy", http.StatusServiceUnavailable},
	52: {"Unavailable", http.StatusServiceUnavailable},
	53: {"Unwilling to perform", http.StatusForbidden},
	64: {"Naming violation", http.StatusBadRequest},
	65: {"Object class violation", http.StatusBadRequest},
	66: {"Not allowed on non-leaf", http.StatusBadRequest},
	67: {"Not allowed on RDN", http.StatusBadRequest},
	68: {"Entry already exists", http.StatusConflict},
	69: {"Object class mods prohibited", http.StatusForbidden},
	80: {"Other error", http.StatusInternalServerError},
	81: {"Server down", http.StatusServiceUnavailable},
	82: {"Local error", http.StatusInternalServerError},
	83: {"Encoding error", http.StatusInternalServerError},
	84: {"Decoding error", http.StatusInternalServerError},
	85: {"Timeout", http.StatusGatewayTimeout},
	89: {"Parameter error", http.StatusBadRequest},
	91: {"Connect error", http.StatusServiceUnavailable},
	112: {"TLS not supported", http.StatusUpgradeRequired},
}

// GetLdapResult は LDAP コマンド実行結果のエラーを解析し、
// LDAPコードの意味と HTTP ステータスコードを返す。
// errorMessage は即時ログ出力される。
func GetLdapResult(err error, errorMessage string) (mean string, responseCode int) {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			code := status.ExitStatus()
			unit, exists := ldapExitCodeMap[code]
			if !exists {
				unit = LdapExitCodeUnit{"Unknown LDAP result code", http.StatusInternalServerError}
			}

			fmt.Fprintf(os.Stderr, "[LDAP] code=%d mean=%s http=%d msg=%s\n", code, unit.Mean, unit.ResponseCode, errorMessage)
			return unit.Mean, unit.ResponseCode
		}
	}

	// 想定外エラーの処理
	debug := os.Getenv("IS_DEBUG")
	if debug == "" {
		fmt.Fprintf(os.Stderr, "warn: unexpected error type %T: %v | msg=%s\n", err, err, errorMessage)
		return "Unexpected error", http.StatusInternalServerError
	} else {
		if debug != "true" {
			fmt.Fprintf(os.Stderr, "fatal: invalid IS_DEBUG value (%s) | msg=%s\n", debug, errorMessage)
		}
		fmt.Fprintf(os.Stderr, "fatal: unexpected error type %T: %v | msg=%s\n", err, err, errorMessage)
		os.Exit(255)
	}

	return "Unreachable", http.StatusInternalServerError
}