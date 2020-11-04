// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_tls

import (
	"github.com/golang/glog"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	runtime2 "runtime/debug"
)

var (
	rootCAConfig = `{
    "signing": {
        "default": {
            "expiry": "43800h"
        },
        "profiles": {
            "server": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "server auth"
                ]
            },
            "client": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "client auth"
                ]
            },
            "peer": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "server auth",
                    "client auth"
                ]
            }
        }
    }
}`
	rootCACsr = `{
    "CN": "RK Demo CA",
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "L": "BJ",
            "O": "RK",
            "ST": "Beijing",
            "OU": "RK Demo"
        }
    ]
}`
	serverCAConfig = `{
    "CN": "example.net",
    "hosts": [
        "localhost",
		"127.0.0.1",
		"0.0.0.0"
    ],
    "key": {
        "algo": "ecdsa",
        "size": 256
    },
    "names": [
        {
            "C": "CN",
            "ST": "Beijing",
            "L": "BJ"
        }
    ]
}`
)

var originalWD, _ = os.Getwd()

type TlsEntry struct {
	logger       *zap.Logger
	port         uint64
	certFilePath string
	keyFilePath  string
	generateCert bool
	generatePath string
}

type TlsOption func(*TlsEntry)

func WithLogger(logger *zap.Logger) TlsOption {
	return func(entry *TlsEntry) {
		entry.logger = logger
	}
}

func WithPort(port uint64) TlsOption {
	return func(entry *TlsEntry) {
		entry.port = port
	}
}

func WithCertFilePath(subPath string) TlsOption {
	return func(entry *TlsEntry) {
		entry.certFilePath = subPath
	}
}

func WithKeyFilePath(subPath string) TlsOption {
	return func(entry *TlsEntry) {
		entry.keyFilePath = subPath
	}
}

func WithGenerateCert(generate bool) TlsOption {
	return func(entry *TlsEntry) {
		entry.generateCert = generate
	}
}

func WithGeneratePath(subPath string) TlsOption {
	return func(entry *TlsEntry) {
		wd, err := os.Getwd()
		if err != nil {
			shutdownWithError(err)
		}
		entry.generatePath = path.Join(wd, subPath)
	}
}

func NewTlsEntry(opts ...TlsOption) *TlsEntry {
	entry := &TlsEntry{
		logger:       zap.NewNop(),
		generateCert: false,
	}

	for i := range opts {
		opts[i](entry)
	}

	// generate tls cert files with cfssl
	// be sure cfssl installed on local machine
	if entry.generateCert {
		generateCertDir(entry.generatePath)
		generateRootCA(entry.generatePath)
		generateServerCA(entry.generatePath)
		clearTlsConfigFile(entry.generatePath)

		entry.keyFilePath = path.Join(entry.generatePath, "server-key.pem")
		entry.certFilePath = path.Join(entry.generatePath, "server.pem")
	}

	return entry
}

func (entry *TlsEntry) GetCertFilePath() string {
	return entry.certFilePath
}

func (entry *TlsEntry) GetKeyFilePath() string {
	return entry.keyFilePath
}

func (entry *TlsEntry) GetPort() uint64 {
	return entry.port
}

// clear ca-config.json, ca-csr.json and server.json
func clearTlsConfigFile(fullPath string) {
	os.RemoveAll(path.Join(fullPath, "ca-config.json"))
	os.RemoveAll(path.Join(fullPath, "ca-csr.json"))
	os.RemoveAll(path.Join(fullPath, "server.json"))
}

// generate root CA
func generateRootCA(fullPath string) {
	// create default one if ca-config.json file does not exist
	if !exists(path.Join(fullPath, "ca-config.json")) {
		if err := ioutil.WriteFile(path.Join(fullPath, "ca-config.json"), []byte(rootCAConfig), 0755); err != nil {
			shutdownWithError(err)
		}
	}

	// create default one if ca-csr.json file does not exist
	if !exists(path.Join(fullPath, "ca-csr.json")) {
		if err := ioutil.WriteFile(path.Join(fullPath, "ca-csr.json"), []byte(rootCACsr), 0755); err != nil {
			shutdownWithError(err)
		}
	}

	// create root cert request, cert and private key if files missing
	if !exists(path.Join(fullPath, "ca.csr")) ||
		!exists(path.Join(fullPath, "ca.pem")) ||
		!exists(path.Join(fullPath, "ca-key.pem")) {
		// generate root CA with bellow command
		// cfssl gencert -initca ca-csr.json | cfssljson -bare ca -
		// first, lets cd to cert directory
		os.Chdir(fullPath)
		defer os.Chdir(originalWD)
		// second, run the command
		cmd := "cfssl gencert -initca ca-csr.json | cfssljson -bare ca -"
		if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
			shutdownWithError(err)
		}
	}
}

// generate server CA
func generateServerCA(fullPath string) {
	// create default one if server.json file does not exist
	if !exists(path.Join(fullPath, "server.json")) {
		if err := ioutil.WriteFile(path.Join(fullPath, "server.json"), []byte(serverCAConfig), 0755); err != nil {
			shutdownWithError(err)
		}
	}

	// create server cert request, cert and private key if files missing
	if !exists(path.Join(fullPath, "server.csr")) ||
		!exists(path.Join(fullPath, "server.pem")) ||
		!exists(path.Join(fullPath, "server-key.pem")) {
		// generate root CA with bellow command
		// cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server server.json | cfssljson -bare server
		// first, lets cd to cert directory
		os.Chdir(fullPath)
		defer os.Chdir(originalWD)
		// second, run the command
		cmd := "cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server server.json | cfssljson -bare server"
		if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
			shutdownWithError(err)
		}
	}
}

// generate certs directory under working directory
func generateCertDir(fullPath string) {
	// create directory of cert
	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		shutdownWithError(err)
	}
}

func exists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func shutdownWithError(err error) {
	runtime2.PrintStack()
	glog.Error(err)
	os.Exit(1)
}
