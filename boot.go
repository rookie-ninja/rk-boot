// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_boot

import (
	"errors"
	"fmt"
	rk_entry "github.com/rookie-ninja/rk-common/entry"
	rk_prom "github.com/rookie-ninja/rk-prom"

	//"github.com/rookie-ninja/rk-boot/prom"
	"github.com/rookie-ninja/rk-common/context"
	"github.com/rookie-ninja/rk-config"
	"github.com/rookie-ninja/rk-gin/boot"
	"github.com/rookie-ninja/rk-grpc/boot"
	"github.com/rookie-ninja/rk-query"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"os"
	"runtime/debug"
	"time"
)

type bootConfig struct {
	Application string `yaml:"application"`
	Event       struct {
		Format      string   `yaml:"format"`
		Quiet       bool     `yaml:"quiet"`
		OutputPaths []string `yaml:"outputPaths"`
		LoggerConf  string   `yaml:"loggerConf"`
		Maxsize     int      `yaml:"maxsize"`
		MaxAge      int      `yaml:"maxage"`
		MaxBackups  int      `yaml:"maxbackups"`
		Localtime   bool     `yaml:"localtime"`
		Compress    bool     `yaml:"compress"`
	} `yaml:"event"`
	Logger []struct {
		Name       string `yaml:"name"`
		Quiet      bool   `yaml:"quiet"`
		OutputPath string `yaml:"outputPath"`
		LoggerConf string `yaml:"loggerConf"`
		Maxsize    int    `yaml:"maxsize"`
		MaxAge     int    `yaml:"maxage"`
		MaxBackups int    `yaml:"maxbackups"`
		Localtime  bool   `yaml:"localtime"`
		Compress   bool   `yaml:"compress"`
	} `yaml:"logger"`
	Config []struct {
		Name   string `yaml:"name"`
		Path   string `yaml:"path"`
		Format string `yaml:"format"`
		Global bool   `yaml:"global"`
	} `yaml:"config"`
	Env struct {
		REALM  string `yaml:"REALM"`
		DOMAIN string `yaml:"DOMAIN"`
		REGION string `yaml:"REGION"`
		AZ     string `yaml:"AZ"`
	} `yaml:"env"`
}

type Boot struct {
	application  string
	startTime    time.Time
	viperConfigs map[string]*viper.Viper
	rkConfigs    map[string]*rk_config.RkConfig
	loggers      map[string]*rk_ctx.LoggerPair
	logger       *zap.Logger
	eventFactory *rk_query.EventFactory
}

type BootOption func(*Boot)

func WithBootConfigPath(filePath string) BootOption {
	return func(boot *Boot) {
		// read config
		bytes := readFile(filePath)
		config := &bootConfig{}
		err := yaml.Unmarshal(bytes, config)
		if err != nil {
			shutdownWithError(err)
		}

		// assign application name
		if len(config.Application) < 1 {
			config.Application = "unknown-application"
		}
		boot.application = config.Application

		// init logger
		boot.loggers = getLoggers(config)

		// init event factory
		boot.eventFactory = getEventFactory(config)

		// init config
		boot.viperConfigs, boot.rkConfigs = getConfigs(config)

		// init entries
		rk_grpc.NewGRpcEntries(filePath, boot.eventFactory, boot.logger)
		rk_gin.NewGinEntries(filePath, boot.eventFactory, boot.logger)

		rk_ctx.GlobalAppCtx.AddEntry(
			rk_prom.PromEntryNameDefault,
			rk_prom.NewPromEntryWithConfig(filePath, boot.eventFactory, boot.logger))
	}
}

func WithEventFactory(fac *rk_query.EventFactory) BootOption {
	return func(boot *Boot) {
		if fac != nil {
			boot.eventFactory = fac
		}
	}
}

func WithRkConfig(name string, config *rk_config.RkConfig) BootOption {
	return func(boot *Boot) {
		if len(name) < 1 {
			// with new name
			name = fmt.Sprintf("rkconfig-%d", len(boot.rkConfigs)+1)
		}

		boot.rkConfigs[name] = config
	}
}

func WithViperConfig(name string, config *viper.Viper) BootOption {
	return func(boot *Boot) {
		if len(name) < 1 {
			// with new name
			name = fmt.Sprintf("viper-%d", len(boot.viperConfigs)+1)
		}

		boot.viperConfigs[name] = config
	}
}

func WithLogger(name string, logger *zap.Logger, loggerConf *zap.Config) BootOption {
	return func(boot *Boot) {
		if len(name) < 1 {
			// with new name
			name = fmt.Sprintf("logger-%d", len(boot.viperConfigs)+1)
		}

		boot.loggers[name] = &rk_ctx.LoggerPair{
			Logger: logger,
			Config: loggerConf,
		}
	}
}

func NewBoot(opts ...BootOption) *Boot {
	// add default logger
	boot := &Boot{
		application:  "rk-application",
		eventFactory: rk_query.NewEventFactory(),
		viperConfigs: make(map[string]*viper.Viper),
		rkConfigs:    make(map[string]*rk_config.RkConfig),
		startTime:    time.Now(),
		logger:       rk_ctx.GlobalAppCtx.GetLogger("default"),
	}

	for i := range opts {
		opts[i](boot)
	}

	// basic information initiated, let's override app context
	rk_ctx.GlobalAppCtx.SetApplication(boot.application)
	rk_ctx.GlobalAppCtx.SetEventFactory(boot.eventFactory)
	for k, v := range boot.loggers {
		rk_ctx.GlobalAppCtx.AddLoggerPair(k, v)
	}
	rk_ctx.GlobalAppCtx.SetStartTime(boot.startTime)
	for k, v := range boot.viperConfigs {
		rk_ctx.GlobalAppCtx.AddViperConfig(k, v)
	}
	for k, v := range boot.rkConfigs {
		rk_ctx.GlobalAppCtx.AddRkConfig(k, v)
	}

	return boot
}

func (boot *Boot) GetApplication() string {
	return boot.application
}

func (boot *Boot) GetRkConfig(name string) *rk_config.RkConfig {
	res, _ := boot.rkConfigs[name]
	return res
}

func (boot *Boot) ListRkConfigs() []*rk_config.RkConfig {
	res := make([]*rk_config.RkConfig, 0)

	for _, v := range boot.rkConfigs {
		res = append(res, v)
	}
	return res
}

func (boot *Boot) GetViperConfig(name string) *viper.Viper {
	res, _ := boot.viperConfigs[name]
	return res
}

func (boot *Boot) ListViperConfigs() []*viper.Viper {
	res := make([]*viper.Viper, 0)

	for _, v := range boot.viperConfigs {
		res = append(res, v)
	}
	return res
}

func (boot *Boot) GetEntry(name string) rk_entry.Entry {
	return rk_ctx.GlobalAppCtx.GetEntry(name)
}

func (boot *Boot) GetGinEntry(name string) *rk_gin.GinEntry {
	return rk_ctx.GlobalAppCtx.GetEntry(name).(*rk_gin.GinEntry)
}

func (boot *Boot) GetGRpcEntry(name string) *rk_grpc.GRpcEntry {
	return rk_ctx.GlobalAppCtx.GetEntry(name).(*rk_grpc.GRpcEntry)
}

func (boot *Boot) GetPromEntry() *rk_prom.PromEntry {
	return rk_ctx.GlobalAppCtx.GetEntry(rk_prom.PromEntryNameDefault).(*rk_prom.PromEntry)
}

func (boot *Boot) Bootstrap() {
	helper := rk_query.NewEventHelper(boot.eventFactory)
	event := helper.Start("rk_app_start")
	defer helper.Finish(event)

	boot.startTime = time.Now()
	event.AddFields(zap.Time("app_start_time", boot.startTime))

	// bootstrap entries
	for _, entry := range rk_ctx.GlobalAppCtx.ListEntries() {
		entry.Bootstrap(event)
	}
}

func (boot *Boot) GetEventFactory() *rk_query.EventFactory {
	return boot.eventFactory
}

func (boot *Boot) RegisterQuitter(name string, input rk_ctx.QuitterFunc) error {
	if input == nil {
		return errors.New("empty quitter function")
	}

	if len(name) < 1 {
		return errors.New("empty quitter name")
	}

	if quitter := rk_ctx.GlobalAppCtx.GetQuitter(name); quitter != nil {
		return errors.New("duplicate quitter function")
	}

	rk_ctx.GlobalAppCtx.AddQuitter(name, input)

	return nil
}

func (boot *Boot) Wait(draining time.Duration) {
	sig := <-rk_ctx.GlobalAppCtx.GetShutdownSig()

	helper := rk_query.NewEventHelper(boot.eventFactory)
	event := helper.Start("rk_app_stop")

	// shutdown entries
	for _, entry := range rk_ctx.GlobalAppCtx.ListEntries() {
		entry.Shutdown(event)
	}

	boot.logger.Info("draining", zap.Duration("draining_duration", draining))
	time.Sleep(draining)

	event.AddFields(
		zap.Duration("app_lifetime_nano", time.Since(boot.startTime)),
		zap.Time("app_start_time", boot.startTime))

	event.AddPair("signal", sig.String())

	if err := recover(); err != nil {
		event.AddFields(zap.Any("recover", err))
		boot.logger.Error(fmt.Sprintf("err: %+v. stack: %s", err, string(debug.Stack())))
	}
	helper.Finish(event)
	os.Exit(0)
}
