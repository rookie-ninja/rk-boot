// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_gin

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/rookie-ninja/rk-boot/sw"
	"github.com/rookie-ninja/rk-boot/tls"
	"github.com/rookie-ninja/rk-gin-interceptor/panic/zap"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"path"
	runtime2 "runtime/debug"
	"strconv"
)

type GinServerEntry struct {
	logger       *zap.Logger
	router       *gin.Engine
	server       *http.Server
	tlsServer    *http.Server
	name         string
	port         uint64
	interceptors []gin.HandlerFunc
	sw           *rk_sw.SWEntry
	tls          *rk_tls.TlsEntry
}

type GinEntryOption func(*GinServerEntry)

func WithInterceptors(inters ...gin.HandlerFunc) GinEntryOption {
	return func(entry *GinServerEntry) {
		if entry.interceptors == nil {
			entry.interceptors = make([]gin.HandlerFunc, 0)
		}

		entry.interceptors = append(entry.interceptors, inters...)
	}
}

func WithRouter(router *gin.Engine) GinEntryOption {
	return func(entry *GinServerEntry) {
		entry.router = router
	}
}

func WithTlsEntry(tls *rk_tls.TlsEntry) GinEntryOption {
	return func(entry *GinServerEntry) {
		entry.tls = tls
	}
}

func WithSWEntry(sw *rk_sw.SWEntry) GinEntryOption {
	return func(entry *GinServerEntry) {
		entry.sw = sw
	}
}

func WithPort(port uint64) GinEntryOption {
	return func(entry *GinServerEntry) {
		entry.port = port
	}
}

func WithTlsPort(port uint64) GinEntryOption {
	return func(entry *GinServerEntry) {
		entry.port = port
	}
}

func WithName(name string) GinEntryOption {
	return func(entry *GinServerEntry) {
		entry.name = name
	}
}

func NewGinServerEntry(opts ...GinEntryOption) *GinServerEntry {
	entry := &GinServerEntry{
		logger: zap.NewNop(),
	}

	for i := range opts {
		opts[i](entry)
	}

	if len(entry.name) < 1 {
		entry.name = "gin-server-" + strconv.FormatUint(entry.port, 10)
	}

	if entry.interceptors == nil {
		entry.interceptors = make([]gin.HandlerFunc, 0)
	}

	if entry.router == nil {
		gin.SetMode(gin.ReleaseMode)
		entry.router = gin.New()
	}

	if entry.sw != nil {
		entry.router.GET(path.Join(entry.sw.GetPath(), "*any"), entry.sw.GinHandler())
		entry.router.GET("/swagger/*any", entry.sw.GinFileHandler())
	}

	entry.interceptors = append(entry.interceptors, rk_gin_inter_panic.RkGinPanicZap())

	if len(entry.interceptors) > 0 {
		entry.router.Use(entry.interceptors...)
	}

	// init tls server only if port is not zero
	if entry.tls != nil && entry.tls.GetPort() != 0 {
		entry.tlsServer = &http.Server{
			Addr:    "0.0.0.0:" + strconv.FormatUint(entry.tls.GetPort(), 10),
			Handler: entry.router,
		}
	}

	// init server only if port is not zero
	if entry.port != 0 {
		if entry.tls != nil && entry.tls.GetPort() != entry.port {
			entry.server = &http.Server{
				Addr:    "0.0.0.0:" + strconv.FormatUint(entry.port, 10),
				Handler: entry.router,
			}
		}
	}

	// make sure we keep only one server


	return entry
}

func (entry *GinServerEntry) GetName() string {
	return entry.name
}

func (entry *GinServerEntry) GetPort() uint64 {
	return entry.port
}

func (entry *GinServerEntry) GetSWEntry() *rk_sw.SWEntry {
	return entry.sw
}

func (entry *GinServerEntry) GetTlsEntry() *rk_tls.TlsEntry {
	return entry.tls
}

func (entry *GinServerEntry) GetServer() *http.Server {
	return entry.server
}

func (entry *GinServerEntry) GetTlsServer() *http.Server {
	return entry.tlsServer
}

func (entry *GinServerEntry) GetRouter() *gin.Engine {
	return entry.router
}

func (entry *GinServerEntry) Start(logger *zap.Logger) {
	if logger == nil {
		logger = zap.NewNop()
	}

	go func(entry *GinServerEntry) {
		// Start server with tls
		if entry.tlsServer != nil {
			fields := []zap.Field{
				zap.Uint64("gin_tls_port", entry.tls.GetPort()),
				zap.String("name", entry.name),
			}

			if entry.sw != nil {
				fields = append(fields, zap.String("sw_path", entry.sw.GetPath()))
			}

			logger.Info("starting gin-tls-server", fields...)

			if err := entry.tlsServer.ListenAndServeTLS(entry.tls.GetCertFilePath(), entry.tls.GetKeyFilePath()); err != nil && err != http.ErrServerClosed {
				fields = append(fields, zap.Error(err))
				logger.Error("err while serving gin-tls-listener", fields...)
				shutdownWithError(err)
			}
		}
	}(entry)

	go func(entry *GinServerEntry) {
		// Start server
		if entry.server != nil {
			fields := []zap.Field{
				zap.Uint64("gin_port", entry.GetPort()),
				zap.String("name", entry.name),
			}

			if entry.sw != nil {
				fields = append(fields, zap.String("sw_path", entry.sw.GetPath()))
			}

			logger.Info("starting gin-server", fields...)

			if err := entry.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fields = append(fields, zap.Error(err))
				logger.Error("err while serving gin-listener", fields...)
				shutdownWithError(err)
			}
		}
	}(entry)
}

func (entry *GinServerEntry) Stop(logger *zap.Logger) {
	if entry.router != nil {
		if logger == nil {
			logger = zap.NewNop()
		}

		if entry.server != nil {
			logger.Info("stopping gin-server",
				zap.Uint64("gin_port", entry.port),
				zap.String("name", entry.name))
			if err := entry.server.Shutdown(context.Background()); err != nil {
				logger.Warn("error occurs while stopping gin server",
					zap.Uint64("gin_port", entry.port),
					zap.String("name", entry.name),
					zap.Error(err))
			}
		}

		if entry.tlsServer != nil {
			logger.Info("stopping gin-tls-server",
				zap.Uint64("gin_tls_port", entry.tls.GetPort()),
				zap.String("name", entry.name))
			if err := entry.tlsServer.Shutdown(context.Background()); err != nil {
				logger.Warn("error occurs while stopping gin tls server",
					zap.Uint64("gin_tls_port", entry.tls.GetPort()),
					zap.String("name", entry.name),
					zap.Error(err))
			}
		}
	}
}

func shutdownWithError(err error) {
	runtime2.PrintStack()
	glog.Error(err)
	os.Exit(1)
}
