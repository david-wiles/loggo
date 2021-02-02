package main

import (
	"github.com/david-wiles/loggo"
	"net/http"
	"os"
	"strconv"
)

func main() {

	// Create a new Loggo instance
	log := loggo.NewLoggo(os.Stdout, 0)

	// Log message with varying log levels
	log.Info("Hello, World!")
	log.Error("Oh no!")
	log.Warn("Hmm...")

	// Create another Loggo that'll ignore everything except errors
	noLoggo := loggo.NewLoggo(os.Stdout, loggo.LogLevelError)

	noLoggo.Info("This won't print")
	noLoggo.Error("This will print")

	// You can also use files for log output
	if f, err := os.Create("log.log"); err == nil {
		fileLoggo := loggo.NewLoggo(f, 0)
		fileLoggo.Info("Hello from file loggo")
		// The file will be closed when cleaned up
		err = fileLoggo.Cleanup()
	}

	// Create logging middleware
	logMiddleware := log.LogHandleFunc(func(w http.ResponseWriter, r *http.Request) {
		// Your http.Handler goes here...
		_, _ = w.Write([]byte("Hello, World!"))
	}, func(resp loggo.LoggedResponse, r *http.Request) {
		// And the recorded response is available here
		log.Info("Got request: " + r.URL.String() + " Response: " + strconv.Itoa(resp.StatusCode) + " " + resp.Body.String())
	})

	_ = http.ListenAndServe(":4321", logMiddleware)
}
