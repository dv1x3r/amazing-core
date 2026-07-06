package pass

import (
	"bytes"
	"crypto/pbkdf2"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func getRandomSalt(size int) []byte {
	const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	salt := make([]byte, size)
	l := len(allowedChars)
	for i := range salt {
		salt[i] = allowedChars[rand.Intn(l)]
	}
	return salt
}

func MakePbkdf2(password string) (string, error) {
	salt := getRandomSalt(22)
	iter := 720000 // https://docs.djangoproject.com/en/5.0/releases/5.0/#django-contrib-auth
	dk, err := pbkdf2.Key(sha256.New, password, salt, iter, sha256.Size)
	if err != nil {
		return "", fmt.Errorf("invalid key: %w", err)
	}
	b64Hash := base64.StdEncoding.EncodeToString(dk)
	return fmt.Sprintf("pbkdf2_sha256$%d$%s$%s", iter, salt, b64Hash), nil
}

func CheckPbkdf2(password string, encoded string) (bool, error) {
	split := strings.SplitN(encoded, "$", 4)
	if len(split) != 4 {
		return false, fmt.Errorf("invalid number of segments: %d", len(split))
	}

	hasher := split[0]
	if hasher != "pbkdf2_sha256" {
		return false, fmt.Errorf("invalid hasher: %s", hasher)
	}

	iter, err := strconv.Atoi(split[1])
	if err != nil {
		return false, fmt.Errorf("invalid integer: %w", err)
	}

	salt := []byte(split[2])

	k, err := base64.StdEncoding.DecodeString(split[3])
	if err != nil {
		return false, fmt.Errorf("invalid encoding: %w", err)
	}

	dk, err := pbkdf2.Key(sha256.New, password, salt, iter, sha256.Size)
	if err != nil {
		return false, fmt.Errorf("invalid key: %w", err)
	}
	return bytes.Equal(k, dk), nil
}
