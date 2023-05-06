package util

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
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

func RandomUUID() uuid.UUID {
	uuid := uuid.New()
	return uuid
	// return uuid.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomNullableString() sql.NullString {
	return sql.NullString{
		String: RandomString(20),
		Valid:  true,
	}
}

func RandomNullableUUID() uuid.NullUUID {
	return uuid.NullUUID{
		UUID:  uuid.Nil,
		Valid: false,
	}
}
