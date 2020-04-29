package main

import "path/filepath"

func getCertificates() (certFile string, keyFile string) {
	certFile, _ = filepath.Abs(collectorCertFile)
	keyFile, _ = filepath.Abs(collectorKeyFile)
	return
}
