package rk_boot

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/rookie-ninja/rk-boot/api/v1"
	"github.com/rookie-ninja/rk-gin-interceptor/context"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

// GC Stub
func GC4Gin(ctx *gin.Context) {
	// Add auto generated request ID
	rk_gin_inter_context.AddRequestIdToOutgoingHeader(ctx)
	event := rk_gin_inter_context.GetEvent(ctx)

	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)

	runtime.GC()
	runtime.ReadMemStats(&after)

	event.AddFields(memStatsToZapFields("before_", &before)...)
	event.AddFields(memStatsToZapFields("after_", &before)...)

	pb := &rk_boot_common_v1.GCResponse{
		MemStatsBefore: memStatsToPB(&before),
		MemStatsAfter:  memStatsToPB(&after),
	}

	ctx.JSON(http.StatusOK, pb)
}

// DumpConfig Stub
func DumpConfig4Gin(ctx *gin.Context) {
	// Add auto generated request ID
	rk_gin_inter_context.AddRequestIdToOutgoingHeader(ctx)

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

	ctx.JSON(http.StatusOK, res)
}

// GetConfig Stub
func GetConfig4Gin(ctx *gin.Context) {
	// Add auto generated request ID
	rk_gin_inter_context.AddRequestIdToOutgoingHeader(ctx)

	key := ctx.Param("key")

	configList := make([]*rk_boot_common_v1.Config, 0)
	res := &rk_boot_common_v1.GetConfigResponse{ConfigList: configList}

	for k, v := range AppCtx.ListRkConfigs() {
		pair := &rk_boot_common_v1.ConfigPair{
			Key:   key,
			Value: cast.ToString(v.Get(key)),
		}

		conf := &rk_boot_common_v1.Config{
			ConfigName: k,
			ConfigPair: []*rk_boot_common_v1.ConfigPair{pair},
		}

		res.ConfigList = append(res.ConfigList, conf)
	}

	for k, v := range AppCtx.ListViperConfigs() {
		pair := &rk_boot_common_v1.ConfigPair{
			Key:   key,
			Value: cast.ToString(v.Get(key)),
		}

		conf := &rk_boot_common_v1.Config{
			ConfigName: k,
			ConfigPair: []*rk_boot_common_v1.ConfigPair{pair},
		}

		res.ConfigList = append(res.ConfigList, conf)
	}

	ctx.JSON(http.StatusOK, res)
}

// Ping Stub
func Ping4Gin(ctx *gin.Context) {
	// Add auto generated request ID
	rk_gin_inter_context.AddRequestIdToOutgoingHeader(ctx)

	res := &rk_boot_common_v1.PongResponse{
		Message: "pong",
	}

	ctx.JSON(http.StatusOK, res)
}

// Log Stub
func Log4Gin(ctx *gin.Context) {
	// Add auto generated request ID
	rk_gin_inter_context.AddRequestIdToOutgoingHeader(ctx)
	event := rk_gin_inter_context.GetEvent(ctx)

	bytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to unmarshal request body",
		})
		return
	}

	var request *rk_boot_common_v1.LogRequest
	if err := proto.Unmarshal(bytes, request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to unmarshal request body",
		})
		return
	}

	for i := range request.Entries {
		entry := request.Entries[i]
		loggerConfig := AppCtx.GetLoggerConfig(entry.LogName)
		if loggerConfig != nil {
			setLogLevel(loggerConfig, entry.LogLevel)
		}

		event.AddPair(entry.LogName, entry.LogLevel)
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// Shutdown Stub
func Shutdown4Gin(ctx *gin.Context) {
	// Add auto generated request ID
	rk_gin_inter_context.AddRequestIdToOutgoingHeader(ctx)
	event := rk_gin_inter_context.GetEvent(ctx)
	event.AddPair("signal", "interrupt")

	res := &rk_boot_common_v1.ShutdownResponse{
		Message: "interrupt",
	}

	ctx.JSON(http.StatusOK, res)

	defer func() {
		AppCtx.GetShutdownSig() <- syscall.SIGINT
	}()
}

// Info Stub
func Info4Gin(ctx *gin.Context) {
	// Add auto generated request ID
	rk_gin_inter_context.AddRequestIdToOutgoingHeader(ctx)

	res := &rk_boot_common_v1.InfoResponse{}

	// Basic info
	fillBasicInfo(res)

	// gRPC info
	fillGRPCInfo(res)

	// gin info
	fillGinInfo(res)

	// Prom info
	fillPromInfo(res)

	ctx.JSON(http.StatusOK, res)
}

// Healthy Stub
func Healthy4Gin(ctx *gin.Context) {
	// Add auto generated request ID
	rk_gin_inter_context.AddRequestIdToOutgoingHeader(ctx)
	event := rk_gin_inter_context.GetEvent(ctx)

	event.AddPair("healthy", "true")

	res := &rk_boot_common_v1.HealthyResponse{
		Healthy: true,
	}

	ctx.JSON(http.StatusOK, res)
}

func fillBasicInfo(res *rk_boot_common_v1.InfoResponse) {
	basicInfo := &rk_boot_common_v1.BasicInfo{
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

func fillPromInfo(res *rk_boot_common_v1.InfoResponse) {
	entry := AppCtx.GetPromEntry()
	if entry == nil {
		return
	}

	promInfo := &rk_boot_common_v1.PromInfo{
		Port: strconv.FormatUint(entry.GetPort(), 10),
		Path: entry.GetPath(),
	}
	res.PromInfo = promInfo
}

func fillGRPCInfo(res *rk_boot_common_v1.InfoResponse) {
	gRPCInfos := make([]*rk_boot_common_v1.GRpcInfo, 0)
	gRPCEntries := AppCtx.ListGRpcEntries()
	for i := range gRPCEntries {
		entry := gRPCEntries[i]
		gRPCInfo := &rk_boot_common_v1.GRpcInfo{
			Name: entry.GetName(),
			Port: strconv.FormatUint(entry.GetPort(), 10),
		}

		if entry.GetGWEntry() != nil {
			gwInfo := &rk_boot_common_v1.GWInfo{
				GwPort: strconv.FormatUint(entry.GetGWEntry().GetHttpPort(), 10),
			}

			gRPCInfo.GwInfo = gwInfo
		}

		if entry.GetSWEntry() != nil {
			swInfo := &rk_boot_common_v1.SWInfo{
				SwPath: entry.GetSWEntry().GetPath(),
				SwPort: strconv.FormatUint(entry.GetSWEntry().GetSWPort(), 10),
			}

			gRPCInfo.SwInfo = swInfo
		}

		gRPCInfos = append(gRPCInfos, gRPCInfo)
	}

	res.GrpcInfoList = gRPCInfos
}

func fillGinInfo(res *rk_boot_common_v1.InfoResponse) {
	ginInfos := make([]*rk_boot_common_v1.GinInfo, 0)
	ginEntries := AppCtx.ListGinEntries()
	for i := range ginEntries {
		entry := ginEntries[i]
		ginInfo := &rk_boot_common_v1.GinInfo{
			Name: entry.GetName(),
			Port: strconv.FormatUint(entry.GetPort(), 10),
		}

		if entry.GetSWEntry() != nil {
			swInfo := &rk_boot_common_v1.SWInfo{
				SwPath: entry.GetSWEntry().GetPath(),
				SwPort: strconv.FormatUint(entry.GetSWEntry().GetSWPort(), 10),
			}

			ginInfo.SwInfo = swInfo
		}

		ginInfos = append(ginInfos, ginInfo)
	}

	res.GinInfoList = ginInfos
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

func memStatsToPB(stats *runtime.MemStats) *rk_boot_common_v1.MemStats {
	pb := rk_boot_common_v1.MemStats{
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
