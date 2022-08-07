package util

import (
	"database/sql"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomTitle() sql.NullString {
	return sql.NullString{String: RandomString(6), Valid: true}
}

func RandomDescription() sql.NullString {
	return sql.NullString{String: RandomString(15), Valid: true}
}
