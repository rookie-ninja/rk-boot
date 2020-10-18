// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_boot

import (
	"context"
	"github.com/rookie-ninja/rk-boot/api/v1"
	"github.com/rookie-ninja/rk-interceptor/context"
	"github.com/spf13/cast"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime"
	"syscall"
)

type CommonServiceGRpc struct{}

func NewCommonServiceGRpc() *CommonServiceGRpc {
	service := &CommonServiceGRpc{}
	return service
}

// GC Stub
func (service *CommonServiceGRpc) GC(ctx context.Context, request *rk_boot_common_v1.GCRequest) (*rk_boot_common_v1.GCResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)
	event := rk_inter_context.GetEvent(ctx)

	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)

	runtime.GC()
	runtime.ReadMemStats(&after)

	event.AddFields(memStatsToZapFields("before_", &before)...)
	event.AddFields(memStatsToZapFields("after_", &before)...)

	res := &rk_boot_common_v1.GCResponse{
		MemStatsBefore: memStatsToPB(&before),
		MemStatsAfter:  memStatsToPB(&after),
	}

	return res, nil
}

// DumpConfig Stub
func (service *CommonServiceGRpc) DumpConfig(ctx context.Context, request *rk_boot_common_v1.DumpConfigRequest) (*rk_boot_common_v1.DumpConfigResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)

	configList := make([]*rk_boot_common_v1.Config, 0)
	res := &rk_boot_common_v1.DumpConfigResponse{ConfigList: configList}

	// rk-configs
	for k, v := range AppCtx.ListRkConfigs() {
		configPairs := make([]*rk_boot_common_v1.ConfigPair, 0)
		for i := range v.GetViper().AllKeys() {
			viperKey := v.GetViper().AllKeys()[i]
			viperValue := cast.ToString(v.GetViper().Get(viperKey))

			pair := &rk_boot_common_v1.ConfigPair{
				Key:   viperKey,
				Value: viperValue,
			}

			configPairs = append(configPairs, pair)
		}

		conf := &rk_boot_common_v1.Config{
			ConfigName: k,
			ConfigPair: configPairs,
		}

		res.ConfigList = append(res.ConfigList, conf)
	}

	// viper-configs
	for k, v := range AppCtx.ListViperConfigs() {
		configPairs := make([]*rk_boot_common_v1.ConfigPair, 0)
		for i := range v.AllKeys() {
			viperKey := v.AllKeys()[i]
			viperValue := cast.ToString(v.Get(viperKey))

			pair := &rk_boot_common_v1.ConfigPair{
				Key:   viperKey,
				Value: viperValue,
			}

			configPairs = append(configPairs, pair)
		}

		conf := &rk_boot_common_v1.Config{
			ConfigName: k,
			ConfigPair: configPairs,
		}

		res.ConfigList = append(res.ConfigList, conf)
	}

	return res, nil
}

// GetConfig Stub
func (service *CommonServiceGRpc) GetConfig(ctx context.Context, request *rk_boot_common_v1.GetConfigRequest) (*rk_boot_common_v1.GetConfigResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)

	configList := make([]*rk_boot_common_v1.Config, 0)
	res := &rk_boot_common_v1.GetConfigResponse{ConfigList: configList}

	for k, v := range AppCtx.ListRkConfigs() {
		pair := &rk_boot_common_v1.ConfigPair{
			Key:   request.Key,
			Value: cast.ToString(v.Get(request.GetKey())),
		}

		conf := &rk_boot_common_v1.Config{
			ConfigName: k,
			ConfigPair: []*rk_boot_common_v1.ConfigPair{pair},
		}

		res.ConfigList = append(res.ConfigList, conf)
	}

	for k, v := range AppCtx.ListViperConfigs() {
		pair := &rk_boot_common_v1.ConfigPair{
			Key:   request.Key,
			Value: cast.ToString(v.Get(request.GetKey())),
		}

		conf := &rk_boot_common_v1.Config{
			ConfigName: k,
			ConfigPair: []*rk_boot_common_v1.ConfigPair{pair},
		}

		res.ConfigList = append(res.ConfigList, conf)
	}

	return res, nil
}

// Ping Stub
func (service *CommonServiceGRpc) Ping(ctx context.Context, request *rk_boot_common_v1.PingRequest) (*rk_boot_common_v1.PongResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)

	res := &rk_boot_common_v1.PongResponse{
		Message: "pong",
	}

	return res, nil
}

// Log Stub
func (service *CommonServiceGRpc) Log(ctx context.Context, request *rk_boot_common_v1.LogRequest) (*rk_boot_common_v1.LogResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)
	event := rk_inter_context.GetEvent(ctx)

	for i := range request.Entries {
		entry := request.Entries[i]
		loggerConfig := AppCtx.GetLoggerConfig(entry.LogName)
		if loggerConfig != nil {
			setLogLevel(loggerConfig, entry.LogLevel)
		}

		event.AddPair(entry.LogName, entry.LogLevel)
	}

	res := &rk_boot_common_v1.LogResponse{}

	return res, nil
}

// Shutdown Stub
func (service *CommonServiceGRpc) Shutdown(ctx context.Context, request *rk_boot_common_v1.ShutdownRequest) (*rk_boot_common_v1.ShutdownResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)
	event := rk_inter_context.GetEvent(ctx)
	event.AddPair("signal", "interrupt")

	res := &rk_boot_common_v1.ShutdownResponse{
		Message: "interrupt",
	}

	defer func() {
		AppCtx.GetShutdownSig() <- syscall.SIGINT
	}()
	return res, nil
}

// Info Stub
func (service *CommonServiceGRpc) Info(ctx context.Context, request *rk_boot_common_v1.InfoRequest) (*rk_boot_common_v1.InfoResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)
	event := rk_inter_context.GetEvent(ctx)

	res := &rk_boot_common_v1.InfoResponse{}

	boot := AppCtx.ListGRpcEntries()

	if boot == nil {
		event.InCCounter("rk_common_service_boot_nil", 1)
		return nil, status.Error(codes.Unimplemented,
			"failed to find pl bootstrapper, system may not started with pl_boot.")
	}

	// Basic info
	fillBasicInfo(res)

	// gRPC info
	fillGRPCInfo(res)

	// Prom info
	fillPromInfo(res)

	return res, nil
}

// Healthy Stub
func (service *CommonServiceGRpc) Healthy(ctx context.Context, request *rk_boot_common_v1.HealthyRequest) (*rk_boot_common_v1.HealthyResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)
	event := rk_inter_context.GetEvent(ctx)

	event.AddPair("healthy", "true")

	res := &rk_boot_common_v1.HealthyResponse{
		Healthy: true,
	}

	return res, nil
}
