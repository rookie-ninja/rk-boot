// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_boot

import (
	"fmt"
	rk_gin "github.com/rookie-ninja/rk-boot/gin"
	"github.com/rookie-ninja/rk-boot/grpc"
	"github.com/rookie-ninja/rk-boot/prom"
	"github.com/rookie-ninja/rk-config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"time"
)

type AppContext struct {
	boot *Boot
}

func NewAppContext(boot *Boot) *AppContext {
	return &AppContext{
		boot: boot,
	}
}

func (ctx *AppContext) AddLogger(name string, logger *zap.Logger, loggerConfig *zap.Config) string {
	if len(name) < 1 {
		name = fmt.Sprintf("logger-%d", len(ctx.boot.userLoggers)+1)
	}

	ctx.boot.userLoggers[name] = &rkLogger{
		logger: logger,
		config: loggerConfig,
	}
	return name
}

func (ctx *AppContext) AddRkConfig(name string, config *rk_config.RkConfig) string {
	if len(name) < 1 {
		name = fmt.Sprintf("rkconfig-%d", len(ctx.boot.rkConfigs)+1)
	}

	ctx.boot.rkConfigs[name] = config
	return name
}

func (ctx *AppContext) AddViperConfig(name string, config *viper.Viper) string {
	if len(name) < 1 {
		name = fmt.Sprintf("viper-%d", len(ctx.boot.viperConfigs)+1)
	}

	ctx.boot.viperConfigs[name] = config
	return name
}

func (ctx *AppContext) GetAppName() string {
	return ctx.boot.GetAppName()
}

func (ctx *AppContext) GetStartTime() time.Time {
	return ctx.boot.startTime
}

func (ctx *AppContext) GetUpTime() time.Duration {
	return time.Since(ctx.boot.startTime)
}

func (ctx *AppContext) GetRkConfig(name string) *viper.Viper {
	val, ok := ctx.boot.viperConfigs[name]
	if ok {
		return val
	}

	return viper.New()
}

func (ctx *AppContext) ListRkConfigs() map[string]*rk_config.RkConfig {
	return ctx.boot.rkConfigs
}

func (ctx *AppContext) GetViperConfig(name string) *viper.Viper {
	res, _ := ctx.boot.viperConfigs[name]
	return res
}

func (ctx *AppContext) ListViperConfigs() map[string]*viper.Viper {
	return ctx.boot.viperConfigs
}

func (ctx *AppContext) GetLogger(name string) *zap.Logger {
	val, ok := ctx.boot.userLoggers[name]
	if ok {
		return val.logger
	}

	return zap.NewNop()
}

func (ctx *AppContext) GetLoggerConfig(name string) *zap.Config {
	val, ok := ctx.boot.userLoggers[name]
	if ok {
		return val.config
	}

	return nil
}

func (ctx *AppContext) ListLoggers() map[string]*rkLogger {
	return ctx.boot.userLoggers
}

func (ctx *AppContext) ListGRpcEntries() []*rk_grpc.GRpcServerEntry {
	return ctx.boot.ListGRpcEntries()
}

func (ctx *AppContext) ListGinEntries() []*rk_gin.GinServerEntry {
	return ctx.boot.ListGinEntries()
}

func (ctx *AppContext) GetPromEntry() *rk_prom.PromEntry {
	return ctx.boot.promEntry
}

func (ctx *AppContext) GetShutdownSig() chan os.Signal {
	return ctx.boot.shutdownSig
}
