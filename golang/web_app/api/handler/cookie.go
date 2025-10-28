package handler

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"strings"
	"time"
	"web_app/cryption"
	"web_app/customError"
)

var ONE_HOUR_SEC = 3600

var (
	hmacSecret = make([]byte, 64)
	hostName string
	isLocalHost = false
)


func InitJwt()  {
	var err error
	hmacSecret, err =  base64.StdEncoding.DecodeString(os.Getenv("GO_MANAGER_COOKIE_SECRET"))
	if err != nil {
		customError.Exit1("Decoding JWT secret error")
	}
	hostName := os.Getenv("GO_MANAGER_HOST_NAME")
	if hostName == "localhost" {
		isLocalHost = true
	} else if strings.HasPrefix("localhost:") {
		isLocalHost = true
	}
}

func setCookie(c *gin.Context, key string, value string) {
	c.SetCookie(key, value, ONE_HOUR_SEC, "/", hostName, !isLocalHost, true)
}

func deleteCookie(c *gin.Context, key string) {
	c.SetCookie(key, "", -1, "/", hostName, !isLocalHost, true)
}

type MyClaims struct {
	UserId string `json:uid`
	Roles []string `json:roles`
	jwt.RegisteredClaims
}

func createJwt(userName string, groups []string) (string, error) {
	// ===== 発行 =====
	now := time.Now()
	uuidV4Arr, err := uuid.NewRandom()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "", err
	}

	uuidV4 := fmt.Sprintf("%s", uuidV4Arr)

	userId, err := cryption.Encrypt(userName)

	if err != nil {
		return "", err
	}

	claims := MyClaims{
		UserId: userId,
		Roles:   groups,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("GO_MANAGER_HOST_NAME"),
			Subject:   "auth",
			Audience:  []string{os.Getenv("GO_MANAGER_HOST_NAME") + "/app"},
			ExpiresAt: jwt.NewNumericDate(now.Add(3600 * time.Second)),
			NotBefore: jwt.NewNumericDate(now.Add(-30 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuidV4,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(hmacSecret)
}

func parseJwt(raw string) (userId string, roles []string) {
	parsed, err := jwt.ParseWithClaims(raw, &MyClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 { return nil, fmt.Errorf("bad alg") }
		return hmacSecret, nil
	}, jwt.WithIssuer(os.Getenv("GO_MANAGER_HOST_NAME")), jwt.WithAudience(os.Getenv("GO_MANAGER_HOST_NAME") + "/app"))
	
	if err != nil || !parsed.Valid {
		return "", []string{}
	}

	userName, err := cryption.Decrypt(parsed.Claims.(*MyClaims).UserId)

	if err != nil {
		return "", []string{}
	}
	return userName, parsed.Claims.(*MyClaims).Roles
}