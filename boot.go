// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_boot

import (
	"errors"
	"fmt"
	"github.com/rookie-ninja/rk-boot/gin"
	"github.com/rookie-ninja/rk-boot/grpc"
	"github.com/rookie-ninja/rk-boot/prom"
	"github.com/rookie-ninja/rk-config"
	"github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"path"
	"runtime/debug"
	"syscall"
	"time"
)

var AppCtx = NewAppContext(nil)

type bootConfig struct {
	AppName string `yaml:"appName"`
	dir     string
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
		Tls  struct {
			Enabled bool   `yaml:"enabled"`
			Port    uint64 `yaml:"port"`
			User    struct {
				Enabled  bool   `yaml:"enabled"`
				CertFile string `yaml:"certFile"`
				KeyFile  string `yaml:"keyFile"`
			} `yaml:"user"`
			Auto struct {
				Enabled    bool   `yaml:"enabled"`
				CertOutput string `yaml:"certOutput"`
			} `yaml:"auto"`
		} `yaml:"tls"`
		GW                  struct {
			Enabled             bool   `yaml:"enabled"`
			Port                uint64 `yaml:"port"`
			EnableCommonService bool   `yaml:"enableCommonService"`
		} `yaml:"gw"`
		SW struct {
			Enabled             bool     `yaml:"enabled"`
			Port                uint64   `yaml:"port"`
			Path                string   `yaml:"path"`
			JsonPath            string   `yaml:"jsonPath"`
			EnableCommonService bool     `yaml:"enableCommonService"`
			Headers             []string `yaml:"headers"`
		} `yaml:"sw"`
		LoggingInterceptor struct {
			Enabled              bool `yaml:"enabled"`
			EnableLogging        bool `yaml:"enableLogging"`
			EnableMetrics        bool `yaml:"enableMetrics"`
			EnablePayloadLogging bool `yaml:"enablePayloadLogging"`
		} `yaml:"loggingInterceptor"`
	} `yaml:"grpc"`
	Gin []struct {
		Name string `yaml:"name"`
		Port uint64 `yaml:"port"`
		Tls  struct {
			Enabled bool   `yaml:"enabled"`
			Port    uint64 `yaml:"port"`
			User    struct {
				Enabled  bool   `yaml:"enabled"`
				CertFile string `yaml:"certFile"`
				KeyFile  string `yaml:"keyFile"`
			} `yaml:"user"`
			Auto struct {
				Enabled    bool   `yaml:"enabled"`
				CertOutput string `yaml:"certOutput"`
			} `yaml:"auto"`
		} `yaml:"tls"`
		SW struct {
			Enabled  bool     `yaml:"enabled"`
			Path     string   `yaml:"path"`
			JsonPath string   `yaml:"jsonPath"`
			Insecure bool     `yaml:"insecure"`
			Headers  []string `yaml:"headers"`
		} `yaml:"sw"`
		EnableCommonService bool `yaml:"enableCommonService"`
		LoggingInterceptor  struct {
			Enabled       bool `yaml:"enabled"`
			EnableLogging bool `yaml:"enableLogging"`
			EnableMetrics bool `yaml:"enableMetrics"`
		} `yaml:"loggingInterceptor"`
		AuthInterceptor struct {
			Enabled     bool     `yaml:"enabled"`
			Realm       string   `yaml:"realm"`
			Credentials []string `yaml:"credentials"`
		} `yaml:"authInterceptor"`
	} `yaml:"gin"`
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
	ginServerEntry  map[string]*rk_gin.GinServerEntry
	promEntry       *rk_prom.PromEntry
	bootLogger      *zap.Logger
	eventLogger     *zap.Logger
	userLoggers     map[string]*rkLogger
	eventFactory    *rk_query.EventFactory
	viperConfigs    map[string]*viper.Viper
	rkConfigs       map[string]*rk_config.RkConfig
	shutdownSig     chan os.Signal
}

type BootOption func(*Boot)

func WithBootConfigPath(filePath string) BootOption {
	return func(boot *Boot) {
		// read config
		bytes, ext := readFile(filePath)
		config := &bootConfig{}
		unMarshal(bytes, ext, config)

		// assign config path
		config.dir = path.Dir(filePath)

		boot.appName = config.AppName

		// init logger
		boot.userLoggers, boot.bootLogger, boot.eventLogger = getLoggers(config)

		// init event factory
		boot.eventFactory = getEventFactory(config, boot.eventLogger)

		// init config
		boot.viperConfigs, boot.rkConfigs = getConfigs(config)

		// init gRpc
		boot.gRpcServerEntry = getGRpcServerEntries(config, boot.eventFactory)

		// init gin
		boot.ginServerEntry = getGinServerEntries(config, boot.eventFactory, boot.bootLogger)

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
	config := &bootConfig{
		AppName: "unknown",
	}
	defaultBootLogger, _, _ := rk_logger.NewZapLoggerWithBytes([]byte(defaultZapConfig), rk_logger.YAML)
	defaultEventLogger, _, _ := rk_logger.NewZapLoggerWithBytes([]byte(defaultZapConfigEvent), rk_logger.YAML)

	boot := &Boot{
		appName:         "unknown",
		gRpcServerEntry: make(map[string]*rk_grpc.GRpcServerEntry),
		userLoggers:     make(map[string]*rkLogger),
		bootLogger:      defaultBootLogger,
		eventFactory:    getEventFactory(config, defaultEventLogger),
		shutdownSig:     make(chan os.Signal),
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

func (boot *Boot) GetGinEntry(name string) *rk_gin.GinServerEntry {
	res, _ := boot.ginServerEntry[name]
	return res
}

func (boot *Boot) ListGRpcEntries() []*rk_grpc.GRpcServerEntry {
	res := make([]*rk_grpc.GRpcServerEntry, 0)
	for _, v := range boot.gRpcServerEntry {
		res = append(res, v)
	}

	return res
}

func (boot *Boot) ListGinEntries() []*rk_gin.GinServerEntry {
	res := make([]*rk_gin.GinServerEntry, 0)
	for _, v := range boot.ginServerEntry {
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
	helper := rk_query.NewEventHelper(boot.eventFactory)
	event := helper.Start("rk_app_start")
	defer helper.Finish(event)

	boot.startTime = time.Now()
	event.AddFields(zap.Time("app_start_time", boot.startTime))

	// gRpc, gateway, swagger
	for _, entry := range boot.gRpcServerEntry {
		entry.Start(boot.bootLogger)
		event.AddFields(zap.Uint64(fmt.Sprintf("%s_gRpc_port", entry.GetName()), entry.GetPort()))

		entry.StartGW(boot.bootLogger)
		event.AddFields(zap.Uint64(fmt.Sprintf("%s_gw_port", entry.GetName()), entry.GetGWEntry().GetHttpPort()))

		entry.StartSW(boot.bootLogger)
		event.AddFields(
			zap.Uint64(fmt.Sprintf("%s_sw_port", entry.GetName()), entry.GetSWEntry().GetSWPort()),
			zap.String(fmt.Sprintf("%s_sw_path", entry.GetName()), entry.GetSWEntry().GetPath()))
	}

	// gin, swagger
	for _, entry := range boot.ginServerEntry {
		entry.Start(boot.bootLogger)
		event.AddFields(
			zap.Uint64(fmt.Sprintf("%s_gin_port", entry.GetName()), entry.GetPort()),
			zap.Uint64(fmt.Sprintf("%s_sw_port", entry.GetName()), entry.GetSWEntry().GetSWPort()),
			zap.String(fmt.Sprintf("%s_sw_path", entry.GetName()), entry.GetSWEntry().GetPath()))
	}

	if boot.promEntry != nil {
		event.AddFields(
			zap.Uint64("prom_port", boot.promEntry.GetPort()),
			zap.String("prom_path", boot.promEntry.GetPath()))
		boot.promEntry.Start(boot.bootLogger)
	}

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

	helper := rk_query.NewEventHelper(boot.eventFactory)
	event := helper.Start("rk_app_stop")

	// shutdown gRpc, gateway, swagger entry
	for _, entry := range boot.gRpcServerEntry {
		event.AddFields(zap.Uint64(fmt.Sprintf("%s_gRpc_port", entry.GetName()), entry.GetPort()))
		entry.Stop(boot.bootLogger.With(zap.Any("signal", sig)))

		event.AddFields(zap.Uint64(fmt.Sprintf("%s_gw_port", entry.GetName()), entry.GetGWEntry().GetHttpPort()))
		entry.StopGW(boot.bootLogger.With(zap.Any("signal", sig)))

		event.AddFields(
			zap.Uint64(fmt.Sprintf("%s_sw_port", entry.GetName()), entry.GetSWEntry().GetSWPort()),
			zap.String(fmt.Sprintf("%s_sw_path", entry.GetName()), entry.GetSWEntry().GetPath()))
		entry.StopSW(boot.bootLogger.With(zap.Any("signal", sig)))
	}

	// shutdown gin entry
	for _, entry := range boot.ginServerEntry {
		event.AddFields(zap.Uint64(fmt.Sprintf("%s_gin_port", entry.GetName()), entry.GetPort()))
		entry.Stop(boot.bootLogger.With(zap.Any("signal", sig)))

		event.AddFields(
			zap.Uint64(fmt.Sprintf("%s_sw_port", entry.GetName()), entry.GetSWEntry().GetSWPort()),
			zap.String(fmt.Sprintf("%s_sw_path", entry.GetName()), entry.GetSWEntry().GetPath()))
	}

	if boot.promEntry != nil {
		event.AddFields(
			zap.Uint64("prom_port", boot.promEntry.GetPort()),
			zap.String("prom_path", boot.promEntry.GetPath()))
		boot.promEntry.Stop(boot.bootLogger)
	}

	boot.bootLogger.Info("draining", zap.Duration("draining_duration", draining))
	time.Sleep(draining)

	event.AddFields(
		zap.Duration("app_lifetime_nano", time.Since(boot.startTime)),
		zap.Time("app_start_time", boot.startTime))

	event.AddPair("signal", sig.String())

	if err := recover(); err != nil {
		event.AddFields(zap.Any("recover", err))
		boot.bootLogger.Error(fmt.Sprintf("err: %+v. stack: %s", err, string(debug.Stack())))
	}
	helper.Finish(event)
	os.Exit(0)
}

func (boot *Boot) shutdownHook() chan os.Signal {
	signal.Notify(boot.shutdownSig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	return boot.shutdownSig
}
