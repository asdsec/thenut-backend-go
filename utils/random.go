package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const numbers = "0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min to max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generate a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// RandomImageUrl generates a random image url
func RandomImageUrl() string {
	return fmt.Sprintf("users/%s/avatar.png", RandomString(6))
}

// RandomPhoneNumber generates a random phone number
func RandomPhoneNumber() string {
	var sb strings.Builder
	k := len(numbers)

	sb.WriteString("+905")

	for i := 0; i < 9; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomGender generates a random gender `m`, or `f`
func RandomGender() string {
	var sb strings.Builder
	genders := "mf"

	sb.WriteByte(genders[rand.Intn(2)])

	return sb.String()
}

// RandomBirthDate generates a random birth date
func RandomBirthDate() time.Time {
	return time.Now()
}
