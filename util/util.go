package util

import (
	"crypto/rand"
	"math/big"
)

const (
	randomStringPattern string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

func GenRandomString(num int) string {
	b := make([]byte, num)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(randomStringPattern))))
		b[i] = randomStringPattern[n.Uint64()]
	}
	return string(b)
}
