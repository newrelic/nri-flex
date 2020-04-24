package main

import (
	"crypto/tls"
	"net/http"

	"github.com/sirupsen/logrus"
)

const collectorCertFile = "/cabundle/cert.pem"
const collectorKeyFile = "/cabundle/key.pem"

// The fake collector is a simple https service that ingests the metrics from the agent, and enables extra
// endpoints to be controlled and monitored from the tests.
// It stores in a queue all the events that it receives from the agent.
func main() {
	logrus.Info("Running fake HTTPS server...")

	mux := http.NewServeMux()
	mux.HandleFunc("/nginx_status", func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte(`Active connections: 43
server accepts handled requests
8000 7368 10993
Reading: 0 Writing: 5 Waiting: 38
`))
	})
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	srv := &http.Server{
		Addr:         ":8043",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	if err := srv.ListenAndServeTLS(collectorCertFile, collectorKeyFile); err != nil {
		logrus.WithError(err).Error("Running fake https server")
	}
}
