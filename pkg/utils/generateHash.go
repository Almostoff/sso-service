package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateHash() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	b := make([]rune, 56)
	for i := range b {
		b[i] = letterRunes[generator.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateCode() string {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	code := generator.Intn(900000) + 100000
	return strconv.Itoa(code)
}
