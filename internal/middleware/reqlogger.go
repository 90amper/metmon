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
		body   bytes.Buffer
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.body.Write(b)
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
		// wbuf, _ := io.ReadAll(w.Body)
		reader := io.NopCloser(bytes.NewBuffer(buf))
		// writer := io.NopCloser(bytes.NewBuffer(wbuf))
		r.Body = reader

		h.ServeHTTP(&lw, r)

		sugar.Infoln(
			"_uri:", r.RequestURI,
			"_method:", r.Method,
			// "body", spew.Sprintf("%#v", buf),
			// "_req_body:", fmt.Sprint(strings.ReplaceAll(string(buf), "\"", "")),
			// "_resp_body:", fmt.Sprint(strings.ReplaceAll(strings.Join(strings.Fields(responseData.body.String()), " "), "\"", "")),
			"_status:", responseData.status,
			"_duration:", duration,
			"_size:", responseData.size,
		)

		sugar.Debugln(
			"_req_body:", fmt.Sprint(strings.ReplaceAll(string(buf), "\"", "")),
			"_resp_body:", fmt.Sprint(strings.ReplaceAll(strings.Join(strings.Fields(responseData.body.String()), " "), "\"", "")),
		)
	}
	return http.HandlerFunc(logFn)
}
