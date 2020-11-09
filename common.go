// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_boot

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rookie-ninja/rk-common/context"
	"github.com/rookie-ninja/rk-config"
	"github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

func getEventFactory(config *bootConfig) *rk_query.EventFactory {
	fields := make([]zap.Field, 0)

	// add username, uid and gid to every event data by default
	if u, err := user.Current(); err == nil {
		fields = append(fields,
			zap.String("user", u.Username),
			zap.String("uid", u.Uid),
			zap.String("gid", u.Gid))
	}

	logger := rk_query.StdoutLogger

	// logger config specified, let's init it first
	if len(config.Event.LoggerConf) > 0 {
		if !filepath.IsAbs(config.Event.LoggerConf) {
			wd, _ := os.Getwd()
			config.Event.LoggerConf = path.Join(wd, config.Event.LoggerConf)
		}

		// merge logger config and build logger
		logger = mergeLogger(readFile(config.Event.LoggerConf), config.Event)
	} else {
		// merge logger config and build logger
		logger = mergeLogger(rk_query.StdLoggerConfigBytes, config.Event)
	}

	return rk_query.NewEventFactory(
		rk_query.WithAppName(config.Application),
		rk_query.WithFormat(rk_query.ToFormat(config.Event.Format)),
		rk_query.WithQuietMode(config.Event.Quiet),
		rk_query.WithLogger(logger),
		rk_query.WithFields(fields))
}

func getLoggers(config *bootConfig) map[string]*rk_ctx.LoggerPair {
	res := map[string]*rk_ctx.LoggerPair{
		"default": {
			Logger: rk_logger.StdoutLogger,
			Config: &rk_logger.StdoutLoggerConfig,
		},
	}

	if config.Logger == nil {
		return res
	}

	for i := range config.Logger {
		confPath := config.Logger[i].LoggerConf
		if !filepath.IsAbs(confPath) {
			wd, _ := os.Getwd()
			confPath = path.Join(wd, confPath)
		}

		bytes := readFile(confPath)
		logger, loggerConf, err := rk_logger.NewZapLoggerWithBytes(bytes, rk_logger.YAML)
		if err != nil {
			shutdownWithError(err)
		}

		name := config.Logger[i].Name
		if len(name) < 1 {
			name = fmt.Sprintf("logger-%d", len(res)+1)
		}

		res[name] = &rk_ctx.LoggerPair{
			Logger: logger,
			Config: loggerConf,
		}
	}

	return res
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

		if !path.IsAbs(element.Path) {
			wd, _ := os.Getwd()
			element.Path = path.Join(wd, element.Path)
		}

		if element.Format == "RK" {
			rks[name] = getRkConfig(element.Path, element.Global)
		} else {
			vipers[name] = getViperConfig(element.Path, element.Global)
		}
	}

	return vipers, rks
}

//func getPromEntry(config *bootConfig) *rk_prom.PromEntry {
//	var res *rk_prom.PromEntry
//	if config.Prom.Enabled {
//		var pgwRemoteAddr string
//		var pgwIntervalMS uint64
//		if config.Prom.PushGateway.Enabled {
//			pgwRemoteAddr = config.Prom.PushGateway.RemoteAddr
//			pgwIntervalMS = config.Prom.PushGateway.IntervalMS
//		}
//
//		res = rk_prom.NewPromEntry(
//			rk_prom.WithPort(config.Prom.Port),
//			rk_prom.WithPath(config.Prom.Path),
//			rk_prom.WithPGWRemoteAddr(pgwRemoteAddr),
//			rk_prom.WithPGWIntervalMS(pgwIntervalMS),
//		)
//	}
//
//	return res
//}

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

func mergeLogger(src []byte, event interface{}) *zap.Logger {
	// get event config as make it as bytes
	override, err := yaml.Marshal(event)
	if err != nil {
		shutdownWithError(err)
	}

	// merge logger config specified and element specified
	merged := mergeBytesAsMap(src, override)
	logger, _, err := rk_logger.NewZapLoggerWithBytes(merged, rk_logger.YAML)
	if err != nil {
		shutdownWithError(err)
	}

	return logger
}

// merge two maps unmarshalled with YAML type
func mergeBytesAsMap(src []byte, override []byte) []byte {
	var srcMap map[string]interface{}
	if err := yaml.Unmarshal(src, &srcMap); err != nil {
		shutdownWithError(err)
	}

	var overrideMap map[string]interface{}
	if err := yaml.Unmarshal(override, &overrideMap); err != nil {
		shutdownWithError(err)
	}

	for k, v := range overrideMap {
		if _, ok := srcMap[k]; ok {
			srcMap[k] = v
		}
	}

	res, err := yaml.Marshal(srcMap)
	if err != nil {
		shutdownWithError(err)
	}

	return res
}

func shutdownWithError(err error) {
	logger, _ := zap.NewDevelopment()
	logger.Error(fmt.Sprintf("%v", err))
	os.Exit(1)
}

func readFile(filePath string) []byte {
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

	return bytes
}
