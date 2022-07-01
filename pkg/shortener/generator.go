package shortener

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/itchyny/base58-go"
)

// sha256Of uses `SHA256` to hash the input string
func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

// base58Encoded encode binary to string using `Base58`
func base58Encoded(bytes []byte) (string, error) {
	encoded, err := base58.BitcoinEncoding.Encode(bytes)
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

// ShortenLink generates a short eight character long
// representation of input
func Shorten(input string) (string, error) {
	hashBytes := sha256Of(input)
	asNumber := new(big.Int).SetBytes(hashBytes).Uint64()
	base58Encoded, err := base58Encoded([]byte(fmt.Sprintf("%d", asNumber)))

	return base58Encoded[:8], err
}
