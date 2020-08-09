// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_boot

import (
	"errors"
	"fmt"
	"github.com/rookie-ninja/rk-boot/grpc"
	"github.com/rookie-ninja/rk-boot/prom"
	"github.com/rookie-ninja/rk-config"
	"github.com/rookie-ninja/rk-query"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var AppCtx = NewAppContext(nil)

type bootConfig struct {
	AppName string `yaml:"appName"`
	Event   struct {
		Format string `yaml:"format"`
		Quiet  bool   `yaml:"quiet"`
	} `yaml:"event"`
	Logger []struct {
		Name     string `yaml:"name"`
		ConfPath string `yaml:"confPath"`
		ForBoot  bool   `yaml:"forBoot,omitempty"`
		ForEvent bool   `yaml:"forEvent,omitempty"`
	} `yaml:"logger"`
	Config []struct {
		Name   string `yaml:"name"`
		Path   string `yaml:"path"`
		Format string `yaml:"format"`
		Global bool   `yaml:"global"`
	} `yaml:"config"`
	GRpc []struct {
		Name                string `yaml:"name"`
		Port                uint64 `yaml:"port"`
		EnableCommonService bool   `yaml:"enableCommonService"`
		GW                  struct {
			Enabled             bool   `yaml:"enabled"`
			Port                uint64 `yaml:"port"`
			Insecure            bool   `yaml:"insecure"`
			EnableCommonService bool   `yaml:"enableCommonService"`
		} `yaml:"gw"`
		SW struct {
			Enabled             bool   `yaml:"enabled"`
			Port                uint64 `yaml:"port"`
			Path                string `yaml:"path"`
			JsonPath            string `yaml:"jsonPath"`
			Insecure            bool   `yaml:"insecure"`
			EnableCommonService bool   `yaml:"enableCommonService"`
		} `yaml:"sw"`
		LoggingInterceptor struct {
			Enabled              bool `yaml:"enabled"`
			EnableLogging        bool `yaml:"enableLogging"`
			EnableMetrics        bool `yaml:"enableMetrics"`
			EnablePayloadLogging bool `yaml:"enablePayloadLogging"`
		} `yaml:"loggingInterceptor"`
	} `yaml:"grpc"`
	Prom struct {
		Enabled     bool   `yaml:"enabled"`
		Port        uint64 `yaml:"port"`
		Path        string `yaml:"path"`
		PushGateway struct {
			Enabled    bool   `yaml:"enabled"`
			RemoteAddr string `yaml:"remoteAddr"`
			IntervalMS uint64 `yaml:"intervalMS"`
			JobName    string `yaml:"jobName"`
		} `yaml:"prom"`
	} `yaml:"prom"`
}

type QuitterFunc func()

type rkLogger struct {
	logger *zap.Logger
	config *zap.Config
}

type Boot struct {
	appName         string
	quitters        map[string]QuitterFunc
	startTime       time.Time
	gRpcServerEntry map[string]*rk_grpc.GRpcServerEntry
	promEntry       *rk_prom.PromEntry
	bootLogger      *zap.Logger
	eventLogger     *zap.Logger
	userLoggers     map[string]*rkLogger
	eventFactory    *rk_query.EventFactory
	viperConfigs    map[string]*viper.Viper
	rkConfigs       map[string]*rk_config.RkConfig
}

type BootOption func(*Boot)

func WithBootConfigPath(filePath string) BootOption {
	return func(boot *Boot) {
		// read config
		bytes, ext := readFile(filePath)
		config := &bootConfig{}
		unMarshal(bytes, ext, config)

		boot.appName = config.AppName

		// init logger
		boot.userLoggers, boot.bootLogger, boot.eventLogger = getLoggers(config)

		// init event factory
		boot.eventFactory = getEventFactory(config, boot.eventLogger)

		// init config
		boot.viperConfigs, boot.rkConfigs = getConfigs(config)

		// init gRpc
		boot.gRpcServerEntry = getGRpcServerEntries(config, boot.eventFactory)

		// init prom
		boot.promEntry = getPromEntry(config)
	}
}

func WithBootLogger(logger *zap.Logger) BootOption {
	return func(boot *Boot) {
		if logger != nil {
			boot.bootLogger = logger
		}
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

		boot.userLoggers[name] = &rkLogger{
			logger: logger,
			config: loggerConf,
		}
	}
}

func NewBoot(opts ...BootOption) *Boot {
	boot := &Boot{
		gRpcServerEntry: make(map[string]*rk_grpc.GRpcServerEntry),
		userLoggers:     make(map[string]*rkLogger),
		bootLogger:      zap.NewNop(),
		eventFactory:    rk_query.NewEventFactory(),
	}

	for i := range opts {
		opts[i](boot)
	}

	return boot
}

func (boot *Boot) GetAppName() string {
	return boot.appName
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

func (boot *Boot) GetGRpcEntry(name string) *rk_grpc.GRpcServerEntry {
	res, _ := boot.gRpcServerEntry[name]
	return res
}

func (boot *Boot) ListGRpcEntries() []*rk_grpc.GRpcServerEntry {
	res := make([]*rk_grpc.GRpcServerEntry, 0)
	for _, v := range boot.gRpcServerEntry {
		res = append(res, v)
	}

	return res
}

func (boot *Boot) RegisterGRpcServer(entry *rk_grpc.GRpcServerEntry) {
	if entry == nil {
		return
	}

	boot.gRpcServerEntry[entry.GetName()] = entry
}

func (boot *Boot) RegisterPromEntry(entry *rk_prom.PromEntry) {
	if entry == nil {
		return
	}

	boot.promEntry = entry
}

func (boot *Boot) Bootstrap() *AppContext {
	// gRpc, gateway, swagger
	for _, v := range boot.gRpcServerEntry {
		v.Start(boot.bootLogger)

		v.StartGW(boot.bootLogger)

		v.StartSW(boot.bootLogger)
	}

	if boot.promEntry != nil {
		boot.promEntry.Start(boot.bootLogger)
	}

	boot.startTime = time.Now()

	AppCtx.boot = boot

	return AppCtx
}

func (boot *Boot) RegisterQuitter(name string, input QuitterFunc) error {
	if input == nil {
		return errors.New("empty quitter function")
	}

	if len(name) < 1 {
		return errors.New("empty quitter name")
	}

	_, contains := boot.quitters[name]

	if contains {
		return errors.New("duplicate quitter function")
	}

	boot.quitters[name] = input

	return nil
}

func (boot *Boot) Quitter(draining time.Duration) {
	sig := <-boot.shutdownHook()

	for _, entry := range boot.gRpcServerEntry {
		entry.Stop(boot.bootLogger.With(zap.Any("signal", sig)))
		entry.StopGW(boot.bootLogger.With(zap.Any("signal", sig)))
	}

	if boot.promEntry != nil {
		boot.promEntry.Stop(boot.bootLogger)
	}

	boot.bootLogger.Info("draining", zap.Duration("draining_duration", draining))
	time.Sleep(draining)

	os.Exit(0)
}

func (boot *Boot) shutdownHook() chan os.Signal {
	shutdownHook := make(chan os.Signal)
	signal.Notify(shutdownHook,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt)
	return shutdownHook
}
