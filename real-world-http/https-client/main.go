package main

import (
	"crypto/tls"
	// "crypto/x509"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {
	clientCert := os.Getenv("CLIENT_CERT")
	clientKey := os.Getenv("CLIENT_KEY")
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		panic(err)
	}

	// certPool := x509.NewCertPool()
	// certPool.AppendCertsFromPEM(cert)
	// tlsConfig := &tls.Config{
	// 	RootCAs: certPool,
	// }
	// tlsConfig.BuildNameToCertificate()

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}
	resp, err := client.Get("https://localhost:18443")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}
