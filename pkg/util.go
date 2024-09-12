package pkg

import (
	"math/rand"
	"strconv"
	"strings"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func Strikethrough(input string) (text string) {
	return "\u0336" + strings.Join(strings.Split(input, ""), "\u0336")

}

func RandomString(length int) (text string) {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func RemoveZeroes(value string) (ret int, err error) {
	intPart := strings.Split(value, ".")[0]
	ret, err = strconv.Atoi(intPart)
	return
}
