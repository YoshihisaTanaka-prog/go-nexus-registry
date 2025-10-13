package handler

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"time"
	"web_app/customError"
)

var ONE_HOUR_SEC = 3600

var hmacSecret []byte

func InitJwt()  {
	_, err := rand.Read(hmacSecret)
	if err != nil {
		customError.Exit1("Generating JWT secret error")
	}
}

func SetCookie(c *gin.Context, key string, value string) {
	secure := true
	hostName := os.Getenv("GO_MANAGER_HOST_NAME")
	if hostName == "localhost" {
		secure = false
	}
	c.SetCookie(key, value, ONE_HOUR_SEC, "/", os.Getenv("GO_MANAGER_HOST_NAME"), secure, true)
}

func DeleteCookie(c *gin.Context, key string) {
	secure := true
	hostName := os.Getenv("GO_MANAGER_HOST_NAME")
	if hostName == "localhost" {
		secure = false
	}
	c.SetCookie(key, "", 1, "/", os.Getenv("GO_MANAGER_HOST_NAME"), secure, true)
}

type MyClaims struct {
	UserId string `json:uid`
	Roles []string `json:roles`
	jwt.RegisteredClaims
}

func CreateJWT(userName string, groups []string) (string, error) {
	// ===== 発行 =====
	now := time.Now()
	uuidV4Arr, err := uuid.NewRandom()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "", err
	}

	uuidV4 := fmt.Sprintf("%s", uuidV4Arr)

	claims := MyClaims{
		UserId: userName,
		Roles:   groups,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("GO_MANAGER_HOST_NAME"),
			Subject:   "auth",
			Audience:  []string{os.Getenv("GO_MANAGER_HOST_NAME")},
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ONE_HOUR_SEC) * time.Second)),
			NotBefore: jwt.NewNumericDate(now.Add(-30 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuidV4,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(hmacSecret)
}