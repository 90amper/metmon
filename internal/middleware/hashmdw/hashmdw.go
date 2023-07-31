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

	// hashRespData struct {
	// 	body   bytes.Buffer
	// 	status int
	// }

	// hashRespWriter struct {
	// 	http.ResponseWriter
	// 	hashRespData *hashRespData
	// }
)

func Init(key, alg string) *Hasher {
	return &Hasher{
		HashKey: key,
		HashAlg: alg,
	}
}

// func (r *hashRespWriter) Write(b []byte) (int, error) {
// 	size, err := r.ResponseWriter.Write(b)
// 	r.hashRespData.body.Write(b)
// 	return size, err
// }

// func (r *hashRespWriter) WriteHeader(statusCode int) {
// 	r.ResponseWriter.WriteHeader(statusCode)
// 	r.hashRespData.status = statusCode
// }

func (hs *Hasher) HashMiddleware(h http.Handler) http.Handler {
	hashFn := func(w http.ResponseWriter, r *http.Request) {
		// hashRespData := &hashRespData{}
		// hw := hashRespWriter{
		// 	ResponseWriter: w,
		// 	hashRespData:   hashRespData,
		// }
		rHashB64 := r.Header.Get("HashSHA256")
		// definedHash := rHash != ""
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
			// w.Header().Set("HashSHA256", fmt.Sprintf("%s", bodyHash))
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

// func tmpGzipMiddleware(h http.Handler) http.Handler {
// 	gzipFn := func(w http.ResponseWriter, r *http.Request) {
// 		ow := w

// 		acceptEncoding := r.Header.Get("Accept-Encoding")
// 		supportsGzip := strings.Contains(acceptEncoding, "gzip")
// 		if supportsGzip {
// 			w.Header().Set("Content-Encoding", "gzip")
// 			cw := compressor.CompressWriter(w)
// 			ow = cw
// 			defer cw.Close()
// 		}

// 		contentEncoding := r.Header.Get("Content-Encoding")
// 		sendsGzip := strings.Contains(contentEncoding, "gzip")
// 		if sendsGzip {
// 			cr, err := compressor.CompressReader(r.Body)
// 			if err != nil {
// 				w.WriteHeader(http.StatusInternalServerError)
// 				return
// 			}
// 			r.Body = cr
// 			defer cr.Close()
// 		}

// 		h.ServeHTTP(ow, r)
// 	}
// 	return http.HandlerFunc(gzipFn)
// }

// func tmpLogger(h http.Handler) http.Handler {
// 	logFn := func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()

// 		responseData := &responseData{
// 			status: 0,
// 			size:   0,
// 		}
// 		lw := loggingResponseWriter{
// 			ResponseWriter: w,
// 			responseData:   responseData,
// 		}

// 		duration := time.Since(start)
// 		buf, _ := io.ReadAll(r.Body)
// 		reader := io.NopCloser(bytes.NewBuffer(buf))
// 		r.Body = reader

// 		h.ServeHTTP(&lw, r)

// 		sugar.Infoln(
// 			"_uri:", r.RequestURI,
// 			"_method:", r.Method,
// 			// "body", spew.Sprintf("%#v", buf),
// 			"_status:", responseData.status,
// 			"_duration:", duration,
// 			"_size:", responseData.size,
// 		)

// 		sugar.Debugln(
// 			"_req_body:", fmt.Sprint(strings.ReplaceAll(string(buf), "\"", "")),
// 			"_resp_body:", fmt.Sprint(strings.ReplaceAll(strings.Join(strings.Fields(responseData.body.String()), " "), "\"", "")),
// 		)
// 	}
// 	return http.HandlerFunc(logFn)
// }
