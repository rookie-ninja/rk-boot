// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_sw

import (
	"context"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rookie-ninja/rk-boot/api/v1"
	"github.com/rookie-ninja/rk-boot/gw"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	swHandlerPrefix = "/swagger/"
	gwHandlerPrefix = "/"
	swAssetsPath    = "./assets/swagger-ui/"
)

var (
	swaggerIndexHTML = `<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>PL Swagger</title>
    <link rel="stylesheet" type="text/css" href="http://pulseline-prod-cdn-cn-south-can-1258344699.file.myqcloud.com/swagger-ui/3.24.2/swagger-ui.css" >
    <link rel="icon" type="image/png" href="http://pulseline-prod-cdn-cn-south-can-1258344699.file.myqcloud.com/swagger-ui/3.24.2/favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="http://pulseline-prod-cdn-cn-south-can-1258344699.file.myqcloud.com/swagger-ui/3.24.2/favicon-16x16.png" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }

      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }

      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>

    <script src="http://pulseline-prod-cdn-cn-south-can-1258344699.file.myqcloud.com/swagger-ui/3.24.2/swagger-ui-bundle.js"> </script>
    <script src="http://pulseline-prod-cdn-cn-south-can-1258344699.file.myqcloud.com/swagger-ui/3.24.2/swagger-ui-standalone-preset.js"> </script>
    <script>
    window.onload = function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
          configUrl: "swagger-config.json",
          dom_id: '#swagger-ui',
          deepLinking: true,
          presets: [
              SwaggerUIBundle.presets.apis,
              SwaggerUIStandalonePreset
          ],
          plugins: [
              SwaggerUIBundle.plugins.DownloadUrl
          ],
          layout: "StandaloneLayout"
      })
      // End Swagger UI call region

      window.ui = ui
    }
  </script>
  </body>
</html>
`
	commonServiceJson = `
{
  "swagger": "2.0",
  "info": {
    "title": "api/rk_common_service.proto",
    "version": "version not set"
  },
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/v1/rk/config": {
      "get": {
        "summary": "DumpConfig Stub",
        "operationId": "RkCommonService_DumpConfig",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/DumpConfigResponse"
            }
          },
          "default": {
            "description": "An unexpected error response",
            "schema": {
              "$ref": "#/definitions/runtimeError"
            }
          }
        },
        "tags": [
          "RkCommonService"
        ]
      }
    },
    "/v1/rk/config/{key}": {
      "get": {
        "summary": "GetConfig Stub",
        "operationId": "RkCommonService_GetConfig",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/GetConfigResponse"
            }
          },
          "default": {
            "description": "An unexpected error response",
            "schema": {
              "$ref": "#/definitions/runtimeError"
            }
          }
        },
        "parameters": [
          {
            "name": "key",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "tags": [
          "RkCommonService"
        ]
      }
    },
    "/v1/rk/gc": {
      "post": {
        "summary": "GC Stub",
        "operationId": "RkCommonService_GC",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/GCResponse"
            }
          },
          "default": {
            "description": "An unexpected error response",
            "schema": {
              "$ref": "#/definitions/runtimeError"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/GCRequest"
            }
          }
        ],
        "tags": [
          "RkCommonService"
        ]
      }
    },
    "/v1/rk/healthy": {
      "get": {
        "summary": "Healthy Stub",
        "operationId": "RkCommonService_Healthy",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/HealthyResponse"
            }
          },
          "default": {
            "description": "An unexpected error response",
            "schema": {
              "$ref": "#/definitions/runtimeError"
            }
          }
        },
        "tags": [
          "RkCommonService"
        ]
      }
    },
    "/v1/rk/help": {
      "get": {
        "summary": "Help Stub",
        "operationId": "RkCommonService_Help",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/HelpResponse"
            }
          },
          "default": {
            "description": "An unexpected error response",
            "schema": {
              "$ref": "#/definitions/runtimeError"
            }
          }
        },
        "tags": [
          "RkCommonService"
        ]
      }
    },
    "/v1/rk/info": {
      "get": {
        "summary": "Info Stub",
        "operationId": "RkCommonService_Info",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/InfoResponse"
            }
          },
          "default": {
            "description": "An unexpected error response",
            "schema": {
              "$ref": "#/definitions/runtimeError"
            }
          }
        },
        "tags": [
          "RkCommonService"
        ]
      }
    },
    "/v1/rk/log": {
      "post": {
        "summary": "Log Stub",
        "operationId": "RkCommonService_Log",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/LogResponse"
            }
          },
          "default": {
            "description": "An unexpected error response",
            "schema": {
              "$ref": "#/definitions/runtimeError"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/LogRequest"
            }
          }
        ],
        "tags": [
          "RkCommonService"
        ]
      }
    },
    "/v1/rk/ping": {
      "post": {
        "summary": "Ping Stub",
        "operationId": "RkCommonService_Ping",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/PongResponse"
            }
          },
          "default": {
            "description": "An unexpected error response",
            "schema": {
              "$ref": "#/definitions/runtimeError"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/PingRequest"
            }
          }
        ],
        "tags": [
          "RkCommonService"
        ]
      }
    },
    "/v1/rk/shutdown": {
      "post": {
        "summary": "Shutdown Stub",
        "operationId": "RkCommonService_Shutdown",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/ShutdownResponse"
            }
          },
          "default": {
            "description": "An unexpected error response",
            "schema": {
              "$ref": "#/definitions/runtimeError"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/ShutdownRequest"
            }
          }
        ],
        "tags": [
          "RkCommonService"
        ]
      }
    }
  },
  "definitions": {
    "BasicInfo": {
      "type": "object",
      "properties": {
        "start_time": {
          "type": "string"
        },
        "up_time": {
          "type": "string"
        },
        "realm": {
          "type": "string"
        },
        "region": {
          "type": "string"
        },
        "az": {
          "type": "string"
        },
        "domain": {
          "type": "string"
        },
        "app_name": {
          "type": "string"
        }
      }
    },
    "Config": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "value": {
          "type": "string"
        }
      }
    },
    "DumpConfigResponse": {
      "type": "object",
      "properties": {
        "config": {
          "type": "string"
        }
      }
    },
    "GCRequest": {
      "type": "object",
      "properties": {
        "operator": {
          "type": "string"
        }
      },
      "title": "GC request, operator must be set"
    },
    "GCResponse": {
      "type": "object",
      "properties": {
        "mem_stats_before": {
          "$ref": "#/definitions/MemStats"
        },
        "mem_stats_after": {
          "$ref": "#/definitions/MemStats"
        }
      },
      "title": "GC response, memory stats would be returned"
    },
    "GRpcInfo": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "port": {
          "type": "string"
        },
        "gw_info": {
          "$ref": "#/definitions/GWInfo"
        },
        "sw_info": {
          "$ref": "#/definitions/SWInfo"
        }
      }
    },
    "GWInfo": {
      "type": "object",
      "properties": {
        "gw_port": {
          "type": "string"
        }
      }
    },
    "GetConfigResponse": {
      "type": "object",
      "properties": {
        "config_list": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Config"
          }
        }
      }
    },
    "HealthyResponse": {
      "type": "object",
      "properties": {
        "healthy": {
          "type": "boolean",
          "format": "boolean"
        }
      }
    },
    "HelpResponse": {
      "type": "object",
      "properties": {
        "stubs": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Stub"
          }
        }
      }
    },
    "InfoResponse": {
      "type": "object",
      "properties": {
        "basic_info": {
          "$ref": "#/definitions/BasicInfo"
        },
        "grpc_info_list": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/GRpcInfo"
          }
        },
        "prom_info": {
          "$ref": "#/definitions/PromInfo"
        }
      }
    },
    "LogEntry": {
      "type": "object",
      "properties": {
        "log_name": {
          "type": "string"
        },
        "log_level": {
          "type": "string"
        }
      }
    },
    "LogRequest": {
      "type": "object",
      "properties": {
        "entries": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/LogEntry"
          }
        }
      }
    },
    "LogResponse": {
      "type": "object"
    },
    "MemStats": {
      "type": "object",
      "properties": {
        "mem_alloc_mb": {
          "type": "string",
          "format": "uint64",
          "description": "Alloc is bytes of allocated heap objects."
        },
        "sys_mem_mb": {
          "type": "string",
          "format": "uint64",
          "description": "Sys is the total bytes of memory obtained from the OS."
        },
        "last_gc_timestamp": {
          "type": "string",
          "title": "LastGC is the time the last garbage collection finished.\nRepresent as RFC3339 time format"
        },
        "num_gc": {
          "type": "integer",
          "format": "int64",
          "description": "NumGC is the number of completed GC cycles."
        },
        "num_force_gc": {
          "type": "integer",
          "format": "int64",
          "description": "/ NumForcedGC is the number of GC cycles that were forced by\nthe application calling the GC function."
        }
      },
      "title": "Memory stats"
    },
    "PingRequest": {
      "type": "object"
    },
    "PongResponse": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string"
        }
      }
    },
    "PromInfo": {
      "type": "object",
      "properties": {
        "port": {
          "type": "string"
        },
        "path": {
          "type": "string"
        }
      }
    },
    "SWInfo": {
      "type": "object",
      "properties": {
        "sw_port": {
          "type": "string"
        },
        "sw_path": {
          "type": "string"
        }
      }
    },
    "ShutdownRequest": {
      "type": "object"
    },
    "ShutdownResponse": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string"
        }
      }
    },
    "Stub": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "description": {
          "type": "string"
        },
        "usage": {
          "type": "string"
        }
      }
    },
    "protobufAny": {
      "type": "object",
      "properties": {
        "type_url": {
          "type": "string"
        },
        "value": {
          "type": "string",
          "format": "byte"
        }
      }
    },
    "runtimeError": {
      "type": "object",
      "properties": {
        "error": {
          "type": "string"
        },
        "code": {
          "type": "integer",
          "format": "int32"
        },
        "message": {
          "type": "string"
        },
        "details": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/protobufAny"
          }
        }
      }
    }
  }
}

`
)

type swURLConfig struct {
	URLs []*swURL `json:"urls"`
}

type swURL struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type SWEntry struct {
	logger              *zap.Logger
	swPort              uint64
	gRpcPort            uint64
	jsonPath            string
	path                string
	enableCommonService bool
	fileHandler         http.Handler
	regFuncs            []rk_gw.RegFunc
	dialOpts            []grpc.DialOption
	muxOpts             []runtime.ServeMuxOption
	server              *http.Server
}

type SWOption func(*SWEntry)

func WithSWPort(port uint64) SWOption {
	return func(entry *SWEntry) {
		entry.swPort = port
	}
}

func WithCommonService(enable bool) SWOption {
	return func(entry *SWEntry) {
		entry.enableCommonService = enable
	}
}

func WithGRpcPort(port uint64) SWOption {
	return func(entry *SWEntry) {
		entry.gRpcPort = port
	}
}

func WithPath(path string) SWOption {
	return func(entry *SWEntry) {
		entry.path = path
	}
}

func WithJsonPath(path string) SWOption {
	return func(entry *SWEntry) {
		entry.jsonPath = path
	}
}

func WithRegFuncs(funcs ...rk_gw.RegFunc) SWOption {
	return func(entry *SWEntry) {
		entry.regFuncs = append(entry.regFuncs, funcs...)
	}
}

func WithDialOptions(opts ...grpc.DialOption) SWOption {
	return func(entry *SWEntry) {
		entry.dialOpts = append(entry.dialOpts, opts...)
	}
}

func NewSWEntry(opts ...SWOption) *SWEntry {
	entry := &SWEntry{
		logger: zap.NewNop(),
	}

	for i := range opts {
		opts[i](entry)
	}

	if entry.dialOpts == nil {
		entry.dialOpts = make([]grpc.DialOption, 0)
	}

	if entry.regFuncs == nil {
		entry.regFuncs = make([]rk_gw.RegFunc, 0)
	}

	entry.fileHandler = func() http.Handler {
		return http.FileServer(http.Dir("./assets/swagger-ui"))
	}()

	if entry.enableCommonService {
		entry.regFuncs = append(entry.regFuncs, rk_boot_common_v1.RegisterRkCommonServiceHandlerFromEndpoint)
	}

	return entry
}

func (entry *SWEntry) AddDialOptions(opts ...grpc.DialOption) {
	entry.dialOpts = append(entry.dialOpts, opts...)
}

func (entry *SWEntry) AddRegFuncs(funcs ...rk_gw.RegFunc) {
	entry.regFuncs = append(entry.regFuncs, funcs...)
}

func (entry *SWEntry) GetSWPort() uint64 {
	return entry.swPort
}

func (entry *SWEntry) GetGRpcPort() uint64 {
	return entry.gRpcPort
}

func (entry *SWEntry) GetPath() string {
	return entry.path
}

func (entry *SWEntry) GetServer() *http.Server {
	return entry.server
}

func (entry *SWEntry) Stop(logger *zap.Logger) {
	if logger == nil {
		logger = zap.NewNop()
	}

	if entry.server != nil {
		logger.Info("stopping swagger",
			zap.Uint64("sw_port", entry.swPort),
			zap.Uint64("gRpc_port", entry.gRpcPort),
			zap.String("sw_path", entry.path))
		entry.server.Shutdown(context.Background())
	}
}

func (entry *SWEntry) Start(logger *zap.Logger) {
	if logger == nil {
		logger = zap.NewNop()
	}

	// Deal with Path
	// add "/" at start and end side if missing
	if !strings.HasPrefix(entry.path, "/") {
		entry.path = "/" + entry.path
	}

	if !strings.HasSuffix(entry.path, "/") {
		entry.path = entry.path + "/"
	}

	// 1: create ./assets/swagger-ui if missing
	entry.createSWAssetsPath()

	// 2: create ./assets/swagger-ui/index.html if missing
	entry.createIndexHtml()

	// 3: create or modify ./assets/swagger-ui/swagger-config.json
	entry.createOrModifySWURLConfig()

	// Init http server
	ctx := context.Background()

	gRPCEndpoint := "0.0.0.0:" + strconv.FormatUint(entry.gRpcPort, 10)
	swEndpoint := "0.0.0.0:" + strconv.FormatUint(entry.swPort, 10)

	gwMux := runtime.NewServeMux()

	for i := range entry.regFuncs {
		err := entry.regFuncs[i](ctx, gwMux, gRPCEndpoint, entry.dialOpts)
		if err != nil {
			entry.logger.Error("failed to register gateway function",
				zap.Uint64("gRpc_port", entry.gRpcPort),
				zap.Uint64("sw_port", entry.swPort),
				zap.String("sw_path", entry.path))
			shutdownWithError(err)
		}
	}

	// Init swagger http mux
	httpMux := http.NewServeMux()
	httpMux.Handle(gwHandlerPrefix, gwMux)
	httpMux.HandleFunc(swHandlerPrefix, entry.swJsonFileHandler)
	httpMux.Handle(entry.path, http.StripPrefix(entry.path, entry.fileHandler))

	entry.server = &http.Server{
		Addr:    swEndpoint,
		Handler: httpMux,
	}

	go func(entry *SWEntry) {
		logger.Info("starting swagger",
			zap.Uint64("sw_port", entry.swPort),
			zap.Uint64("gRpc_port", entry.gRpcPort),
			zap.String("sw_path", entry.path))
		if err := entry.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			entry.logger.Error("failed to start swagger",
				zap.Uint64("gRpc_port", entry.gRpcPort),
				zap.Uint64("sw_port", entry.swPort),
				zap.String("sw_path", entry.path),
				zap.Error(err))
			shutdownWithError(err)
		}
	}(entry)
}

func (entry *SWEntry) createSWAssetsPath() {
	err := os.MkdirAll(swAssetsPath, os.ModePerm)
	if err != nil {
		entry.logger.Error("failed to create folder to swagger",
			zap.String("sw_assets_path", swAssetsPath),
			zap.Error(err))
		shutdownWithError(err)
	}
}

func (entry *SWEntry) createIndexHtml() {
	_, err := os.Stat(path.Join(swAssetsPath, "index.html"))
	if os.IsNotExist(err) {
		// create a default index.html file
		err := ioutil.WriteFile(path.Join(swAssetsPath, "index.html"), []byte(swaggerIndexHTML), os.ModePerm)
		if err != nil {
			entry.logger.Error("failed to create index.html for swagger",
				zap.String("sw_assets_path", swAssetsPath),
				zap.Error(err))
			shutdownWithError(err)
		}
	} else if err != nil {
		entry.logger.Error("failed to stat index.html for swagger",
			zap.String("sw_assets_path", swAssetsPath),
			zap.Error(err))
		shutdownWithError(err)
	}
}

func (entry *SWEntry) createOrModifySWURLConfig() {
	// 1: Get swagger-config.json if exists
	swaggerURLConfig := entry.getSWURLConfigByName()

	// 2: Add user API swagger JSON
	swaggerJSONFiles := entry.listFilesWithSuffix()
	for i := range swaggerJSONFiles {
		swaggerURL := &swURL{
			Name: strings.TrimSuffix(swaggerJSONFiles[i], ".json"),
			URL:  path.Join("/swagger", swaggerJSONFiles[i]),
		}
		entry.appendAndDeduplication(swaggerURLConfig, swaggerURL)
	}

	// 3: Add pl-common
	entry.appendAndDeduplication(swaggerURLConfig, &swURL{
		Name: "rk-common",
		URL:  "/swagger/rk_common_service.swagger.json",
	})

	// 4: Marshal to swagger-config.json
	bytes, err := json.Marshal(swaggerURLConfig)
	if err != nil {
		entry.logger.Warn("failed to unmarshal swagger-config.json",
			zap.Uint64("gRpc_port", entry.gRpcPort),
			zap.Uint64("sw_port", entry.swPort),
			zap.String("sw_assets_path", swAssetsPath),
			zap.Error(err))
		shutdownWithError(err)
	}

	// 5: Create swagger-config.json
	err = ioutil.WriteFile(path.Join(swAssetsPath, "swagger-config.json"), bytes, os.ModePerm)
	if err != nil {
		entry.logger.Warn("failed to create swagger-config.json",
			zap.Uint64("gRpc_port", entry.gRpcPort),
			zap.Uint64("sw_port", entry.swPort),
			zap.String("sw_assets_path", swAssetsPath),
			zap.Error(err))
		shutdownWithError(err)
	}
}

func (entry *SWEntry) getSWURLConfigByName() *swURLConfig {
	config := &swURLConfig{
		URLs: make([]*swURL, 0),
	}

	_, err := os.Stat(path.Join(swAssetsPath, "swagger-config.json"))

	// swagger-config.json exists, read value from it
	if err == nil {
		// Exist! override it
		content, err := ioutil.ReadFile(path.Join(swAssetsPath, "swagger-config.json"))
		if err != nil {
			entry.logger.Warn("failed to read swagger-config.json",
				zap.Uint64("gRpc_port", entry.gRpcPort),
				zap.Uint64("sw_port", entry.swPort),
				zap.String("sw_assets_path", swAssetsPath),
				zap.Error(err))
			return config
		}

		err = json.Unmarshal(content, config)
		if err != nil {
			entry.logger.Warn("failed to unmarshal swagger-config.json",
				zap.Uint64("gRpc_port", entry.gRpcPort),
				zap.Uint64("sw_port", entry.swPort),
				zap.String("sw_assets_path", swAssetsPath),
				zap.Error(err))
			return config
		}
	}

	return config
}

func (entry *SWEntry) listFilesWithSuffix() []string {
	res := make([]string, 0)

	jsonPath := entry.jsonPath
	suffix := ".json"
	// re-path it with working directory if not absolute path
	if !path.IsAbs(entry.jsonPath) {
		wd, err := os.Getwd()
		if err != nil {
			entry.logger.Info("failed to get working directory",
				zap.String("error", err.Error()))
			return res
		}
		jsonPath = path.Join(wd, jsonPath)
	}

	files, err := ioutil.ReadDir(jsonPath)
	if err != nil {
		entry.logger.Info("failed to list files with suffix",
			zap.String("path", jsonPath),
			zap.String("suffix", suffix),
			zap.String("error", err.Error()))
		return res
	}

	for i := range files {
		file := files[i]
		if !file.IsDir() && strings.HasSuffix(file.Name(), suffix) {
			res = append(res, file.Name())
		}
	}

	return res
}

func (entry *SWEntry) appendAndDeduplication(config *swURLConfig, url *swURL) {
	urls := config.URLs

	for i := range urls {
		element := urls[i]

		if element.Name == url.Name {
			return
		}
	}

	config.URLs = append(config.URLs, url)
}

func (entry *SWEntry) swJsonFileHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "swagger.json") {
		http.NotFound(w, r)
		return
	}

	p := strings.TrimPrefix(r.URL.Path, swHandlerPrefix)
	// This is common file
	if p == "rk_common_service.swagger.json" {
		http.ServeContent(w, r, "rk-common", time.Now(), strings.NewReader(commonServiceJson))
		return
	}

	p = path.Join(entry.jsonPath, p)

	http.ServeFile(w, r, p)
}

func shutdownWithError(err error) {
	glog.Error(err)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}
