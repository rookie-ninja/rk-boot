// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_boot

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rookie-ninja/rk-boot/api/v1"
	"github.com/rookie-ninja/rk-boot/grpc"
	"github.com/rookie-ninja/rk-boot/gw"
	"github.com/rookie-ninja/rk-boot/prom"
	"github.com/rookie-ninja/rk-boot/sw"
	"github.com/rookie-ninja/rk-config"
	"github.com/rookie-ninja/rk-interceptor/logging/zap"
	"github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

func shutdownWithError(err error) {
	logger, _ := zap.NewDevelopment()
	logger.Error(fmt.Sprintf("%v", err))
	os.Exit(1)
}

func readFile(filePath string) ([]byte, string) {
	if !path.IsAbs(filePath) {
		wd, err := os.Getwd()
		if err != nil {
			shutdownWithError(err)
		}
		filePath = path.Join(wd, filePath)
	}

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		shutdownWithError(err)
	}

	ext := path.Ext(filePath)

	return bytes, ext
}

func unMarshal(bytes []byte, ext string, target interface{}) {
	if ext == ".yaml" || ext == ".yml" {
		err := yaml.Unmarshal(bytes, target)
		if err != nil {
			shutdownWithError(err)
		}
	} else if ext == ".json" {
		err := json.Unmarshal(bytes, target)
		if err != nil {
			shutdownWithError(err)
		}
	} else {
		shutdownWithError(errors.New(fmt.Sprintf("unsupported file type::%s", ext)))
	}
}

func getEventFactory(config *bootConfig, logger *zap.Logger) *rk_query.EventFactory {
	return rk_query.NewEventFactory(
		rk_query.WithAppName(config.AppName),
		rk_query.WithFormat(rk_query.ToFormat(config.Event.Format)),
		rk_query.WithQuietMode(config.Event.Quiet),
		rk_query.WithLogger(logger))
}

func getLoggers(config *bootConfig) (map[string]*rkLogger, *zap.Logger, *zap.Logger) {
	res := make(map[string]*rkLogger)
	bootLogger := rk_logger.StdoutLogger
	eventLogger := rk_query.StdoutLogger

	if config.Logger == nil {
		return res, bootLogger, eventLogger
	}

	for i := range config.Logger {
		bytes, ext := readFile(config.Logger[i].ConfPath)
		logger, loggerConf, err := rk_logger.NewZapLoggerWithBytes(bytes, rk_logger.ToFileType(ext))
		if err != nil {
			shutdownWithError(err)
		}

		name := config.Logger[i].Name
		if len(name) < 1 {
			name = fmt.Sprintf("logger-%d", len(res)+1)
		}

		res[name] = &rkLogger{
			logger: logger,
			config: loggerConf,
		}

		if config.Logger[i].ForBoot {
			bootLogger = logger
		}

		if config.Logger[i].ForEvent {
			eventLogger = logger
		}
	}

	return res, bootLogger, eventLogger
}

func getConfigs(config *bootConfig) (map[string]*viper.Viper, map[string]*rk_config.RkConfig) {
	vipers := make(map[string]*viper.Viper)
	rks := make(map[string]*rk_config.RkConfig)

	for i := range config.Config {
		element := config.Config[i]
		name := element.Name
		if len(name) < 1 {
			name = uuid.New().String()
		}

		if element.Format == "RK" {
			rks[name] = getRkConfig(element.Path, element.Global)
		} else {
			vipers[name] = getViperConfig(element.Path, element.Global)
		}
	}

	return vipers, rks
}

func getPromEntry(config *bootConfig) *rk_prom.PromEntry {
	var res *rk_prom.PromEntry
	if config.Prom.Enabled {
		var pgwRemoteAddr string
		var pgwIntervalMS uint64
		if config.Prom.PushGateway.Enabled {
			pgwRemoteAddr = config.Prom.PushGateway.RemoteAddr
			pgwIntervalMS = config.Prom.PushGateway.IntervalMS
		}

		res = rk_prom.NewPromEntry(
			rk_prom.WithPort(config.Prom.Port),
			rk_prom.WithPath(config.Prom.Path),
			rk_prom.WithPGWRemoteAddr(pgwRemoteAddr),
			rk_prom.WithPGWIntervalMS(pgwIntervalMS),
		)
	}

	return res
}

func getGRpcServerEntries(config *bootConfig, eventFactory *rk_query.EventFactory) map[string]*rk_grpc.GRpcServerEntry {
	res := make(map[string]*rk_grpc.GRpcServerEntry)

	for i := range config.GRpc {
		element := config.GRpc[i]
		name := element.Name

		// did we enabled gateway?
		var gwEntry *rk_gw.GRpcGWEntry
		if element.GW.Enabled {
			opts := make([]grpc.DialOption, 0)
			if element.GW.Insecure {
				opts = append(opts, grpc.WithInsecure())
			}

			gwEntry = rk_gw.NewGRpcGWEntry(
				rk_gw.WithHttpPort(element.GW.Port),
				rk_gw.WithGRpcPort(element.Port),
				rk_gw.WithDialOptions(opts...),
				rk_gw.WithCommonService(element.GW.EnableCommonService))
		}

		// did we enabled swagger?
		var swEntry *rk_sw.SWEntry
		if element.SW.Enabled {
			opts := make([]grpc.DialOption, 0)
			if element.SW.Insecure {
				opts = append(opts, grpc.WithInsecure())
			}

			swEntry = rk_sw.NewSWEntry(
				rk_sw.WithGRpcPort(element.Port),
				rk_sw.WithSWPort(element.SW.Port),
				rk_sw.WithPath(element.SW.Path),
				rk_sw.WithJsonPath(element.SW.JsonPath),
				rk_sw.WithDialOptions(opts...),
				rk_sw.WithCommonService(element.SW.EnableCommonService))
		}

		entry := rk_grpc.NewGRpcServerEntry(
			rk_grpc.WithName(name),
			rk_grpc.WithPort(element.Port),
			rk_grpc.WithGWEntry(gwEntry),
			rk_grpc.WithRegFuncs(registerRkCommonServiceGRPC),
			rk_grpc.WithSWEntry(swEntry),
			rk_grpc.WithCommonService(element.EnableCommonService))

		// did we enabled logging interceptor?
		if element.LoggingInterceptor.Enabled {
			opts := make([]rk_inter_logging.Option, 0)
			if !element.LoggingInterceptor.EnableLogging {
				opts = append(opts, rk_inter_logging.EnableLoggingOption(rk_inter_logging.DisableLogging))
			}

			if !element.LoggingInterceptor.EnableMetrics {
				opts = append(opts, rk_inter_logging.EnableLoggingOption(rk_inter_logging.DisableMetrics))
			}

			if !element.LoggingInterceptor.EnablePayloadLogging {
				opts = append(opts, rk_inter_logging.EnablePayloadLoggingOption(rk_inter_logging.DisablePayloadLogging))
			}

			entry.AddUnaryInterceptors(rk_inter_logging.UnaryServerInterceptor(eventFactory, opts...))
			entry.AddStreamInterceptors(rk_inter_logging.StreamServerInterceptor(eventFactory, opts...))
		}

		res[name] = entry
	}

	return res
}

func getRkConfig(path string, global bool) *rk_config.RkConfig {
	if len(path) < 1 {
		shutdownWithError(errors.New("empty config path"))
	}

	if global {
		rkConfig, err := rk_config.NewRkConfigGlobal(path)
		if err != nil {
			shutdownWithError(errors.New("empty config path"))
		}

		return rkConfig
	} else {
		rkConfig, err := rk_config.NewRkConfig(path)
		if err != nil {
			shutdownWithError(errors.New("empty config path"))
		}

		return rkConfig
	}
}

func getViperConfig(path string, global bool) *viper.Viper {
	if len(path) < 1 {
		if global {
			return viper.GetViper()
		} else {
			return viper.New()
		}
	}

	if global {
		viperConfig, err := rk_config.NewViperConfigGlobal(path)
		if err != nil {
			shutdownWithError(errors.New("empty config path"))
		}

		return viperConfig
	} else {
		viperConfig, err := rk_config.NewViperConfig(path)
		if err != nil {
			shutdownWithError(errors.New("empty config path"))
		}

		return viperConfig
	}
}

// Register common service
func registerRkCommonServiceGRPC(server *grpc.Server) {
	rk_boot_common_v1.RegisterRkCommonServiceServer(server, NewCommonService())
}
