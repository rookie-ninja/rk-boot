// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// read cert
	caCert, err := ioutil.ReadFile(path.Join(wd, "example/grpc/server/cert/ca.pem"))
	if err != nil {
		log.Fatal(err)
	}

	// create cert pool
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// init http client with cert
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	// send request to swagger
	_, err = client.Get("https://localhost:8081/sw")
	if err != nil {
		log.Fatal(err)
	}
}
