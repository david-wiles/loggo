package loggo

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"
)

type LogLevel int

const (
	LogLevelInfo  LogLevel = 0
	LogLevelWarn  LogLevel = 1
	LogLevelError LogLevel = 2
	LogLevelOff   LogLevel = 3
)

// LogMiddlewareFunc is a function that acts on a request and response after the server
// has acted on the request. The LoggedResponse is not the response itself, but just a copy of it
type LogMiddlewareFunc func(LoggedResponse, *http.Request)

// Loggo is just an io.Writer and an associated LogLevel setting
// The functions defined for Loggo are primarily just for improving formatting and ergonomics,
// however, the package does include some useful functions for using Loggo as middleware in
// http server applications.
type Loggo struct {
	writer io.Writer
	level  LogLevel
}

func NewLoggo(w io.Writer, level LogLevel) *Loggo {
	return &Loggo{w, level}
}

func (log *Loggo) Cleanup() error {
	// If our writer implements io.Closer, then we should
	// close the writer when panicking or exiting the program
	if w, ok := log.writer.(io.Closer); ok {
		return w.Close()
	}
	return nil
}

// LogHandler will allow the response to be recorded and acted on after the specified http.Handler
// has completed. The response is passed to logFunc as a LoggedResponse. The log instance should be
// defined in the closure of the LogMiddlewareFunc
func (log *Loggo) LogHandler(next http.Handler, logFunc LogMiddlewareFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if log.level > 0 {
			// If the log level is not set to include info messages
			// the logger will just pass all requests to the next
			// middleware in the chain
			next.ServeHTTP(w, r)
		} else {
			var buf bytes.Buffer
			loggedWriter := newLoggingRW(w, &buf)
			next.ServeHTTP(loggedWriter, r)
			logFunc(LoggedResponse{
				Body:       &buf,
				StatusCode: loggedWriter.status,
				Header:     loggedWriter.Header().Clone(),
			}, r)
		}
	})
}

func (log Loggo) LogHandleFunc(next http.HandlerFunc, logFunc LogMiddlewareFunc) http.Handler {
	return log.LogHandler(http.HandlerFunc(next), logFunc)
}

func (log Loggo) Fatal(message string) {
	_, _ = log.writer.Write([]byte("[" + timeNow() + "] FATAL: " + message + "\n"))
	panic(message)
}

func (log Loggo) Error(message string) {
	if log.level < 3 {
		_, _ = log.writer.Write([]byte("[" + timeNow() + "] ERROR: " + message + "\n"))
	}
}

func (log Loggo) Warn(message string) {
	if log.level < 2 {
		_, _ = log.writer.Write([]byte("[" + timeNow() + "] WARN: " + message + "\n"))
	}
}

func (log Loggo) Info(message string) {
	if log.level < 1 {
		_, _ = log.writer.Write([]byte("[" + timeNow() + "] INFO: " + message + "\n"))
	}
}

// Convenience method for a common use-case for log middleware: logging requests
// Each request will be logged using the specified format:
// $remote_addr [$time_local] "$request" $path $status $http_user_agent $request_time
func (log Loggo) LogAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record time that the request was received
		start := time.Now()
		log.LogHandler(next, func(w LoggedResponse, r *http.Request) {
			remoteAddr := w.Header.Get("X-Forwarded-For")
			userAgent := w.Header.Get("User-Agent")
			timing := int(time.Now().Sub(start).Milliseconds())
			_, _ = log.writer.Write([]byte("[" + timeNow() + "] " + remoteAddr + " " + r.Method + " " + r.URL.Path + " " + strconv.Itoa(w.StatusCode) + " " + userAgent + " " + strconv.Itoa(timing) + "ms\n"))
		}).ServeHTTP(w, r)
	})
}

func (log Loggo) LogError(err error) {
	log.Error(err.Error())
}

func timeNow() string {
	return time.Now().Format(time.RFC3339Nano)
}
