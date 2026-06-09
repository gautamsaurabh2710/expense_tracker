package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP() string {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "123456"
	}
	return fmt.Sprintf("%06d", n.Int64()+100000)
}
