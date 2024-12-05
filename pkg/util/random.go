package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const number = "1234567890"

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomInt(min, max int64) int64 {
	return min + rng.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rng.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomNumberString(n int) string {
	var sb strings.Builder
	k := len(number)

	for i := 0; i < n; i++ {
		c := number[rng.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomEmail() string {
	return RandomString(6) + "@gmail.com"
}

func RandomPhoneNumber() string {
	return "+234" + RandomNumberString(8)
}

func GenerateAccountNumber() string {
	randomPart := rng.Int63n(1e10)
	return fmt.Sprintf("%010d", randomPart)
}
