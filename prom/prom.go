// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_prom

import (
	"context"
	"github.com/rookie-ninja/rk-prom"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

type PromEntry struct {
	logger        *zap.Logger
	port          uint64
	path          string
	pgwRemoteAddr string
	pgwIntervalMS uint64
	pgwPublisher  rk_prom.Publisher
	pgwJobName    string
	server        *http.Server
}

type PromOption func(*PromEntry)

func WithPort(port uint64) PromOption {
	return func(entry *PromEntry) {
		entry.port = port
	}
}

func WithPath(path string) PromOption {
	return func(entry *PromEntry) {
		entry.path = path
	}
}

func WithLogger(logger *zap.Logger) PromOption {
	return func(entry *PromEntry) {
		entry.logger = logger
	}
}

func WithPGWRemoteAddr(addr string) PromOption {
	return func(entry *PromEntry) {
		entry.pgwRemoteAddr = addr
	}
}

func WithPGWJobName(name string) PromOption {
	return func(entry *PromEntry) {
		entry.pgwJobName = name
	}
}

func WithPGWIntervalMS(interval uint64) PromOption {
	return func(entry *PromEntry) {
		entry.pgwIntervalMS = interval
	}
}

func NewPromEntry(opts ...PromOption) *PromEntry {
	entry := &PromEntry{
		logger: zap.NewNop(),
	}

	for i := range opts {
		opts[i](entry)
	}

	// pushGateway publisher enabled
	if len(entry.pgwRemoteAddr) > 0 {
		// set default value
		if entry.pgwIntervalMS == 0 {
			entry.pgwIntervalMS = uint64(1 * time.Millisecond)
		}

		// set default job name
		if len(entry.pgwJobName) < 1 {
			entry.pgwJobName = "rk-job"
		}
	}

	return entry
}

func (entry *PromEntry) GetPort() uint64 {
	return entry.port
}

func (entry *PromEntry) GetPath() string {
	return entry.path
}

func (entry *PromEntry) GetPGWRemoteAddr() string {
	return entry.pgwRemoteAddr
}

func (entry *PromEntry) GetPGWIntervalMS() uint64 {
	return entry.pgwIntervalMS
}

func (entry *PromEntry) GetServer() *http.Server {
	return entry.server
}

func (entry *PromEntry) Start(logger *zap.Logger) {
	rk_prom.SetZapLogger(logger)
	server, err := rk_prom.StartProm(entry.port, entry.path)
	if err != nil {
		shutdownWithError(err)
	}

	// PushGateway publisher enabled
	if len(entry.pgwRemoteAddr) > 0 {
		entry.pgwPublisher, err = rk_prom.NewPushGatewayPublisher(
			time.Duration(entry.pgwIntervalMS)*time.Second,
			entry.pgwRemoteAddr,
			entry.pgwJobName)
	}

	entry.server = server
}

func (entry *PromEntry) Stop(logger *zap.Logger) {
	logger.Info("stopping prometheus client",
		zap.Uint64("port", entry.port),
		zap.String("path", entry.path))
	if entry.server != nil {
		entry.server.Shutdown(context.Background())
	}

	if entry.pgwPublisher != nil {
		logger.Info("stopping pushGateway publisher",
			zap.String("remote_addr", entry.pgwRemoteAddr),
			zap.String("job_name", entry.pgwJobName),
			zap.Uint64("interval_ms", entry.pgwIntervalMS))
		entry.pgwPublisher.Shutdown()
	}
}

func shutdownWithError(err error) {
	// log it
	os.Exit(1)
}
