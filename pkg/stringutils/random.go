package stringutils

import (
	"math/rand"
	"strings"
	"time"
)

// RandomString : Generates random strings of given length with Uppercase charset
func RandomString(size int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var b strings.Builder
	for i := 0; i < size; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()

	return str
}
