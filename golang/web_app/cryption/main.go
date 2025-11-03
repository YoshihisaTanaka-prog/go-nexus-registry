package cryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"web_app/customError"
)

var cryptionKey = make([]byte, 64)

func InitCription() {
	var err error
	cryptionKey, err = base64.StdEncoding.DecodeString(os.Getenv("NEPLUS_KEY"))
	if err != nil {
		customError.Exit1("鍵のデコードに失敗しました:", err)
	}
}

// 暗号化
func Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(cryptionKey)
	if err != nil {
		return "", fmt.Errorf("AES初期化失敗: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GCM初期化失敗: %v", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce生成失敗: %v", err)
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// 復号化
func Decrypt(cipherBase64 string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherBase64)
	if err != nil {
		return "", fmt.Errorf("暗号文のデコードに失敗しました: %v", err)
	}

	block, err := aes.NewCipher(cryptionKey)
	if err != nil {
		return "", fmt.Errorf("AES初期化失敗: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GCM初期化失敗: %v", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("データが短すぎます")
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("復号化失敗: %v", err)
	}

	return string(plainText), nil
}