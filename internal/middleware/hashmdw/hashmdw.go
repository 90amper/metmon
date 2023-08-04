package hashmdw

import (
	"bytes"
	"crypto/hmac"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/90amper/metmon/internal/logger"
	"github.com/90amper/metmon/pkg/hasher"
)

type (
	Hasher struct {
		HashKey string
		HashAlg string
	}
)

func Init(key, alg string) *Hasher {
	return &Hasher{
		HashKey: key,
		HashAlg: alg,
	}
}

func (hs *Hasher) HashMiddleware(h http.Handler) http.Handler {
	hashFn := func(w http.ResponseWriter, r *http.Request) {
		rHashB64 := r.Header.Get("HashSHA256")
		if rHashB64 != "" {
			rHash, err := base64.StdEncoding.DecodeString(rHashB64)
			if err != nil {
				logger.Error(err)
			}
			buf, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(buf))

			bodyHash, err := hasher.Hash(buf, []byte(hs.HashKey), hs.HashAlg)
			if err != nil {
				logger.Error(err)
			}
			if !hmac.Equal(rHash, bodyHash) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			buf, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(buf))
			bodyHash, err := hasher.Hash(buf, []byte(hs.HashKey), hs.HashAlg)
			bodyHashB64 := base64.StdEncoding.EncodeToString(bodyHash)
			if err != nil {
				logger.Error(err)
			}
			r.Header.Set("HashSHA256", string(bodyHashB64))
		}
		h.ServeHTTP(w, r)

	}
	return http.HandlerFunc(hashFn)
}
