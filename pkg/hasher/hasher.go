package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func Hash(value, key []byte, alg string) ([]byte, error) {
	switch alg {
	case "SHA256":
		// подписываем алгоритмом HMAC, используя SHA-256
		h := hmac.New(sha256.New, key)
		h.Write(value)
		dst := h.Sum(nil)
		// logger.Debug("%v", dst)
		return dst, nil
	default:
		return nil, fmt.Errorf("invalid argument")
	}

}
