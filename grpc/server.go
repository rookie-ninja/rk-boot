// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_grpc

import (
	"crypto/tls"
	"github.com/golang/glog"
	"github.com/rookie-ninja/rk-boot/gw"
	"github.com/rookie-ninja/rk-boot/sw"
	rk_tls "github.com/rookie-ninja/rk-boot/tls"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"strconv"
	"syscall"
)

type GRpcServerEntry struct {
	logger              *zap.Logger
	name                string
	port                uint64
	enableCommonService bool
	serverOpts          []grpc.ServerOption
	unaryInterceptors   []grpc.UnaryServerInterceptor
	streamInterceptors  []grpc.StreamServerInterceptor
	regFuncs            []RegFunc
	gw                  *rk_gw.GRpcGWEntry
	sw                  *rk_sw.SWEntry
	tls                 *rk_tls.TlsEntry
	server              *grpc.Server
	listener            net.Listener
}

type RegFunc func(server *grpc.Server)

type GRpcEntryOption func(*GRpcServerEntry)

func WithRegFuncs(funcs ...RegFunc) GRpcEntryOption {
	return func(entry *GRpcServerEntry) {
		entry.regFuncs = append(entry.regFuncs, funcs...)
	}
}

func WithGWEntry(gw *rk_gw.GRpcGWEntry) GRpcEntryOption {
	return func(entry *GRpcServerEntry) {
		entry.gw = gw
	}
}

func WithSWEntry(sw *rk_sw.SWEntry) GRpcEntryOption {
	return func(entry *GRpcServerEntry) {
		entry.sw = sw
	}
}

func WithTlsEntry(tls *rk_tls.TlsEntry) GRpcEntryOption {
	return func(entry *GRpcServerEntry) {
		entry.tls = tls
	}
}

func WithPort(port uint64) GRpcEntryOption {
	return func(entry *GRpcServerEntry) {
		entry.port = port
	}
}

func WithCommonService(enable bool) GRpcEntryOption {
	return func(entry *GRpcServerEntry) {
		entry.enableCommonService = enable
	}
}

func WithName(name string) GRpcEntryOption {
	return func(entry *GRpcServerEntry) {
		entry.name = name
	}
}

func WithServerOptions(opts ...grpc.ServerOption) GRpcEntryOption {
	return func(entry *GRpcServerEntry) {
		entry.serverOpts = append(entry.serverOpts, opts...)
	}
}

func WithUnaryInterceptors(opts ...grpc.UnaryServerInterceptor) GRpcEntryOption {
	return func(entry *GRpcServerEntry) {
		entry.unaryInterceptors = append(entry.unaryInterceptors, opts...)
	}
}

func WithStreamInterceptors(opts ...grpc.StreamServerInterceptor) GRpcEntryOption {
	return func(entry *GRpcServerEntry) {
		entry.streamInterceptors = append(entry.streamInterceptors, opts...)
	}
}

func NewGRpcServerEntry(opts ...GRpcEntryOption) *GRpcServerEntry {
	entry := &GRpcServerEntry{
		logger: zap.NewNop(),
	}

	for i := range opts {
		opts[i](entry)
	}

	if len(entry.name) < 1 {
		entry.name = "gRpc-server-" + strconv.FormatUint(entry.port, 10)
	}

	if entry.serverOpts == nil {
		entry.serverOpts = make([]grpc.ServerOption, 0)
	}

	if entry.unaryInterceptors == nil {
		entry.unaryInterceptors = make([]grpc.UnaryServerInterceptor, 0)
	}

	if entry.streamInterceptors == nil {
		entry.streamInterceptors = make([]grpc.StreamServerInterceptor, 0)
	}

	if entry.regFuncs == nil {
		entry.regFuncs = make([]RegFunc, 0)
	}

	return entry
}

func (entry *GRpcServerEntry) AddServerOptions(opts ...grpc.ServerOption) {
	entry.serverOpts = append(entry.serverOpts, opts...)
}

func (entry *GRpcServerEntry) AddUnaryInterceptors(inter ...grpc.UnaryServerInterceptor) {
	entry.unaryInterceptors = append(entry.unaryInterceptors, inter...)
}

func (entry *GRpcServerEntry) AddStreamInterceptors(inter ...grpc.StreamServerInterceptor) {
	entry.streamInterceptors = append(entry.streamInterceptors, inter...)
}

func (entry *GRpcServerEntry) AddRegFuncs(funcs ...RegFunc) {
	entry.regFuncs = append(entry.regFuncs, funcs...)
}

func (entry *GRpcServerEntry) AddGWRegFuncs(funcs ...rk_gw.RegFunc) {
	if entry.gw != nil {
		entry.gw.AddRegFuncs(funcs...)
	}

	if entry.sw != nil {
		entry.sw.AddRegFuncs(funcs...)
	}
}

func (entry *GRpcServerEntry) GetPort() uint64 {
	return entry.port
}

func (entry *GRpcServerEntry) GetName() string {
	return entry.name
}

func (entry *GRpcServerEntry) GetServer() *grpc.Server {
	return entry.server
}

func (entry *GRpcServerEntry) GetListener() net.Listener {
	return entry.listener
}

func (entry *GRpcServerEntry) GetGWEntry() *rk_gw.GRpcGWEntry {
	return entry.gw
}

func (entry *GRpcServerEntry) GetSWEntry() *rk_sw.SWEntry {
	return entry.sw
}

func (entry *GRpcServerEntry) Stop(logger *zap.Logger) {
	if entry.server != nil {
		if logger == nil {
			logger = zap.NewNop()
		}

		fields := []zap.Field{
			zap.Uint64("gRpc_port", entry.GetPort()),
			zap.String("name", entry.name),
		}

		if entry.tls != nil {
			fields = append(fields, zap.Bool("tls", true))
		}

		logger.Info("stopping gRpc-server", fields...)
		entry.server.GracefulStop()
	}
}

func (entry *GRpcServerEntry) StopGW(logger *zap.Logger) {
	entry.gw.Stop(logger)
}

func (entry *GRpcServerEntry) StopSW(logger *zap.Logger) {
	entry.sw.Stop(logger)
}

func (entry *GRpcServerEntry) Start(logger *zap.Logger) {
	if logger == nil {
		logger = zap.NewNop()
	}

	listener, err := net.Listen("tcp4", ":"+strconv.FormatUint(entry.port, 10))
	if err != nil {
		shutdownWithError(err)
	}

	entry.listener = listener
	// make unary server opts
	entry.serverOpts = append(entry.serverOpts, grpc.ChainUnaryInterceptor(entry.unaryInterceptors...))
	// make stream server opts
	entry.serverOpts = append(entry.serverOpts, grpc.ChainStreamInterceptor(entry.streamInterceptors...))

	if entry.tls != nil {
		cert, err := tls.LoadX509KeyPair(entry.tls.GetCertFilePath(), entry.tls.GetKeyFilePath())
		if err != nil {
			shutdownWithError(err)
		}
		entry.serverOpts = append(entry.serverOpts, grpc.Creds(credentials.NewServerTLSFromCert(&cert)))
	}

	entry.server = grpc.NewServer(entry.serverOpts...)
	for _, regFunc := range entry.regFuncs {
		regFunc(entry.server)
	}

	go func(entry *GRpcServerEntry) {
		fields := []zap.Field{
			zap.Uint64("gRpc_port", entry.GetPort()),
			zap.String("name", entry.name),
		}

		if entry.tls != nil {
			fields = append(fields, zap.Bool("tls", true))
		}

		logger.Info("starting gRpc-server", fields...)
		if err := entry.server.Serve(listener); err != nil {
			fields = append(fields, zap.Error(err))
			logger.Error("err while serving gRpc-listener", fields...)
			shutdownWithError(err)
		}
	}(entry)
}

func (entry *GRpcServerEntry) StartGW(logger *zap.Logger) {
	entry.gw.Start(logger)
}

func (entry *GRpcServerEntry) StartSW(logger *zap.Logger) {
	entry.sw.Start(logger)
}

func shutdownWithError(err error) {
	glog.Error(err)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}
