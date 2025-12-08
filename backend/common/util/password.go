package util

import (
	"crypto/rand"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func CreatePasswordHash(password string, salt []byte) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{string(hashedPassword), string(salt)}, "."), nil
}

func ComparePasswordHash(password string, storedPassword string) (bool, error) {
	splittedPassword := strings.Split(storedPassword, ".")
	if len(splittedPassword) != 2 {
		return false, errors.New("invalid stored password")
	}

	salt := []byte(splittedPassword[1])

	hashedPassword, err := CreatePasswordHash(password, salt)
	if err != nil {
		return false, err
	}

	return storedPassword == hashedPassword, nil
}
