package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// The fake collector is a simple https service that ingests the metrics from the agent, and enables extra
// endpoints to be controlled and monitored from the tests.
// It stores in a queue all the events that it receives from the agent.
func main() {
	logrus.Info("Running fake HTTP server...")

	mux := http.NewServeMux()
	mux.HandleFunc("/nginx_status", func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte(`Active connections: 43
server accepts handled requests
8000 7368 10993
Reading: 0 Writing: 5 Waiting: 38
`))
	})
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != nil {
		logrus.WithError(err).Error("Running fake http server")
	}
}
