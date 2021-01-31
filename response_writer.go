package loggo

import (
	"bytes"
	"io"
	"net/http"
)

// loggingRW is used internally to record the response written to a
// http.ResponseWriter so that it can be used later for logging or caching.
// The response is exposed through a new LoggedResponse for each request.
// Use LogHandler to access the LoggedResponse
type loggingRW struct {
	w      http.ResponseWriter
	status int
	writer io.Writer
}

// LoggedResponse represents a finished response as recorded by a loggingRW
// It contains a copy of the headers from the original http.ResponseWriter,
// the status code, and a buffer containing a copy of the bytes written to
// the original request
type LoggedResponse struct {
	Body       *bytes.Buffer
	StatusCode int
	Header     http.Header
}

// Create a new loggingRW using the given ResponseWriter
func newLoggingRW(w http.ResponseWriter, buf *bytes.Buffer) *loggingRW {
	return &loggingRW{
		w:      w,
		status: 200,
		writer: io.MultiWriter(w, buf),
	}
}

func (l *loggingRW) Header() http.Header {
	return l.w.Header()
}

func (l *loggingRW) Write(b []byte) (int, error) {
	return l.writer.Write(b)
}

func (l *loggingRW) WriteHeader(statusCode int) {
	l.w.WriteHeader(statusCode)
	l.status = statusCode
}
