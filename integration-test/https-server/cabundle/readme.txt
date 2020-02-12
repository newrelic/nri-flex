Certificates in this file are created for a host named `https-server`, so this must be run inside a docker
compose file, as well as the tests that make use of HTTPS endpoints using this certificate.

To regenerate them, you can run the following command from your $GOROOT/src/lib/crypto/tls folder:

go run generate_cert.go  --rsa-bits 1024 --host https-server --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h
