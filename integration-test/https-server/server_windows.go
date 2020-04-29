package main

import "path/filepath"

func getCertificates() (certFile string, keyFile string) {
	certFile, _ = filepath.Abs(filepath.Join("https-server", certFile))
	keyFile, _ = filepath.Abs(filepath.Join("https-server", keyFile))
	return
}
