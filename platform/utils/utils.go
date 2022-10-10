package utils

import (
	"2f-authorization/platform/logger"
	"context"
	"crypto/rand"
	"io"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321"
const specialBytes = `!@#$%^&*:.`

func HashAndSalt(ctx context.Context, pwd []byte, logger logger.Logger) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, 14)
	if err != nil {
		logger.Error(ctx, "could not hash password", zap.Error(err))
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(hashedPwd, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPassword))
	return err == nil
}

func GenerateRandomString(length int, includeSpecial bool) string {
	str := letterBytes
	if includeSpecial {
		str += specialBytes
	}

	randString := make([]byte, length)
	_, _ = io.ReadAtLeast(rand.Reader, randString, length)
	for i := 0; i < len(randString); i++ {
		randString[i] = str[int(randString[i])%len(str)]
	}

	return string(randString)
}
