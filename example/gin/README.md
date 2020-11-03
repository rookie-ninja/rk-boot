# gRPC Example
gRpc service bootstrapper for goLang.
With rk-boot, users can start gRpc service with yaml formatted config file.
Easy to compile, run and debug your gRpc service, gRpc gateway and swagger UI.

- [rk-config](https://github.com/uber-go/zap)
- [rk-query](https://github.com/rookie-ninja/rk-query)
- [rk-logger](https://github.com/rookie-ninja/rk-logger)
- [rk-interceptor](https://github.com/rookie-ninja/rk-interceptor)
- [rk-prom](https://github.com/rookie-ninja/rk-prom)

## Installation
`go get -u rookie-ninja/rk-boot`

## Quick Start
There are two ways users can run Gin service. one is yaml formatted config file.
The other one is through goLang code.

### Example
With human readable yaml config.
All you need to do is compile swagger json file with swag

Directory struct
```shell script
.
├── configs
│   ├── boot.yaml
│   ├── rk.yaml
│   ├── zap-app.yaml
│   └── zap-query.yaml
└── server
    └── gin.go
```

boot.yaml:
```yaml
---
appName: rk-server
event:
  format: RK
  quiet: false
logger:
  - name: app
    confPath: "example/gin/configs/zap-app.yaml"
    forBoot: true
  - name: query
    confPath: "example/gin/configs/zap-query.yaml"
    forEvent: true
config:
  - name: rk-main
    path: "example/gin/configs/rk.yaml"
    format: RK
    global: false
gin:
  - name: greeter
    port: 8080
    tls:
      enabled: true
      port: 8443
      user:
        enabled: false
        certFile: "example/gin/server/cert/server.pem"
        keyFile: "example/gin/server/cert/server-key.pem"
      auto:
        enabled: true
        certOutput: "example/gin/server/cert"
    sw:
      enabled: true
      path: "sw"
      jsonPath: "example/gin/server/docs"
      insecure: true
      headers:
        - "cache-control: no-cache"
    enableCommonService: true
    loggingInterceptor:
      enabled: true
      enableLogging: true
      enableMetrics: true
    authInterceptor:
      enabled: false
      realm: "rk"
      credentials:
        - "foo:pass"
        - "bar:pass"
prom:
  enabled: true
  port: 1608
  path: metrics
  pushGateway:
    enabled: false
    remoteAddr: xxx
    intervalMS: 2000
    jobName: xxx
ut:
  output: docs
```

```go
// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rookie-ninja/rk-boot"
	"net/http"
	"time"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample rk-demo server.
// @termsOfService http://swagger.io/terms/

// @securityDefinitions.basic BasicAuth

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	boot := rk_boot.NewBoot(rk_boot.WithBootConfigPath("example/gin/configs/boot.yaml"))

	boot.GetGinEntry("greeter").GetRouter().GET("/v1/hello", hello)

	boot.Bootstrap()
	boot.Quitter(5 * time.Second)
}

// @Summary Hello
// @Id 1
// @Tags Hello
// @version 1.0
// @produce application/json
// @Success 200 string string
// @Router /v1/hello [get]
func hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "hello!",
	})
}
```

YAML config Explanation

| Name | Description | Option |
| ------ | ------ | ------ |
| appName | Your application name, would be logged into logger | string |
| event.format | The format of event data, please refer [rk-query](https://github.com/rookie-ninja/rk-query) for details | RK, JSON |
| event.quiet | Set event to quiet mode, please refer [rk-query](https://github.com/rookie-ninja/rk-query) for details | true, false |
| logger.name | Zap logger name | string |
| logger.confPath | Zap logger config path, if relative path was given, then os.GetWd() would be used for abs path | string |
| logger.forBoot | Use this logger for log while bootstrapping | true, false |
| logger.forEvent | Use this logger for rk-query | true, false |
| config.name | Viper config name | string |
| config.path | Viper config file path | string |
| config.format | Config format, viper support standard yaml, hcl, toml, json. For RK format please refer [rk-config](https://github.com/uber-go/zap) | RK, Viper |
| config.global | Whether access config globally, like rk_config.Get() | true, false |
| gin.name | The name of gin server | string |
| gin.port | The port of gin server | number |
| gin.enableCommonService | Enable embedded common service | true, false |
| gin.tls.enabled | Enable tls | true, false |
| gin.tls.port | The port of tls server | number |
| gin.tls.user.enabled | With user CA files | true, false |
| gin.tls.user.certFile | Cert file path | string |
| gin.tls.user.keyFile | Key file path | string |
| gin.tls.auto.enabled | Generate tls CA files by default | true, false |
| gin.tls.auto.certOutput | Output path for auto generated CA files | string |
| gin.sw.enabled | Enable swagger service over gRpc server | true, false |
| gin.sw.port | The port of swagger | true, false |
| gin.sw.insecure | Run swagger with insecure mode | true, false |
| gin.sw.enableCommonService | Enable embedded common service | true, false |
| gin.sw.path | The path access swagger service from web | string |
| gin.sw.jsonPath | Where the swagger.json files are stored locally | string |
| gin.sw.headers | Default headers append to swagger | string array |
| gin.loggingInterceptor.enabled | Enable rk-interceptor logging interceptor | true, false |
| gin.loggingInterceptor.enableLogging | Enable rk-interceptor logging interceptor specifically for each Rpc with rk-query | true, false |
| gin.loggingInterceptor.enableMetrics | Enable rk-interceptor logging interceptor specifically for each Rpc with prometheus | true, false |
| gin.authInterceptor.enabled | Enable rk-interceptor basic auth interceptor | true, false |
| gin.authInterceptor.realm | Enable rk-interceptor basic auth realm | string |
| gin.authInterceptor.credentials | Credentials for basic auth | string array |
| prom.enabled | Enable local prometheus client | true, false |
| prom.port | The port of prometheus client | number |
| prom.path | The path of prometheus client | string |
| prom.pushGateway.enabled | Enable pushGateway jobs locally | true, false |
| prom.pushGateway.remoteAddr | pushGateway remote address | string |
| prom.pushGateway.intervalMS | Push job intervals with milliseconds | number |
| prom.pushGateway.jobName | pushGateway job name | string |

### Development Status: Stable

### Appendix
#### swagger json compilation
With ptoroc command
Compile to go file, swagger file would be generated under docs/
```shell script
swag init -g gin.go
```

Default CA files would be generated if tls enalbed