package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alpa = "abcdfghjklmnpqrstvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}
func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(alpa[rand.Intn(len(alpa))])
	}
	return sb.String()
}
func RandomOwner() string {
	return RandomString(8)
}
func RandomMoney() int64 {
	return randomInt(0, 1000)
}

func RandomCurrency() string {
	c := []string{"GHC", "USD", "EUR"}
	n := len(c)

	return c[rand.Intn(n)]
}
func RandomEmail() string {
	return fmt.Sprintf("%s@example.com", RandomOwner())
}
