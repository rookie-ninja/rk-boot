# rk-boot
Bootstrapper for rkentry.Entry.
With rk-boot, users can start grpc, gin, prometheus client or custom entry service with yaml formatted config file.
Easy to compile, run and debug your grpc service, grpc gateway, swagger UI and rk-tv web UI.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Installation](#installation)
- [Quick Start](#quick-start)
  - [YAML config](#yaml-config)
    - [grpc](#grpc)
    - [gin](#gin)
  - [Common Service](#common-service)
  - [TV Service](#tv-service)
  - [Development Status: Stable](#development-status-stable)
  - [Appendix](#appendix)
    - [Proto file compilation](#proto-file-compilation)
  - [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation
`go get -u github.com/rookie-ninja/rk-boot`

## Quick Start
There are two ways users can run gRpc or Gin service. one is yaml formatted config file.
The other one is through golang code.

### YAML config
With human readable yaml config.
All you need to do is compile .proto file with buf.

There is example in Makefile.

#### grpc

Example:
```yaml
---
# Bellow sections are RK common entries which is not necessary.
rk:
  appName: myapp
  version: v0.0.1                                      # Optional, default: "v0.0.0"
  description: "this is description"                   # Optional, default: ""
  keywords: [ "rk", "golang" ]                         # Optional, default: []
  homeURL: "http://example.com"                        # Optional, default: ""
  iconURL: "http://example.com"                        # Optional, default: ""
  docsURL: [ "http://example.com" ]                    # Optional, default: []
  maintainers: [ "rk-dev" ]                            # Optional, default: []
zapLogger:                                             # Optional
  - name: zap-logger                                   # Required
    description: "Description of entry"                # Optional
eventLogger:                                           # Optional
  - name: event-logger                                 # Required
    description: "Description of entry"                # Optional
cert:                                                  # Optional
  - name: "local-cert"                                 # Required
    description: "Description of entry"                # Optional
    provider: "localFs"                                # Required, etcd, consul, localFs, remoteFs are supported options
    locale: "*::*::*::*"                               # Optional, default: *::*::*::*
    serverCertPath: "example/grpc/certs/server.pem"    # Optional, default: "", path of certificate on local FS
    serverKeyPath: "example/grpc/certs/server-key.pem" # Optional, default: "", path of certificate on local FS
    clientCertPath: "example/grpc/certs/server.pem"    # Optional, default: "", path of certificate on local FS
config:
  - name: rk-main                                      # Required
    path: "example/grpc/config/config.yaml"            # Required
    description: "Description of entry"                # Optional
    locale: "*::*::*::*"                               # Optional, default: *::*::*::*
# Grpc entry
grpc:
  - name: greeter                                      # Required
    port: 1949                                         # Optional
    commonService:
      enabled: true                                    # Optional, default: false
    cert:
      ref: "local-cert"                                # Optional, default: "", reference of cert entry declared above
    gw:
      enabled: true                                    # Optional, default: false
      port: 8080                                       # Required
      gwMappingFilePaths:
        - "example/grpc/api/v1/gw_mapping.yaml"
      cert:
        ref: "local-cert"                              # Optional, default: "", reference of cert entry declared above
      tv:
        enabled: true                                  # Optional, default: false
      sw:
        enabled: true                                  # Optional, default: false
        path: "sw"                                     # Optional, default: "sw"
        jsonPath: "example/grpc/api/gen/v1"            # Optional
        headers: [ "sw:rk" ]                           # Optional, default: []
      prom:
        enabled: true                                  # Optional, default: false
        path: "metrics"                                # Optional, default: ""
        pusher:
          enabled: false                               # Optional, default: false
          jobName: "greeter-pusher"                    # Required
          remoteAddress: "localhost:9091"              # Required
          basicAuth: "user:pass"                       # Optional, default: ""
          intervalMS: 1000                             # Optional, default: 1000
          cert:
            ref: "local-cert"                          # Optional, default: "", reference of cert entry declared above
    logger:                                            # Optional
      zapLogger:                                       # Optional
        ref: zap-logger                                # Optional, default: logger of STDOUT, reference of logger entry declared above
      eventLogger:                                     # Optional
        ref: event-logger                              # Optional, default: logger of STDOUT, reference of logger entry declared above
    interceptors:
      loggingZap:
        enabled: true                                  # Optional, default: false
      metricsProm:
        enabled: true                                  # Optional, default: false
      basicAuth:
        enabled: false                                 # Optional, default: false
        credentials:
          - "user:pass"                                # Optional, default: ""
      tokenAuth:
        enabled: false
        tokens:
          - token: ""
            expired: false                             # Optional, default: ""
```

```go
package main

import (
	"context"
	"github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-boot/example/grpc/api/gen/v1"
	"google.golang.org/grpc"
)

func main() {
	boot := rkboot.NewBoot(rkboot.WithBootConfigPath("example/grpc/boot.yaml"))

	// register gRpc
	boot.GetGrpcEntry("greeter").AddGrpcRegFuncs(registerGreeter)
	boot.GetGrpcEntry("greeter").AddGwRegFuncs(hello.RegisterGreeterHandlerFromEndpoint)

	// Bootstrap
	boot.Bootstrap(context.TODO())
}

func registerGreeter(server *grpc.Server) {
	hello.RegisterGreeterServer(server, &GreeterServer{})
}

type GreeterServer struct{}

func (server *GreeterServer) SayHello(ctx context.Context, request *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{
		Message: "Hello " + request.Name,
	}, nil
}

```

YAML config Explanation

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.name | The name of gRpc server | string | N/A |
| grpc.port | The port of gRpc server | integer | nil, server won't start |
| grpc.commonService.enabled | Enable embedded common service | boolean | false |
| grpc.cert.ref | Reference of cert entry declared in cert section | string | "" |
| grpc.gw.enabled | Enable gateway service over gRpc server | boolean | false |
| grpc.gw.port | The port of gRpc gateway | integer | nil, server won't start |
| grpc.gw.gwMappingFilePaths | The grpc gateway mapping file path | string array | empty array |
| grpc.gw.cert.ref | Reference of cert entry declared in cert section | string | "" |
| grpc.gw.tv.enabled | Enable RK TV | boolean | false |
| grpc.gw.sw.enabled | Enable swagger service over gRpc server | boolean | false |
| grpc.sw.path | The path access swagger service from web | string | sw |
| grpc.sw.jsonPath | Where the swagger.json files are stored locally | string | "" |
| grpc.sw.headers | Headers would be sent to caller | map<string, string> | nil |
| grpc.prom.enabled | Enable prometheus | boolean | false |
| grpc.prom.path | Path of prometheus | string | metrics |
| grpc.prom.pusher.enabled | Enable prometheus pusher | bool | false |
| grpc.prom.pusher.jobName | Job name would be attached as label while pushing to remote pushgateway | string | "" |
| grpc.prom.pusher.remoteAddress | PushGateWay address, could be form of http://x.x.x.x or x.x.x.x | string | "" |
| grpc.prom.pusher.intervalMs | Push interval in milliseconds | string | 1000 |
| grpc.prom.pusher.basicAuth | Basic auth used to interact with remote pushgateway, form of \<user:pass\> | string | "" |
| grpc.prom.pusher.cert.ref | Reference of rkentry.CertEntry | string | "" |
| grpc.logger.zapLogger.ref | Reference of logger entry declared above | string | "" |
| grpc.logger.eventLogger.ref | Reference of logger entry declared above | string | "" |
| grpc.interceptors.loggingZap.enabled | Enable logging interceptor | boolean | false |
| grpc.interceptors.metricsProm.enabled | Enable prometheus metrics for every request | boolean | false |
| grpc.interceptors.basicAuth.enabled | Enable auth interceptor | boolean | false |
| grpc.interceptors.basicAuth.credentials | Provide basic auth credentials, form of \<user:pass\> | string | false |

#### gin

Example:
```yaml
---
# Bellow sections are RK common entries which is not necessary.
rk: # NOT required
  appName: rk-example-entry                           # Optional, default: "rkApp"
  version: v0.0.1                                     # Optional, default: "v0.0.0"
  description: "this is description"                  # Optional, default: ""
  keywords: ["rk", "golang"]                          # Optional, default: []
  homeUrl: "http://example.com"                       # Optional, default: ""
  iconUrl: "http://example.com"                       # Optional, default: ""
  docsUrl: ["http://example.com"]                     # Optional, default: []
  maintainers: ["rk-dev"]                             # Optional, default: []
zapLogger:                                            # Optional
  - name: zap-logger                                  # Required
    description: "Description of entry"               # Optional
eventLogger:                                          # Optional
  - name: event-logger                                # Required
    description: "Description of entry"               # Optional
cert:                                                 # Optional
  - name: "local-cert"                                # Required
    description: "Description of entry"               # Optional
    provider: "localFs"                               # Required, etcd, consul, localFs, remoteFs are supported options
    locale: "*::*::*::*"                              # Optional, default: *::*::*::*
    serverCertPath: "example/gin/certs/server.pem"    # Optional, default: "", path of certificate on local FS
    serverKeyPath: "example/gin/certs/server-key.pem" # Optional, default: "", path of certificate on local FS
#    clientCertPath: "example/client.pem"             # Optional, default: "", path of certificate on local FS
#    clientKeyPath: "example/client.pem"              # Optional, default: "", path of certificate on local FS
config:
  - name: rk-main                                     # Required
    path: "example/gin/config/config.yaml"            # Required
    description: "Description of entry"               # Optional
    locale: "*::*::*::*"                              # Optional, default: *::*::*::*
# Gin entry
gin:
  - name: greeter                                     # Required
    port: 8080                                        # Required
    cert:                                             # Optional
      ref: "local-cert"                               # Optional, default: "", reference of cert entry declared above
    sw:                                               # Optional
      enabled: true                                   # Optional, default: false
      jsonPath: "example/gin/docs"
      path: "sw"                                      # Optional, default: "sw"
      headers: ["sw:rk"]                              # Optional, default: []
    commonService:                                    # Optional
      enabled: true                                   # Optional, default: false
    tv:                                               # Optional
      enabled:  true                                  # Optional, default: false
    prom:                                             # Optional
      enabled: true                                   # Optional, default: false
      path: "metrics"                                 # Optional, default: ""
      pusher:                                         # Optional
        enabled: false                                # Optional, default: false
        jobName: "greeter-pusher"                     # Required
        remoteAddress: "localhost:9091"               # Required
        basicAuth: "user:pass"                        # Optional, default: ""
        intervalMs: 1000                              # Optional, default: 1000
        cert:                                         # Optional
          ref: "local-cert"                           # Optional, default: "", reference of cert entry declared above
    logger:                                           # Optional
      zapLogger:                                      # Optional
        ref: zap-logger                               # Optional, default: logger of STDOUT, reference of logger entry declared above
      eventLogger:                                    # Optional
        ref: event-logger                             # Optional, default: logger of STDOUT, reference of logger entry declared above
    interceptors:                                     # Optional
      loggingZap:
        enabled: true                                 # Optional, default: false
      metricsProm:
        enabled: true                                 # Optional, default: false
#      basicAuth:
#        enabled: true                                # Optional, default: false
#       credentials:
#         - "user:pass"                               # Optional, default: ""

```

```go
package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rookie-ninja/rk-boot"
	"net/http"
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
	// Create a new boot instance.
	boot := rkboot.NewBoot(rkboot.WithBootConfigPath("example/gin/boot.yaml"))

	// Register handler
	boot.GetGinEntry("greeter").Router.GET("/v1/hello", hello)

	// Bootstrap
	boot.Bootstrap(context.TODO())
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
| gin.name | Name of gin server entry | string | N/A |
| gin.port | Port of server | integer | nil, server won't start |
| gin.cert.ref | Reference of cert entry declared in cert section | string | "" |
| gin.sw.enabled | Enable swagger | boolean | false | 
| gin.sw.path | Swagger path | string | / |
| gin.sw.jsonPath | Swagger json file path | string | / |
| gin.sw.headers | Headers will send with swagger response | array | [] |
| gin.commonService.enabled | Enable common service | boolean | false |
| gin.tv.enabled | Enable RK TV whose path is /rk/v1/tv | boolean | false |
| gin.prom.enabled | Enable prometheus | boolean | false |
| gin.prom.path | Path of prometheus | string | metrics |
| gin.prom.cert.ref |  Reference of cert entry declared in cert section | string | "" |
| gin.prom.pusher.enabled | Enable prometheus pusher | bool | false |
| gin.prom.pusher.jobName | Job name would be attached as label while pushing to remote pushgateway | string | "" |
| gin.prom.pusher.remoteAddress | PushGateWay address, could be form of http://x.x.x.x or x.x.x.x | string | "" |
| gin.prom.pusher.intervalMs | Push interval in milliseconds | string | 1000 |
| gin.prom.pusher.basicAuth | Basic auth used to interact with remote pushgateway, form of \<user:pass\> | string | "" |
| gin.prom.pusher.cert.ref | Reference of rkentry.CertEntry | string | "" |
| gin.logger.zapLogger.ref | Reference of logger entry declared above | string | "" |
| gin.logger.eventLogger.ref | Reference of logger entry declared above | string | "" |
| gin.interceptors.loggingZap.enabled | Enable logging interceptor | boolean | false |
| gin.interceptors.metricsProm.enabled | Enable prometheus metrics for every request | boolean | false |
| gin.interceptors.basicAuth.enabled | Enable auth interceptor | boolean | false |
| gin.interceptors.basicAuth.credentials | Provide basic auth credentials, form of \<user:pass\> | string | false |

### Common Service

| path | description |
| ------ | ------ |
| /rk/v1/apis | List API |
| /rk/v1/certs | List CertEntry |
| /rk/v1/configs | List ConfigEntry |
| /rk/v1/entries | List all Entry |
| /rk/v1/gc | Trigger GC |
| /rk/v1/healthy | Get application healthy status, returns true if application is running |
| /rk/v1/info | Get application and process info |
| /rk/v1/logs | List logger related entries |
| /rk/v1/req | List prometheus metrics of requests |
| /rk/v1/sys | Get OS stat |
| /rk/v1/tv | Get HTML page of /tv |

### TV Service

| path | description |
| ------ | ------ |
| /rk/v1/tv or /rk/v1/tv/overview | Get application and process info of HTML page |
| /rk/v1/tv/api | Get API of HTML page |
| /rk/v1/tv/entry | Get entry of HTML page |
| /rk/v1/tv/config | Get config of HTML page |
| /rk/v1/tv/cert | Get cert of HTML page |
| /rk/v1/tv/os | Get OS of HTML page |
| /rk/v1/tv/env | Get Go environment of HTML page |
| /rk/v1/tv/prometheus | Get metrics of HTML page |
| /rk/v1/log | Get log of HTML page |

### Development Status: Stable

### Appendix
#### Proto file compilation
Use bellow command to rebuild proto files, we are using [buf](https://docs.buf.build/generate-usage) to generate proto related files.
Configuration could be found at root path of project.

- make buf

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