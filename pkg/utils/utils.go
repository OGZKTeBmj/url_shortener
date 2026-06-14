package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func ErrWrap(msg any, err error) error {
	if err != nil {
		return fmt.Errorf("%s : %w", msg, err)
	}
	return nil
}

func GenerateShortCode() (string, error) {
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	code, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 6)
	if err != nil {
		return "", err
	}
	return code, nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func HashRefreshToken(token string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(token))

	sum := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(sum)
}
