package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

var (
	certFile = "cabundle/cert.pem"
	keyFile  = "cabundle/key.pem"
)

// The fake collector is a simple https service that ingests the metrics from the agent, and enables extra
// endpoints to be controlled and monitored from the tests.
// It stores in a queue all the events that it receives from the agent.
func main() {

	useTLS, err := strconv.ParseBool(os.Args[1])
	if err != nil {
		logrus.Errorf("invalid argument: %s", err.Error())
	}

	logrus.Infof("Starting HTTP server with tls: %v", useTLS)

	mux := http.NewServeMux()
	mux.HandleFunc("/nginx_status", serveNginx)
	mux.HandleFunc("/json", serveJSON)

	if useTLS {
		startHTTPS(mux)
	} else {
		startHTTP(mux)
	}
}

func startHTTP(mux *http.ServeMux) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	defer srv.Close()

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		logrus.WithError(err).Error("failed to start http server")
	}
}

func startHTTPS(mux *http.ServeMux) {
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
		InsecureSkipVerify: true,
	}
	srv := &http.Server{
		Addr:         ":8043",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	defer srv.Close()

	certFilePath, keyFilePath := getCertificates(certFile, keyFile)
	err := srv.ListenAndServeTLS(certFilePath, keyFilePath)
	if err != http.ErrServerClosed {
		logrus.WithError(err).Error("failed to start https server")
	}
}

func serveNginx(rw http.ResponseWriter, r *http.Request) {
	_, _ = rw.Write([]byte(`Active connections: 43
	server accepts handled requests
	8000 7368 10993
	Reading: 0 Writing: 5 Waiting: 38
	`))
}

func serveJSON(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "application/json")
	_, _ = rw.Write([]byte(`
	{
		"metrics": [
			{
			 "cpu": 10.0,
			 "memory": 3500,
			 "disk": 500
			} 
		]
	}
	`))
}

// In Linux we are running in a docker container and the cabundle dir is in the root alongside the executables,
// but in windows we are running the server file with "go run" which changes the working dir to "integration-test"
// so we have to build the path from there
func getCertificates(certFile string, keyFile string) (certFilePath string, keyFilePath string) {
	if runtime.GOOS == "windows" {
		certFilePath, _ = filepath.Abs(filepath.Join("https-server", certFile))
		keyFilePath, _ = filepath.Abs(filepath.Join("https-server", keyFile))
	} else {
		certFilePath, _ = filepath.Abs(certFile)
		keyFilePath, _ = filepath.Abs(keyFile)
	}
	return
}
