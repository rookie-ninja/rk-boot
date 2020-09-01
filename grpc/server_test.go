// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_grpc

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rookie-ninja/rk-boot/gw"
	"github.com/rookie-ninja/rk-boot/sw"
	"github.com/rookie-ninja/rk-interceptor/logging/zap"
	"github.com/rookie-ninja/rk-query"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"strings"
	"testing"
)

// Test WithRegFuncs with empty input
func TestWithRegFuncs_WithEmptyInput(t *testing.T) {
	entry := &GRpcServerEntry{
		regFuncs: make([]RegFunc, 0, 0),
	}

	option := WithRegFuncs()
	option(entry)

	assert.Empty(t, entry.regFuncs)
}

// Test WithRegFuncs happy case
func TestWithRegFuncs_HappyCase(t *testing.T) {
	entry := &GRpcServerEntry{
		regFuncs: make([]RegFunc, 0, 0),
	}

	one := func(*grpc.Server) {}
	two := func(*grpc.Server) {}
	option := WithRegFuncs(one, two)
	option(entry)

	assert.Len(t, entry.regFuncs, 2)
}

// Test WithGWEntry with nil
func TestWithGWEntry_WithNil(t *testing.T) {
	entry := &GRpcServerEntry{}

	option := WithGWEntry(nil)
	option(entry)

	assert.Nil(t, entry.gw)
}

// Test WithGWEntry happy case
func TestWithGWEntry_HappyCase(t *testing.T) {
	entry := &GRpcServerEntry{}

	option := WithGWEntry(&rk_gw.GRpcGWEntry{})
	option(entry)

	assert.NotNil(t, entry.gw)
}

// Test WithSWEntry with nil
func TestWithSWEntry_WithNil(t *testing.T) {
	entry := &GRpcServerEntry{}

	option := WithSWEntry(nil)
	option(entry)

	assert.Nil(t, entry.gw)
}

// Test WithSWEntry happy case
func TestWithSWEntry_HappyCase(t *testing.T) {
	entry := &GRpcServerEntry{}

	option := WithSWEntry(&rk_sw.SWEntry{})
	option(entry)

	assert.NotNil(t, entry.sw)
}

// Test WithPort happy case
func TestWithPort_HappyCase(t *testing.T) {
	entry := &GRpcServerEntry{}

	option := WithPort(1)
	option(entry)

	assert.Equal(t, uint64(1), entry.port)
}

// Test WithCommonService with false
func TestWithCommonService_WithFalse(t *testing.T) {
	entry := &GRpcServerEntry{}

	option := WithCommonService(false)
	option(entry)

	assert.False(t, entry.enableCommonService)
}

// Test WithCommonService with true
func TestWithCommonService_WithTrue(t *testing.T) {
	entry := &GRpcServerEntry{}

	option := WithCommonService(true)
	option(entry)

	assert.True(t, entry.enableCommonService)
}

// Test WithName with empty string
func TestWithName_WithEmptyString(t *testing.T) {
	entry := &GRpcServerEntry{}

	option := WithName("")
	option(entry)

	assert.Empty(t, entry.name)
}

// Test WithName happy case
func TestWithName_HappyCase(t *testing.T) {
	entry := &GRpcServerEntry{}

	option := WithName("fake-name")
	option(entry)

	assert.Equal(t, "fake-name", entry.name)
}

// Test WithServerOptions with empty input
func TestWithServerOptions_WithEmptyInput(t *testing.T) {
	entry := &GRpcServerEntry{
		regFuncs: make([]RegFunc, 0, 0),
	}

	option := WithServerOptions()
	option(entry)

	assert.Empty(t, entry.serverOpts)
}

// Test WithServerOptions happy case
func TestWithServerOptions_HappyCase(t *testing.T) {
	entry := &GRpcServerEntry{
		regFuncs: make([]RegFunc, 0, 0),
	}

	option := WithServerOptions(grpc.EmptyServerOption{}, grpc.EmptyServerOption{})
	option(entry)

	assert.Len(t, entry.serverOpts, 2)
}

// Test WithUnaryInterceptors with empty input
func TestWithUnaryInterceptors_WithEmptyInput(t *testing.T) {
	entry := &GRpcServerEntry{
		regFuncs: make([]RegFunc, 0, 0),
	}

	option := WithUnaryInterceptors()
	option(entry)

	assert.Empty(t, entry.unaryInterceptors)
}

// Test WithUnaryInterceptors happy case
func TestWithUnaryInterceptors_HappyCase(t *testing.T) {
	entry := &GRpcServerEntry{
		regFuncs: make([]RegFunc, 0, 0),
	}

	one := rk_inter_logging.UnaryServerInterceptor(rk_query.NewEventFactory())
	two := rk_inter_logging.UnaryServerInterceptor(rk_query.NewEventFactory())
	option := WithUnaryInterceptors(one, two)
	option(entry)

	assert.Len(t, entry.unaryInterceptors, 2)
}

// Test WithStreamInterceptors with empty input
func TestWithStreamInterceptors_WithEmptyInput(t *testing.T) {
	entry := &GRpcServerEntry{
		regFuncs: make([]RegFunc, 0, 0),
	}

	option := WithStreamInterceptors()
	option(entry)

	assert.Empty(t, entry.streamInterceptors)
}

// Test WithStreamInterceptors happy case
func TestWithStreamInterceptors_HappyCase(t *testing.T) {
	entry := &GRpcServerEntry{
		regFuncs: make([]RegFunc, 0, 0),
	}

	one := rk_inter_logging.StreamServerInterceptor(rk_query.NewEventFactory())
	two := rk_inter_logging.StreamServerInterceptor(rk_query.NewEventFactory())
	option := WithStreamInterceptors(one, two)
	option(entry)

	assert.Len(t, entry.streamInterceptors, 2)
}

// Test NewGRpcServerEntry with empty options
func TestNewGRpcServerEntry_WithEmptyInput(t *testing.T) {
	entry := NewGRpcServerEntry()

	assert.NotNil(t, entry)
	assert.NotNil(t, entry.logger)
	assert.True(t, strings.HasPrefix(entry.name, "gRpc-server-"))
	assert.Empty(t, entry.serverOpts)
	assert.Empty(t, entry.unaryInterceptors)
	assert.Empty(t, entry.streamInterceptors)
	assert.Empty(t, entry.regFuncs)
}

// Test NewGRpcServerEntry with name
func TestNewGRpcServerEntry_WithName(t *testing.T) {
	entry := NewGRpcServerEntry(WithName("fake-name"))
	assert.Equal(t, "fake-name", entry.name)

	assert.NotNil(t, entry)
	assert.NotNil(t, entry.logger)
	assert.Empty(t, entry.serverOpts)
	assert.Empty(t, entry.unaryInterceptors)
	assert.Empty(t, entry.streamInterceptors)
	assert.Empty(t, entry.regFuncs)
}

// Test NewGRpcServerEntry with server options
func TestNewGRpcServerEntry_WithServerOpts(t *testing.T) {
	entry := NewGRpcServerEntry(WithServerOptions(grpc.EmptyServerOption{}))

	assert.NotNil(t, entry)
	assert.NotNil(t, entry.logger)
	assert.NotEmpty(t, entry.serverOpts)
	assert.True(t, strings.HasPrefix(entry.name, "gRpc-server-"))
	assert.Empty(t, entry.unaryInterceptors)
	assert.Empty(t, entry.streamInterceptors)
	assert.Empty(t, entry.regFuncs)
}

// Test NewGRpcServerEntry with UnaryServerInterceptor
func TestNewGRpcServerEntry_WithUnaryServerInterceptor(t *testing.T) {
	one := rk_inter_logging.UnaryServerInterceptor(rk_query.NewEventFactory())
	entry := NewGRpcServerEntry(WithUnaryInterceptors(one))

	assert.NotNil(t, entry)
	assert.NotNil(t, entry.logger)
	assert.Empty(t, entry.serverOpts)
	assert.True(t, strings.HasPrefix(entry.name, "gRpc-server-"))
	assert.NotEmpty(t, entry.unaryInterceptors)
	assert.Empty(t, entry.streamInterceptors)
	assert.Empty(t, entry.regFuncs)
}

// Test NewGRpcServerEntry with StreamServerInterceptor
func TestNewGRpcServerEntry_WithStreamServerInterceptor(t *testing.T) {
	one := rk_inter_logging.StreamServerInterceptor(rk_query.NewEventFactory())
	entry := NewGRpcServerEntry(WithStreamInterceptors(one))

	assert.NotNil(t, entry)
	assert.NotNil(t, entry.logger)
	assert.Empty(t, entry.serverOpts)
	assert.True(t, strings.HasPrefix(entry.name, "gRpc-server-"))
	assert.Empty(t, entry.unaryInterceptors)
	assert.NotEmpty(t, entry.streamInterceptors)
	assert.Empty(t, entry.regFuncs)
}

// Test NewGRpcServerEntry with RegFunc
func TestNewGRpcServerEntry_WithRegFunc(t *testing.T) {
	one := func(*grpc.Server) {}
	entry := NewGRpcServerEntry(WithRegFuncs(one))

	assert.NotNil(t, entry)
	assert.NotNil(t, entry.logger)
	assert.Empty(t, entry.serverOpts)
	assert.True(t, strings.HasPrefix(entry.name, "gRpc-server-"))
	assert.Empty(t, entry.unaryInterceptors)
	assert.Empty(t, entry.streamInterceptors)
	assert.NotEmpty(t, entry.regFuncs)
}

// Test AddServerOptions with empty input
func TestGRpcServerEntry_AddServerOptions_WithEmptyInput(t *testing.T) {
	entry := NewGRpcServerEntry()
	entry.AddServerOptions()

	assert.Empty(t, entry.serverOpts)
}

// Test AddServerOptions happy case
func TestGRpcServerEntry_AddServerOptions_HappyCases(t *testing.T) {
	entry := NewGRpcServerEntry()
	entry.AddServerOptions(grpc.EmptyServerOption{})

	assert.NotEmpty(t, entry.serverOpts)
}

// Test AddUnaryInterceptors with empty input
func TestGRpcServerEntry_AddUnaryInterceptors_WithEmptyInput(t *testing.T) {
	entry := NewGRpcServerEntry()
	entry.AddUnaryInterceptors()

	assert.Empty(t, entry.unaryInterceptors)
}

// Test AddUnaryInterceptors happy case
func TestGRpcServerEntry_AddUnaryInterceptors_HappyCases(t *testing.T) {
	entry := NewGRpcServerEntry()
	one := rk_inter_logging.UnaryServerInterceptor(rk_query.NewEventFactory())
	entry.AddUnaryInterceptors(one)

	assert.NotEmpty(t, entry.unaryInterceptors)
}

// Test AddStreamInterceptors with empty input
func TestGRpcServerEntry_AddStreamInterceptors_WithEmptyInput(t *testing.T) {
	entry := NewGRpcServerEntry()
	entry.AddStreamInterceptors()

	assert.Empty(t, entry.streamInterceptors)
}

// Test AddStreamInterceptors happy case
func TestGRpcServerEntry_AddStreamInterceptors_HappyCases(t *testing.T) {
	entry := NewGRpcServerEntry()
	one := rk_inter_logging.StreamServerInterceptor(rk_query.NewEventFactory())
	entry.AddStreamInterceptors(one)

	assert.NotEmpty(t, entry.streamInterceptors)
}

// Test AddRegFuncs with empty input
func TestGRpcServerEntry_AddRegFuncs_WithEmptyInput(t *testing.T) {
	entry := NewGRpcServerEntry()
	entry.AddRegFuncs()

	assert.Empty(t, entry.regFuncs)
}

// Test AddRegFuncs happy case
func TestGRpcServerEntry_AddRegFuncs_HappyCases(t *testing.T) {
	entry := NewGRpcServerEntry()
	one := func(*grpc.Server) {}
	entry.AddRegFuncs(one)

	assert.NotEmpty(t, entry.regFuncs)
}

// Test AddGWRegFuncs with empty input
func TestGRpcServerEntry_AddGWRegFuncs_WithEmptyInput(t *testing.T) {
	entry := NewGRpcServerEntry()
	// There should be no panic
	entry.AddGWRegFuncs()
}

// Test AddGWRegFuncs happy case
func TestGRpcServerEntry_AddGWRegFuncs_HappyCases(t *testing.T) {
	entry := NewGRpcServerEntry()
	one := func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error { return nil }
	// There should not be panic
	entry.AddGWRegFuncs(one)
}

// Test GetPort
func TestGRpcServerEntry_GetPort(t *testing.T) {
	entry := NewGRpcServerEntry()
	assert.Equal(t, entry.port, entry.GetPort())
}

// Test GetName
func TestGRpcServerEntry_GetName(t *testing.T) {
	entry := NewGRpcServerEntry()
	assert.Equal(t, entry.name, entry.GetName())
}

// Test GetServer
func TestGRpcServerEntry_GetServer(t *testing.T) {
	entry := NewGRpcServerEntry()
	assert.Equal(t, entry.server, entry.GetServer())
}

// Test GetListener
func TestGRpcServerEntry_GetListener(t *testing.T) {
	entry := NewGRpcServerEntry()
	assert.Equal(t, entry.listener, entry.GetListener())
}

// Test GetGWEntry
func TestGRpcServerEntry_GetGWEntry(t *testing.T) {
	entry := NewGRpcServerEntry()
	assert.Equal(t, entry.gw, entry.GetGWEntry())
}

// Test GetSWEntry
func TestGRpcServerEntry_GetSWEntry(t *testing.T) {
	entry := NewGRpcServerEntry()
	assert.Equal(t, entry.sw, entry.GetSWEntry())
}

// Test Stop function
func TestGRpcServerEntry_Stop(t *testing.T) {
	entry := NewGRpcServerEntry()
	// There should be no panic
	entry.Stop(nil)
}

// Test StopGW function
func TestGRpcServerEntry_StopGW(t *testing.T) {
	entry := NewGRpcServerEntry()
	// There should be no panic
	entry.StopGW(nil)
}

// Test StopSW function
func TestGRpcServerEntry_StopSW(t *testing.T) {
	entry := NewGRpcServerEntry()
	// There should be no panic
	entry.StopSW(nil)
}
