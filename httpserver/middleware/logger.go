package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"monitoring/service/logging"
	"time"
)

type internalLogger struct {
	http.ResponseWriter
	statusCode int
	result     []byte
}

func NewLoggingResponseWriter(rw http.ResponseWriter) *internalLogger {
	return &internalLogger{ResponseWriter: rw, statusCode: http.StatusOK, result: nil}
}

func (logger *internalLogger) WriteHeader(statusCode int) {
	logger.statusCode = statusCode
	logger.ResponseWriter.WriteHeader(statusCode)
}

func (logger *internalLogger) Write(b []byte) (int, error) {
	resultCopy := make([]byte, len(b))
	copy(resultCopy, b)
	logger.result = resultCopy
	return logger.ResponseWriter.Write(resultCopy)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		generatedId := time.Now().UnixNano() / int64(time.Millisecond)

		var buf bytes.Buffer
		tee := io.TeeReader(r.Body, &buf)
		bodyBytes, err := io.ReadAll(tee)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		bodyString := strings.ReplaceAll(strings.ReplaceAll(string(bodyBytes), " ", ""), "\r\n", "")
		logging.Httpf("-> | req_id_%d | %s | %s | body: %s ", generatedId, r.URL.Path, r.Method, bodyString)

		logger := NewLoggingResponseWriter(rw)
		next.ServeHTTP(logger, r)

		result := fmt.Sprintf("<- | req_id_%d | %s | %s | %d", generatedId, r.URL.Path, r.Method, logger.statusCode)

		reSpaces := regexp.MustCompile(`\s+`)
		reNewlines := regexp.MustCompile(`\r\n`)
		inputString := reSpaces.ReplaceAllString(string(logger.result), " ")
		s := reNewlines.ReplaceAllString(inputString, " ")

		if s != "" {
			result += fmt.Sprintf(" | body: %s", s)
		}

		logging.Httpf(result)
	})
}
