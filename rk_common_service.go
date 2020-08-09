// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_boot

import (
	"context"
	"github.com/rookie-ninja/rk-boot/api"
	"github.com/rookie-ninja/rk-interceptor/context"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

type CommonService struct{}

func NewCommonService() *CommonService {
	service := &CommonService{}
	return service
}

// GC Stub
func (service *CommonService) GC(ctx context.Context, request *api.GCRequest) (*api.GCResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)
	event := rk_inter_context.GetEvent(ctx)

	event.AddPair("operator", request.Operator)
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)

	runtime.GC()
	runtime.ReadMemStats(&after)

	event.AddFields(memStatsToZapFields("before_", &before)...)
	event.AddFields(memStatsToZapFields("after_", &before)...)

	res := &api.GCResponse{
		MemStatsBefore: memStatsToPB(&before),
		MemStatsAfter:  memStatsToPB(&after),
	}

	return res, nil
}

// DumpConfig Stub
func (service *CommonService) DumpConfig(ctx context.Context, request *api.DumpConfigRequest) (*api.DumpConfigResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)

	configList := make([]*api.Config, 0)
	res := &api.DumpConfigResponse{ConfigList: configList}

	// rk-configs
	for k, v := range AppCtx.ListRkConfigs() {
		configPairs := make([]*api.ConfigPair, 0)
		for i := range v.GetViper().AllKeys() {
			viperKey := v.GetViper().AllKeys()[i]
			viperValue := cast.ToString(v.GetViper().Get(viperKey))

			pair := &api.ConfigPair{
				Key:   viperKey,
				Value: viperValue,
			}

			configPairs = append(configPairs, pair)
		}

		conf := &api.Config{
			ConfigName: k,
			ConfigPair: configPairs,
		}

		res.ConfigList = append(res.ConfigList, conf)
	}

	// viper-configs
	for k, v := range AppCtx.ListViperConfigs() {
		configPairs := make([]*api.ConfigPair, 0)
		for i := range v.AllKeys() {
			viperKey := v.AllKeys()[i]
			viperValue := cast.ToString(v.Get(viperKey))

			pair := &api.ConfigPair{
				Key:   viperKey,
				Value: viperValue,
			}

			configPairs = append(configPairs, pair)
		}

		conf := &api.Config{
			ConfigName: k,
			ConfigPair: configPairs,
		}

		res.ConfigList = append(res.ConfigList, conf)
	}

	return res, nil
}

// GetConfig Stub
func (service *CommonService) GetConfig(ctx context.Context, request *api.GetConfigRequest) (*api.GetConfigResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)

	configList := make([]*api.Config, 0)
	res := &api.GetConfigResponse{ConfigList: configList}

	for k, v := range AppCtx.ListRkConfigs() {
		pair := &api.ConfigPair{
			Key:   request.Key,
			Value: cast.ToString(v.Get(request.GetKey())),
		}

		conf := &api.Config{
			ConfigName: k,
			ConfigPair: []*api.ConfigPair{pair},
		}

		res.ConfigList = append(res.ConfigList, conf)
	}

	for k, v := range AppCtx.ListViperConfigs() {
		pair := &api.ConfigPair{
			Key:   request.Key,
			Value: cast.ToString(v.Get(request.GetKey())),
		}

		conf := &api.Config{
			ConfigName: k,
			ConfigPair: []*api.ConfigPair{pair},
		}

		res.ConfigList = append(res.ConfigList, conf)
	}

	return res, nil
}

// Ping Stub
func (service *CommonService) Ping(ctx context.Context, request *api.PingRequest) (*api.PongResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)

	res := &api.PongResponse{
		Message: "pong",
	}

	return res, nil
}

// Log Stub
func (service *CommonService) Log(ctx context.Context, request *api.LogRequest) (*api.LogResponse, error) {
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

	res := &api.LogResponse{}

	return res, nil
}

// Shutdown Stub
func (service *CommonService) Shutdown(ctx context.Context, request *api.ShutdownRequest) (*api.ShutdownResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)
	event := rk_inter_context.GetEvent(ctx)
	event.AddPair("signal", "interrupt")

	res := &api.ShutdownResponse{
		Message: "interrupt",
	}

	defer syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	return res, nil
}

// Info Stub
func (service *CommonService) Info(ctx context.Context, request *api.InfoRequest) (*api.InfoResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)
	event := rk_inter_context.GetEvent(ctx)

	res := &api.InfoResponse{}

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
func (service *CommonService) Healthy(ctx context.Context, request *api.HealthyRequest) (*api.HealthyResponse, error) {
	// Add auto generated request ID
	rk_inter_context.AddRequestIdToOutgoingMD(ctx)
	event := rk_inter_context.GetEvent(ctx)

	event.AddPair("healthy", "true")

	res := &api.HealthyResponse{
		Healthy: true,
	}

	return res, nil
}

func fillBasicInfo(res *api.InfoResponse) {
	basicInfo := &api.BasicInfo{
		AppName:   AppCtx.GetAppName(),
		StartTime: AppCtx.GetStartTime().Format(time.RFC3339),
		UpTime:    AppCtx.GetUpTime().String(),
		Realm:     os.Getenv("REALM"),
		Region:    os.Getenv("REGION"),
		Az:        os.Getenv("AZ"),
		Domain:    os.Getenv("DOMAIN"),
	}
	res.BasicInfo = basicInfo
}

func fillPromInfo(res *api.InfoResponse) {
	entry := AppCtx.GetPromEntry()
	if entry == nil {
		return
	}

	promInfo := &api.PromInfo{
		Port: strconv.FormatUint(entry.GetPort(), 10),
		Path: entry.GetPath(),
	}
	res.PromInfo = promInfo
}

func fillGRPCInfo(res *api.InfoResponse) {
	gRPCInfos := make([]*api.GRpcInfo, 0)
	gRPCEntries := AppCtx.ListGRpcEntries()
	for i := range gRPCEntries {
		entry := gRPCEntries[i]
		gRPCInfo := &api.GRpcInfo{
			Name: entry.GetName(),
			Port: strconv.FormatUint(entry.GetPort(), 10),
		}

		if entry.GetGWEntry() != nil {
			gwInfo := &api.GWInfo{
				GwPort: strconv.FormatUint(entry.GetGWEntry().GetHttpPort(), 10),
			}

			gRPCInfo.GwInfo = gwInfo
		}

		if entry.GetSWEntry() != nil {
			swInfo := &api.SWInfo{
				SwPath: entry.GetSWEntry().GetPath(),
				SwPort: strconv.FormatUint(entry.GetSWEntry().GetSWPort(), 10),
			}

			gRPCInfo.SwInfo = swInfo
		}

		gRPCInfos = append(gRPCInfos, gRPCInfo)
	}

	res.GrpcInfoList = gRPCInfos
}

func setLogLevel(config *zap.Config, level string) string {
	res := level
	if level == "debug" {
		config.Level.SetLevel(zapcore.DebugLevel)
	} else if level == "info" {
		config.Level.SetLevel(zapcore.InfoLevel)
	} else if level == "warn" {
		config.Level.SetLevel(zapcore.WarnLevel)
	} else if level == "error" {
		config.Level.SetLevel(zapcore.ErrorLevel)
	} else if level == "dpanic" {
		config.Level.SetLevel(zapcore.DPanicLevel)
	} else if level == "panic" {
		config.Level.SetLevel(zapcore.PanicLevel)
	} else if level == "fatal" {
		config.Level.SetLevel(zapcore.FatalLevel)
	} else {
		res = "unknown level, should be one of [info, warn,error, dpanic, panic, fatal]"
	}

	return res
}

func memStatsToZapFields(prefix string, stats *runtime.MemStats) []zap.Field {
	res := make([]zap.Field, 0)

	res = append(res,
		zap.Uint64(prefix+"MemAllocMB", bytesToMB(stats.Alloc)),
		zap.Uint64(prefix+"SysMemMB", bytesToMB(stats.Alloc)),
		zap.Time(prefix+"LastGCTimestamp", time.Unix(int64(stats.LastGC), 0)),
		zap.Uint32(prefix+"NumGC", stats.NumGC),
		zap.Uint32(prefix+"NumForceGC", stats.NumForcedGC),
	)

	return nil
}

func memStatsToPB(stats *runtime.MemStats) *api.MemStats {
	pb := api.MemStats{
		MemAllocMb:      bytesToMB(stats.Alloc),
		SysMemMb:        bytesToMB(stats.Sys),
		LastGcTimestamp: time.Unix(int64(stats.LastGC)/int64(time.Second), 0).Format(time.RFC3339),
		NumGc:           stats.NumGC,
		NumForceGc:      stats.NumForcedGC,
	}

	return &pb
}

func bytesToMB(b uint64) uint64 {
	return b / 1024 / 1024
}
