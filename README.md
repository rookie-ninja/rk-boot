# rk-boot
gRpc service bootstrapper for goLang.
With rk-boot, users can start gRpc service with yaml formatted config file.
Easy to compile, run and debug your gRpc service, gRpc gateway and swagger UI.

- [rk-config](https://github.com/uber-go/zap)
- [rk-query](https://github.com/rookie-ninja/rk-query)
- [rk-logger](https://github.com/rookie-ninja/rk-logger)
- [rk-interceptor](https://github.com/rookie-ninja/rk-interceptor)
- [rk-prom](https://github.com/rookie-ninja/rk-prom)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Installation](#installation)
- [Quick Start](#quick-start)
  - [YAML config](#yaml-config)
    - [gRpc](#grpc)
  - [Development Status: Stable](#development-status-stable)
  - [Appendix](#appendix)
    - [Proto file compilation](#proto-file-compilation)
  - [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation
`go get -u rookie-ninja/rk-boot`

## Quick Start
There are two ways users can run gRpc or Gin service. one is yaml formatted config file.
The other one is through goLang code.

### YAML config
With human readable yaml config.
All you need to do is compile .proto file with protoc, protoc-gen-go, protoc-gen-grpc-gateway and protoc-gen-swagger

#### gRpc

Example:
```yaml
---
appName: rk-demo
event:
  format: RK
  quiet: false
logger:
  - name: app
    confPath: "example/configs/zap-app.yaml"
    forBoot: true
  - name: query
    confPath: "example/configs/zap-query.yaml"
    forEvent: true
config:
  - name: rk-main
    path: "example/configs/rk.yaml"
    format: RK
    global: false
grpc:
  - name: greeter
    port: 8080
    enableCommonService: true
    gw:
      enabled: true
      port: 8081
      insecure: true
      enableCommonService: true
      enableTV: true
    sw:
      enabled: true
      port: 8090
      path: sw
      jsonPath: "example/api/v1"
      insecure: true
      enableCommonService: true
    loggingInterceptor:
      enabled: true
      enableLogging: true
      enableMetrics: true
      enablePayloadLogging: false
prom:
  enabled: true
  port: 1608
  path: metrics
  pushGateway:
    enabled: false
    remoteAddr: xxx
    intervalMS: 2000
    jobName: xxx
proto:
  source:
    - api/v1/*.proto
  import:
    - third-party/googleapis
  doc:
    output: docs
    name: rk-server-demo
    type:
      - html
      - markdown
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
	"context"
	"github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-boot/example/api/v1"
	"google.golang.org/grpc"
	"time"
)

func main() {
	boot := rk_boot.NewBoot(rk_boot.WithBootConfigPath("example/configs/boot.yaml"))

	// register gRpc
	boot.GetGRpcEntry("greeter").AddRegFuncs(registerGreeter)
	boot.GetGRpcEntry("greeter").AddGWRegFuncs(hello_v1.RegisterGreeterHandlerFromEndpoint)

	boot.Bootstrap()
	boot.Quitter(5 * time.Second)
}

func registerGreeter(server *grpc.Server) {
	hello_v1.RegisterGreeterServer(server, &GreeterServer{})
}

type GreeterServer struct{}

func (server *GreeterServer) SayHello(ctx context.Context, request *hello_v1.HelloRequest) (*hello_v1.HelloResponse, error) {
	return &hello_v1.HelloResponse{
		Message: "hello",
	}, nil
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
| grpc.name | The name of gRpc server | string |
| grpc.port | The port of gRpc server | integer |
| grpc.enableCommonService | Enable embedded common service | true, false |
| grpc.gw.enabled | Enable gateway service over gRpc server | true, false |
| grpc.gw.port | The port of gRpc gateway | true, false |
| grpc.gw.enableTV | Enable RK TV | true, false |
| grpc.gw.insecure | Run gateway with insecure mode | true, false |
| grpc.gw.enableCommonService | Enable embedded common service | true, false |
| grpc.sw.enabled | Enable swagger service over gRpc server | true, false |
| grpc.sw.port | The port of swagger | true, false |
| grpc.sw.insecure | Run swagger with insecure mode | true, false |
| grpc.sw.enableCommonService | Enable embedded common service | true, false |
| grpc.sw.path | The path access swagger service from web | string |
| grpc.sw.jsonPath | Where the swagger.json files are stored locally | string |
| grpc.loggingInterceptor.enabled | Enable rk-interceptor logging interceptor | true, false |
| grpc.loggingInterceptor.enableLogging | Enable rk-interceptor logging interceptor specifically for each Rpc with rk-query | true, false |
| grpc.loggingInterceptor.enableMetrics | Enable rk-interceptor logging interceptor specifically for each Rpc with prometheus | true, false |
| grpc.loggingInterceptor.enablePayloadLogging | Enable rk-interceptor logging interceptor specifically for each Rpc's payload | true, false |
| gin.name | name of gin server entry| string | unknown application |
| gin.port | port of server | integer | nil, server won't start |
| gin.tls.enabled | enable tls or not | boolean | false | 
| gin.tls.user.enabled | enable user provided CA file? | boolean | false |
| gin.tls.user.certFile | cert file path | string | empty string |
| gin.tls.user.keyFile | key file path | string | empty string | 
| gin.tls.auth.enabled | server will generate CA files | string | false |
| gin.tls.auth.certOutput | cert file output path | string | current working directory | 
| gin.sw.enabled | enable swagger | boolean | false | 
| gin.sw.path | swagger path | string | / |
| gin.sw.jsonPath | swagger json file path | string | / |
| gin.sw.headers | headers will send with swagger response | array | empty array |
| gin.enableCommonService | enable common service | boolean | false |
| gin.enableTV | enable RK TV whose path is /v1/rk/tv | boolean | false |
| gin.loggingInterceptor.enabled | enable logging interceptor | boolean | false |
| gin.loggingInterceptor.enableLogging | enable logging for every request | boolean | false |
| gin.loggingInterceptor.enableMetrics | enable prometheus metrics for every request | boolean | false |
| gin.authInterceptor.enabled | enable auth interceptor | boolean | false |
| gin.authInterceptor.realm | realm for basic auth interceptor | string | Authorization Required |
| gin.authInterceptor.credentials | array of credentials such as "user:pass" | string array | empty array |
| prom.enabled | Enable local prometheus client | true, false |
| prom.port | The port of prometheus client | integer |
| prom.path | The path of prometheus client | string |
| prom.pusher.enabled | Enable pushGateway jobs locally | true, false |
| prom.pusher.url | pushGateway remote address | string |
| prom.pusher.interval | Push job intervals with seconds | integer |
| prom.pusher.job | pushGateway job name | string |

#### gin

Example:
```yaml
---
application: rk-server
event:
  format: RK
  quiet: false
  outputPaths:
    - stdout
  loggerConf: "example/gin/configs/query.yaml"
logger:
  - name: app
    quiet: false
    outputPath: stdout
    loggerConf: "example/gin/configs/app.yaml"
    maxsize: 1
    maxage: 7
    maxbackups: 3
    localtime: true
    compress: true
config:
  - name: rk-main
    path: "example/gin/configs/app.yaml"
    format: RK
    global: false
env:
  REALM: rk
  DOMAIN: dev
  REGION: cn-north-pek
  AZ: cn-north-pek-1
gin:
  - name: greeter
    port: 8080
    tls:
      enabled: false
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
      headers:
        - "cache-control: no-cache"
    enableCommonService: true
    enableTV: true
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
ut:
  output: docs
build:
  goos: darwin
  goarch: amd64

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
	boot.Wait(5 * time.Second)
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

Available configuration
User can start multiple servers at the same time

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.name | name of gin server entry| string | unknown application |
| gin.port | port of server | integer | nil, server won't start |
| gin.tls.enabled | enable tls or not | boolean | false | 
| gin.tls.user.enabled | enable user provided CA file? | boolean | false |
| gin.tls.user.certFile | cert file path | string | empty string |
| gin.tls.user.keyFile | key file path | string | empty string | 
| gin.tls.auth.enabled | server will generate CA files | string | false |
| gin.tls.auth.certOutput | cert file output path | string | current working directory | 
| gin.sw.enabled | enable swagger | boolean | false | 
| gin.sw.path | swagger path | string | / |
| gin.sw.jsonPath | swagger json file path | string | / |
| gin.sw.headers | headers will send with swagger response | array | empty array |
| gin.enableCommonService | enable common service | boolean | false |
| gin.enableTV | enable RK TV whose path is /v1/rk/tv | boolean | false |
| gin.loggingInterceptor.enabled | enable logging interceptor | boolean | false |
| gin.loggingInterceptor.enableLogging | enable logging for every request | boolean | false |
| gin.loggingInterceptor.enableMetrics | enable prometheus metrics for every request | boolean | false |
| gin.authInterceptor.enabled | enable auth interceptor | boolean | false |
| gin.authInterceptor.realm | realm for basic auth interceptor | string | Authorization Required |
| gin.authInterceptor.credentials | array of credentials such as "user:pass" | string array | empty array |

### Common Services
User can start multiple servers at the same time

| path | description |
| ------ | ------ |
| /v1/rk/healthy | always return true if service is available |
| /v1/rk/gc | trigger gc and return memory stats |
| /v1/rk/info | return basic info |
| /v1/rk/config | return configs in memory |
| /v1/rk/apis | list all apis |
| /v1/rk/sys | return system information including cpu and memory usage |
| /v1/rk/req | return requests stats recorded by prometheus client |
| /v1/rk/tv | web ui for metrics |

### Development Status: Stable

### Appendix
#### Proto file compilation
With rk command line tools.
Please refer [rk-cmd](https://github.com/rookie-ninja/rk-cmd)

With ptoroc command
Compile to go file
```shell script
protoc -I. -I third-party/googleapis --go_out=plugins=grpc:. --go_opt=paths=source_relative api/v1/*.proto
```

Compile to gw.go file
```shell script
protoc -I. -I third-party/googleapis --grpc-gateway_out=logtostderr=true,paths=source_relative:. api/v1/*.proto
```

Compile to gw.go and swagger.json file
```shell script
protoc -I. -I third-party/googleapis --grpc-gateway_out=logtostderr=true,paths=source_relative:. --swagger_out=logtostderr=true:. api/v1/*.proto
```

### Contributing
We encourage and support an active, healthy community of contributors &mdash;
including you! Details are in the [contribution guide](CONTRIBUTING.md) and
the [code of conduct](CODE_OF_CONDUCT.md). The rk maintainers keep an eye on
issues and pull requests, but you can also report any negative conduct to
dongxuny@gmail.com. That email list is a private, safe space; even the zap
maintainers don't have access, so don't hesitate to hold us to a high
standard.

<hr>

Released under the [MIT License](LICENSE).