package utils

import (
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
