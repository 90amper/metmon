package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/90amper/metmon/internal/logger"
)

var sugar = logger.NewDebugLogger()

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func Logger(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		duration := time.Since(start)
		buf, _ := io.ReadAll(r.Body)

		reader := io.NopCloser(bytes.NewBuffer(buf))
		r.Body = reader

		h.ServeHTTP(&lw, r)

		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			// "body", spew.Sprintf("%#v", buf),
			"body", fmt.Sprint(strings.ReplaceAll(string(buf), "\"", "")),
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}
	return http.HandlerFunc(logFn)
}
