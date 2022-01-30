# rk-boot/grpc
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
  - [Supported bootstrap](#supported-bootstrap)
  - [Supported instances](#supported-instances)
  - [Supported middlewares](#supported-middlewares)
- [Installation](#installation)
- [Quick Start](#quick-start)
  - [1.Prepare .proto files](#1prepare-proto-files)
  - [2.Generate .pb.go files with buf](#2generate-pbgo-files-with-buf)
  - [3.Create boot.yaml](#3create-bootyaml)
  - [4.Create main.go](#4create-maingo)
  - [5.Start server](#5start-server)
  - [6.Validation](#6validation)
    - [6.1 gRPC & grpc-gateway server](#61-grpc--grpc-gateway-server)
    - [6.2 Swagger & TV & Prometheus Client](#62-swagger--tv--prometheus-client)
    - [6.3 Logging](#63-logging)
    - [6.4 Meta](#64-meta)
    - [6.5 Send request](#65-send-request)
    - [6.6 RPC logs](#66-rpc-logs)
- [YAML options](#yaml-options)
  - [gRPC Service](#grpc-service)
  - [gRPC gateway options](#grpc-gateway-options)
  - [Common Service](#common-service)
  - [Prometheus Client](#prometheus-client)
  - [TV Service](#tv-service)
  - [Swagger Service](#swagger-service)
  - [Static file handler Service](#static-file-handler-service)
  - [Middlewares](#middlewares)
    - [Log](#log)
    - [Metrics](#metrics)
    - [Auth](#auth)
    - [Meta](#meta)
    - [Tracing](#tracing)
    - [RateLimit](#ratelimit)
    - [Timeout](#timeout)
    - [CORS](#cors)
    - [JWT](#jwt)
    - [Secure](#secure)
    - [CSRF](#csrf)
  - [Full YAML](#full-yaml)
- [Development Status: Stable](#development-status-stable)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Overview
Middlewares & bootstrapper designed for [gRPC](https://grpc.io/docs/languages/go/) and [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway), see [official docs](https://rkdev.info/docs/bootstrapper/user-guide/grpc-golang/).

With rk-boot/grpc, user can start [gRPC](https://grpc.io/docs/languages/go/) and [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) server via boot.yaml file easily.

**Dependency**: [rk-grpc](https://github.com/rookie-ninja/rk-grpc)

## Architecture
![image](https://github.com/rookie-ninja/rk-grpc/blob/master/docs/img/grpc-arch.png?raw=true)

## Features
### Supported bootstrap
| Bootstrap | Description |
| --- | --- |
| YAML based | Start [gRPC](https://grpc.io/docs/languages/go/) and [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) microservice from YAML |
| Code based | Start [gRPC](https://grpc.io/docs/languages/go/) and [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) microservice from code |

### Supported instances
All instances could be configured via YAML or Code.

**User can enable anyone of those as needed! No mandatory binding!**

| Instance | Description |
| --- | --- |
| [gRPC](https://grpc.io/docs/languages/go/) | [gRPC](https://grpc.io/docs/languages/go/) defined with protocol buffer. |
| [gRPC](https://grpc.io/docs/languages/go/) proxy | Proxy gRPC request to another gRPC server. |
| [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) | [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) service with same port. |
| [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) options | Well defined [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) options. |
| Config | Configure [spf13/viper](https://github.com/spf13/viper) as config instance and reference it from YAML |
| Logger | Configure [uber-go/zap](https://github.com/uber-go/zap) logger configuration and reference it from YAML |
| EventLogger | Configure logging of RPC with [rk-query](https://github.com/rookie-ninja/rk-query) and reference it from YAML |
| Credential | Fetch credentials from remote datastore like ETCD. |
| Cert | Fetch TLS/SSL certificates from remote datastore like ETCD and start microservice. |
| Prometheus | Start prometheus client at client side and push metrics to [pushgateway](https://github.com/prometheus/pushgateway) as needed. |
| Swagger | Builtin swagger UI handler. |
| CommonService | List of common APIs. |
| TV | A Web UI shows microservice and environment information. |
| StaticFileHandler | A Web UI shows files could be downloaded from server, currently support source of local and pkger. |

### Supported middlewares
All middlewares could be configured via YAML or Code.

**User can enable anyone of those as needed! No mandatory binding!**

| Middleware | Description |
| --- | --- |
| Metrics | Collect RPC metrics and export to [prometheus](https://github.com/prometheus/client_golang) client. |
| Log | Log every RPC requests as event with [rk-query](https://github.com/rookie-ninja/rk-query). |
| Trace | Collect RPC trace and export it to stdout, file or jaeger with [open-telemetry/opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go). |
| Panic | Recover from panic for RPC requests and log it. |
| Meta | Send micsroservice metadata as header to client. |
| Auth | Support [Basic Auth] and [API Key] authorization types. |
| RateLimit | Limiting RPC rate globally or per path. |
| Timeout | Timing out request by configuration. |
| CORS | Server side CORS validation. |
| JWT | Server side JWT validation. |
| Secure | Server side secure validation. |
| CSRF | Server side csrf validation. |

## Installation
`go get github.com/rookie-ninja/rk-boot/grpc`

## Quick Start
In the bellow example, we will start microservice with bellow functionality and middlewares enabled via YAML.

- [gRPC](https://grpc.io/docs/languages/go/) and [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) server
- [gRPC](https://grpc.io/docs/languages/go/) server reflection
- Swagger UI
- CommonService
- TV
- Prometheus Metrics (middleware)
- Logging (middleware)
- Meta (middleware)

### 1.Prepare .proto files
- api/v1/greeter.proto

```protobuf
syntax = "proto3";

package api.v1;

option go_package = "api/v1/greeter";

service Greeter {
  rpc Greeter (GreeterRequest) returns (GreeterResponse) {}
}

message GreeterRequest {
  string name = 1;
}

message GreeterResponse {
  string msg = 1;
}
```

- api/v1/gw_mapping.yaml

```yaml
type: google.api.Service
config_version: 3

# Please refer google.api.Http in https://github.com/googleapis/googleapis/blob/master/google/api/http.proto file for details.
http:
  rules:
    - selector: api.v1.Greeter.Greeter
      get: /v1/greeter
```

- buf.yaml

```yaml
version: v1beta1
name: github.com/rk-dev/rk-boot
build:
  roots:
    - api
```

- buf.gen.yaml

```yaml
version: v1beta1
plugins:
  # protoc-gen-go needs to be installed, generate go files based on proto files
  - name: go
    out: api/gen
    opt:
     - paths=source_relative
  # protoc-gen-go-grpc needs to be installed, generate grpc go files based on proto files
  - name: go-grpc
    out: api/gen
    opt:
      - paths=source_relative
      - require_unimplemented_servers=false
  # protoc-gen-grpc-gateway needs to be installed, generate grpc-gateway go files based on proto files
  - name: grpc-gateway
    out: api/gen
    opt:
      - paths=source_relative
      - grpc_api_configuration=api/v1/gw_mapping.yaml
  # protoc-gen-openapiv2 needs to be installed, generate swagger config files based on proto files
  - name: openapiv2
    out: api/gen
    opt:
      - grpc_api_configuration=api/v1/gw_mapping.yaml
```

### 2.Generate .pb.go files with [buf](https://docs.buf.build/introduction)
```
$ buf generate --path api/v1
```

- directory hierarchy

```
.
├── api
│   ├── gen
│   │   └── v1
│   │       ├── greeter.pb.go
│   │       ├── greeter.pb.gw.go
│   │       ├── greeter.swagger.json
│   │       └── greeter_grpc.pb.go
│   └── v1
│       ├── greeter.proto
│       └── gw_mapping.yaml
├── boot.yaml
├── buf.gen.yaml
├── buf.yaml
├── go.mod
├── go.sum
└── main.go
```

### 3.Create boot.yaml
- Enable [gRPC](https://grpc.io/docs/languages/go/) reflection via [grpc.enableReflection: true]
- Use RK style [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) [option](https://github.com/rookie-ninja/rk-grpc/blob/master/boot/gw_server_options.go) via [enableRkGwOption: true]
- Configure json file path of swagger config file via [grpc.sw.jsonPath: "api/gen/v1"]

```yaml
---
grpc:
  - name: greeter                     # Required
    port: 8080                        # Required
    enabled: true                     # Required
    enableReflection: true            # Optional, default: false
    enableRkGwOption: true            # Optional, default: false
    commonService:
      enabled: true                   # Optional, default: false
    tv:
      enabled: true                   # Optional, default: false
    sw:
      enabled: true                   # Optional, default: false
      jsonPath: "api/gen/v1"          # Optional, default: ""
    prom:
      enabled: true                   # Optional, default: false
    interceptors:
      loggingZap:
        enabled: true                 # Optional, default: false
      metricsProm:
        enabled: true                 # Optional, default: false
      meta:
        enabled: true                 # Optional, default: false
```

### 4.Create main.go
Implement proto method. Both gRPC and grpc-gateway registration functions needs to be added via [grpcEntry.AddRegFuncGrpc] and [grpcEntry.AddRegFuncGw]

```go
// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"fmt"
	"github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-boot/grpc"
	"github.com/rookie-ninja/rk-demo/api/gen/v1"
	"google.golang.org/grpc"
)

// Application entrance.
func main() {
	// Create a new boot instance.
	boot := rkboot.NewBoot()

	// Get grpc entry with name
	grpcEntry := rkbootgrpc.GetGrpcEntry("greeter")
	grpcEntry.AddRegFuncGrpc(registerGreeter)
	grpcEntry.AddRegFuncGw(greeter.RegisterGreeterHandlerFromEndpoint)

	// Bootstrap
	boot.Bootstrap(context.Background())

	// Wait for shutdown sig
	boot.WaitForShutdownSig(context.Background())
}

func registerGreeter(server *grpc.Server) {
	greeter.RegisterGreeterServer(server, &GreeterServer{})
}

type GreeterServer struct{}

func (server *GreeterServer) Greeter(ctx context.Context, request *greeter.GreeterRequest) (*greeter.GreeterResponse, error) {
	return &greeter.GreeterResponse{
		Msg: fmt.Sprintf("Hello %s!", request.Name),
	}, nil
}
```

### 5.Start server
```
$ go run main.go
```

### 6.Validation
#### 6.1 gRPC & grpc-gateway server
Try to test [gRPC](https://grpc.io/docs/languages/go/) & [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) Service with [curl](https://curl.se/) & [grpcurl](https://github.com/fullstorydev/grpcurl)

```shell script
# Curl to common service
$ curl localhost:8080/rk/v1/healthy
{"healthy":true}
```

[grpcurl](https://github.com/fullstorydev/grpcurl) is working when [enableReflection: true] is set.

```shell script
# List grpc services at port 8080 without TLS
# Expect RkCommonService since we enabled common services.
$ grpcurl -plaintext localhost:8080 list                           
api.v1.Greeter
grpc.reflection.v1alpha.ServerReflection
rk.api.v1.RkCommonService

# List grpc methods in rk.api.v1.RkCommonService
$ grpcurl -plaintext localhost:8080 list rk.api.v1.RkCommonService            
rk.api.v1.RkCommonService.Apis
rk.api.v1.RkCommonService.Certs
rk.api.v1.RkCommonService.Configs
rk.api.v1.RkCommonService.Deps
rk.api.v1.RkCommonService.Entries
rk.api.v1.RkCommonService.Gc
rk.api.v1.RkCommonService.GwErrorMapping
rk.api.v1.RkCommonService.Healthy
rk.api.v1.RkCommonService.Info
rk.api.v1.RkCommonService.License
rk.api.v1.RkCommonService.Logs
rk.api.v1.RkCommonService.Readme
rk.api.v1.RkCommonService.Req
rk.api.v1.RkCommonService.Sys
rk.api.v1.RkCommonService.Git

# Send request to rk.api.v1.RkCommonService.Healthy
$ grpcurl -plaintext localhost:8080 rk.api.v1.RkCommonService.Healthy
{
    "healthy": true
}
```

#### 6.2 Swagger & TV & Prometheus Client
- Swagger UI: [http://localhost:8080/sw](http://localhost:8080/sw)
- TV: [http://localhost:8080/rk/v1/tv](http://localhost:8080/rk/v1/tv)
- Prometheus Client: [http://localhost:8080/metrics](http://localhost:8080/metrics)

#### 6.3 Logging
By default, we enable zap logger and event logger with encoding type of [console]. Encoding type of [json] is also supported.

```shell script
2022-01-08T20:04:58.389+0800    INFO    boot/grpc_entry.go:1515 Bootstrap grpcEntry     {"eventId": "95c613c8-1589-4357-b386-682c29077eea", "entryName": "greeter"}
------------------------------------------------------------------------
endTime=2022-01-08T20:04:58.391675+08:00
startTime=2022-01-08T20:04:58.389112+08:00
elapsedNano=2562884
timezone=CST
ids={"eventId":"95c613c8-1589-4357-b386-682c29077eea"}
app={"appName":"rk","appVersion":"","entryName":"greeter","entryType":"GrpcEntry"}
env={"arch":"amd64","az":"*","domain":"*","hostname":"lark.local","localIP":"10.8.0.2","os":"darwin","realm":"*","region":"*"}
payloads={"commonServiceEnabled":true,"commonServicePathPrefix":"/rk/v1/","grpcPort":8080,"gwPort":8080,"promEnabled":true,"promPath":"/metrics","promPort":8080,"swEnabled":true,"swPath":"/sw/","tvEnabled":true,"tvPath":"/rk/v1/tv/"}
error={}
counters={}
pairs={}
timing={}
remoteAddr=localhost
operation=Bootstrap
resCode=OK
eventStatus=Ended
EOE
```

#### 6.4 Meta
By default, we will send back some metadata to client with headers.

```shell script
$ curl -vs localhost:8080/rk/v1/healthy
...
< HTTP/1.1 200 OK
< Content-Type: application/json
< X-Request-Id: df25d9f8-0658-4b80-834b-4d8af0690eb3
< X-Rk-App-Name: rk
< X-Rk-App-Unix-Time: 2022-01-08T20:05:20.27004+08:00
< X-Rk-App-Version: 
< X-Rk-Received-Time: 2022-01-08T20:05:20.27004+08:00
< Date: Sat, 08 Jan 2022 12:05:20 GMT
...
```

#### 6.5 Send request
We registered /v1/greeter API in [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) server and let's validate it!

```shell script
$ curl "localhost:8080/v1/greeter?name=rk-dev"
{"msg":"Hello rk-dev!"}
```

We registered api.v1.Greeter.Greeter API in [gRPC](https://grpc.io/docs/languages/go/) server and let's validate it!

```shell script
$ grpcurl -d '{"name":"rk-dev"}' -plaintext localhost:8080 api.v1.Greeter.Greeter
{
  "msg": "Hello rk-dev!"
}
```

#### 6.6 RPC logs
Bellow logs would be printed in stdout.

The first block of log is from [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) request.

The second block of log is from [gRPC](https://grpc.io/docs/languages/go/) request.

```
------------------------------------------------------------------------
endTime=2022-01-08T23:51:02.876796+08:00
startTime=2022-01-08T23:51:02.876533+08:00
elapsedNano=263232
timezone=CST
ids={"eventId":"ebd063f1-c413-4bf9-8dfb-38e4e109c869","requestId":"ebd063f1-c413-4bf9-8dfb-38e4e109c869"}
app={"appName":"rk","appVersion":"","entryName":"greeter","entryType":"GrpcEntry"}
env={"arch":"amd64","az":"*","domain":"*","hostname":"lark.local","localIP":"10.8.0.2","os":"darwin","realm":"*","region":"*"}
payloads={"grpcMethod":"Greeter","grpcService":"api.v1.Greeter","grpcType":"unaryServer","gwMethod":"GET","gwPath":"/v1/greeter","gwScheme":"http","gwUserAgent":"curl/7.64.1"}
error={}
counters={}
pairs={}
timing={}
remoteAddr=127.0.0.1:55386
operation=/api.v1.Greeter/Greeter
resCode=OK
eventStatus=Ended
EOE
------------------------------------------------------------------------
endTime=2022-01-08T23:51:19.761892+08:00
startTime=2022-01-08T23:51:19.761844+08:00
elapsedNano=48023
timezone=CST
ids={"eventId":"17e9fd95-037f-4d7c-a647-08809e8680ed","requestId":"17e9fd95-037f-4d7c-a647-08809e8680ed"}
app={"appName":"rk","appVersion":"","entryName":"greeter","entryType":"GrpcEntry"}
env={"arch":"amd64","az":"*","domain":"*","hostname":"lark.local","localIP":"10.8.0.2","os":"darwin","realm":"*","region":"*"}
payloads={"grpcMethod":"Greeter","grpcService":"api.v1.Greeter","grpcType":"unaryServer","gwMethod":"","gwPath":"","gwScheme":"","gwUserAgent":""}
error={}
counters={}
pairs={}
timing={}
remoteAddr=127.0.0.1:56483
operation=/api.v1.Greeter/Greeter
resCode=OK
eventStatus=Ended
EOE
```

## YAML options
User can start multiple [gRPC](https://grpc.io/docs/languages/go/) and [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) instances at the same time. Please make sure use different port and name.

### gRPC Service
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.name | The name of [gRPC](https://grpc.io/docs/languages/go/) server | string | N/A |
| grpc.enabled | Enable [gRPC](https://grpc.io/docs/languages/go/) entry | bool | false |
| grpc.port | The port of [gRPC](https://grpc.io/docs/languages/go/) server | integer | nil, server won't start |
| grpc.description | Description of [gRPC](https://grpc.io/docs/languages/go/) entry. | string | "" |
| grpc.enableReflection | Enable [gRPC](https://grpc.io/docs/languages/go/) server reflection | boolean | false |
| grpc.enableRkGwOption | Enable RK style [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) server options. [detail](boot/gw_server_options.go) | false |
| grpc.noRecvMsgSizeLimit | Disable [gRPC](https://grpc.io/docs/languages/go/) server side receive message size limit | false |
| grpc.gwMappingFilePaths | The grpc [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) mapping file path. [example](boot/api/v1/gw_mapping.yaml) | string array | [] |
| grpc.cert.ref | Reference of cert entry declared in [cert entry](https://github.com/rookie-ninja/rk-entry#certentry) | string | "" |
| grpc.logger.zapLogger.ref | Reference of zapLoggerEntry declared in [zapLoggerEntry](https://github.com/rookie-ninja/rk-entry#zaploggerentry) | string | "" |
| grpc.logger.eventLogger.ref | Reference of eventLoggerEntry declared in [eventLoggerEntry](https://github.com/rookie-ninja/rk-entry#eventloggerentry) | string | "" |

```yaml
grpc:
  - name: greeter                                          # Required
    enabled: true                                          # Required
    port: 8080                                             # Required
    description: "greeter server"                          # Optional, default: ""
    enableReflection: true                                 # Optional, default: false, gRPC reflection for grpcurl cli
    enableRkGwOption: true                                 # Optional, default: false, RK style grpc-gateway option
    noRecvMsgSizeLimit: true                               # Optional, default: false, disable server side receive size limit
    gwMappingFilePaths: []                                 # Optional, default: []
    cert:
      ref: "local-cert"                                    # Optional, default: "", reference of cert entry declared above
    logger:
      zapLogger:
        ref: zap-logger                                    # Optional, default: logger of STDOUT, reference of logger entry declared above
      eventLogger:
        ref: event-logger                                  # Optional, default: logger of STDOUT, reference of logger entry declared above
```

### gRPC gateway options
Please refer to bellow repository for detailed explanations.
- [protobuf-go/encoding/protojson/encode.go](https://github.com/protocolbuffers/protobuf-go/blob/master/encoding/protojson/encode.go#L43)
- [protobuf-go/encoding/protojson/decode.go ](https://github.com/protocolbuffers/protobuf-go/blob/master/encoding/protojson/decode.go#L33)

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.gwOption.marshal.multiline | Enable multiline in [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) marshaller | bool | false |
| grpc.gwOption.marshal.emitUnpopulated | Enable emitUnpopulated in [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) marshaller | bool | false |
| grpc.gwOption.marshal.indent | Set indent in [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) marshaller | string | "  " |
| grpc.gwOption.marshal.allowPartial | Enable allowPartial in [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) marshaller | bool | false |
| grpc.gwOption.marshal.useProtoNames | Enable useProtoNames in [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) marshaller | bool | false |
| grpc.gwOption.marshal.useEnumNumbers | Enable useEnumNumbers in [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) marshaller | bool | false |
| grpc.gwOption.unmarshal.allowPartial | Enable allowPartial in [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) unmarshaler | bool | false |
| grpc.gwOption.unmarshal.discardUnknown | Enable discardUnknown in [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) unmarshaler | bool | false |

```yaml
grpc:
  - name: greeter                                     # Required
    port: 8080                                        # Required
    enabled: true                                     # Required
    enableRkGwOption: true                            # Optional, default: false
    gwOption:                                         # Optional, default: nil
      marshal:                                        # Optional, default: nil
        multiline: false                              # Optional, default: false
        emitUnpopulated: false                        # Optional, default: false
        indent: ""                                    # Optional, default: false
        allowPartial: false                           # Optional, default: false
        useProtoNames: false                          # Optional, default: false
        useEnumNumbers: false                         # Optional, default: false
      unmarshal:                                      # Optional, default: nil
        allowPartial: false                           # Optional, default: false
        discardUnknown: false                         # Optional, default: false
```

### Common Service
```yaml
http:
  rules:
    - selector: rk.api.v1.RkCommonService.Healthy
      get: /rk/v1/healthy
    - selector: rk.api.v1.RkCommonService.Gc
      get: /rk/v1/gc
    - selector: rk.api.v1.RkCommonService.Info
      get: /rk/v1/info
    - selector: rk.api.v1.RkCommonService.Configs
      get: /rk/v1/configs
    - selector: rk.api.v1.RkCommonService.Apis
      get: /rk/v1/apis
    - selector: rk.api.v1.RkCommonService.Sys
      get: /rk/v1/sys
    - selector: rk.api.v1.RkCommonService.Req
      get: /rk/v1/req
    - selector: rk.api.v1.RkCommonService.Entries
      get: /rk/v1/entries
    - selector: rk.api.v1.RkCommonService.Certs
      get: /rk/v1/certs
    - selector: rk.api.v1.RkCommonService.Logs
      get: /rk/v1/logs
    - selector: rk.api.v1.RkCommonService.Deps
      get: /rk/v1/deps
    - selector: rk.api.v1.RkCommonService.License
      get: /rk/v1/license
    - selector: rk.api.v1.RkCommonService.Readme
      get: /rk/v1/readme
    - selector: rk.api.v1.RkCommonService.Git
      get: /rk/v1/git
    - selector: rk.api.v1.RkCommonService.GwErrorMapping
      get: /rk/v1/gwErrorMapping
```

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.commonService.enabled | Enable embedded common service | boolean | false |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    commonService:
      enabled: true                                        # Optional, default: false
```

### Prometheus Client
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.prom.enabled | Enable prometheus | boolean | false |
| grpc.prom.path | Path of prometheus | string | /metrics |
| grpc.prom.pusher.enabled | Enable prometheus pusher | bool | false |
| grpc.prom.pusher.jobName | Job name would be attached as label while pushing to remote [pushgateway](https://github.com/prometheus/pushgateway) | string | "" |
| grpc.prom.pusher.remoteAddress | [pushgateway](https://github.com/prometheus/pushgateway) address, could be form of http://x.x.x.x or x.x.x.x | string | "" |
| grpc.prom.pusher.intervalMs | Push interval in milliseconds | string | 1000 |
| grpc.prom.pusher.basicAuth | Basic auth used to interact with remote [pushgateway](https://github.com/prometheus/pushgateway), form of [user:pass] | string | "" |
| grpc.prom.pusher.cert.ref | Reference of rkentry.CertEntry | string | "" |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    prom:
      enabled: true                                        # Optional, default: false
      path: ""                                             # Optional, default: "metrics"
      pusher:                                              # pushgateway configuration
        enabled: false                                     # Optional, default: false
        jobName: "greeter-pusher"                          # Required
        remoteAddress: "localhost:9091"                    # Required
        basicAuth: "user:pass"                             # Optional, default: ""
        intervalMs: 10000                                  # Optional, default: 1000
        cert:                                              # Optional
          ref: "local-test"                                # Optional, default: "", reference of cert entry declared above
```

### TV Service
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.tv.enabled | Enable RK TV | boolean | false |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    tv:
      enabled: true                                        # Optional, default: false
```

### Swagger Service
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.sw.enabled | Enable swagger service over [gRPC](https://grpc.io/docs/languages/go/) server | boolean | false |
| grpc.sw.path | The path access swagger service from web | string | /sw |
| grpc.sw.jsonPath | Where the swagger.json files are stored locally | string | "" |
| grpc.sw.headers | Headers would be sent to caller as scheme of [key:value] | []string | [] |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    sw:
      enabled: true                                        # Optional, default: false
      path: "sw"                                           # Optional, default: "sw"
      jsonPath: ""                                         # Optional, default: ""
      headers: ["sw:rk"]                                   # Optional, default: []
```

### Static file handler Service
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.static.enabled | Optional, Enable static file handler | boolean | false |
| grpc.static.path | Optional, path of static file handler | string | /rk/v1/static |
| grpc.static.sourceType | Required, local and pkger supported | string | "" |
| grpc.static.sourcePath | Required, full path of source directory | string | "" |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    static:
      enabled: true                                        # Optional, default: false
      path: "/rk/v1/static"                                # Optional, default: /rk/v1/static
      sourceType: local                                    # Required, options: pkger, local
      sourcePath: "."                                      # Required, full path of source directory
```

- About [pkger](https://github.com/markbates/pkger)
User can use pkger command line tool to embed static files into .go files.

Please use sourcePath like: github.com/rookie-ninja/rk-grpc:/boot/assets

### Middlewares
#### Log
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.loggingZap.enabled | Enable log interceptor | boolean | false |
| grpc.interceptors.loggingZap.zapLoggerEncoding | json or console | string | console |
| grpc.interceptors.loggingZap.zapLoggerOutputPaths | Output paths | []string | stdout |
| grpc.interceptors.loggingZap.eventLoggerEncoding | json or console | string | console |
| grpc.interceptors.loggingZap.eventLoggerOutputPaths | Output paths | []string | false |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    interceptors:
      loggingZap:
        enabled: true                                      # Optional, default: false
        zapLoggerEncoding: "json"                          # Optional, default: "console"
        zapLoggerOutputPaths: ["logs/app.log"]             # Optional, default: ["stdout"]
        eventLoggerEncoding: "json"                        # Optional, default: "console"
        eventLoggerOutputPaths: ["logs/event.log"]         # Optional, default: ["stdout"]
```

We will log two types of log for every RPC call.
- zapLogger

Contains user printed logging with requestId or traceId.

- eventLogger

Contains per RPC metadata, response information, environment information and etc.

| Field | Description |
| ---- | ---- |
| endTime | As name described |
| startTime | As name described |
| elapsedNano | Elapsed time for RPC in nanoseconds |
| timezone | As name described |
| ids | Contains three different ids(eventId, requestId and traceId). If meta interceptor was enabled or event.SetRequestId() was called by user, then requestId would be attached. eventId would be the same as requestId if meta interceptor was enabled. If trace interceptor was enabled, then traceId would be attached. |
| app | Contains [appName, appVersion](https://github.com/rookie-ninja/rk-entry#appinfoentry), entryName, entryType. |
| env | Contains arch, az, domain, hostname, localIP, os, realm, region. realm, region, az, domain were retrieved from environment variable named as REALM, REGION, AZ and DOMAIN. "*" means empty environment variable.|
| payloads | Contains RPC related metadata |
| error | Contains errors if occur |
| counters | Set by calling event.SetCounter() by user. |
| pairs | Set by calling event.AddPair() by user. |
| timing | Set by calling event.StartTimer() and event.EndTimer() by user. |
| remoteAddr |  As name described |
| operation | RPC method name |
| resCode | Response code of RPC |
| eventStatus | Ended or InProgress |

- example

```shell script
------------------------------------------------------------------------
endTime=2021-06-24T05:58:48.282193+08:00
startTime=2021-06-24T05:58:48.28204+08:00
elapsedNano=153005
timezone=CST
ids={"eventId":"573ce6a8-308b-4fc0-9255-33608b9e41d4","requestId":"573ce6a8-308b-4fc0-9255-33608b9e41d4"}
app={"appName":"rk-grpc","appVersion":"master-xxx","entryName":"greeter","entryType":"GrpcEntry"}
env={"arch":"amd64","az":"*","domain":"*","hostname":"lark.local","localIP":"10.8.0.6","os":"darwin","realm":"*","region":"*"}
payloads={"grpcMethod":"Healthy","grpcService":"rk.api.v1.RkCommonService","grpcType":"unaryServer","gwMethod":"GET","gwPath":"/rk/v1/healthy","gwScheme":"http","gwUserAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36"}
error={}
counters={}
pairs={"healthy":"true"}
timing={}
remoteAddr=localhost:57135
operation=/rk.api.v1.RkCommonService/Healthy
resCode=OK
eventStatus=Ended
EOE
```

#### Metrics
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.metricsProm.enabled | Enable metrics interceptor | boolean | false |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    prom:
      enabled: true                                        # Optional, default: false
    interceptors:
      metricsProm:
        enabled: true                                      # Optional, default: false
```

#### Auth
Enable the server side auth. codes.Unauthenticated would be returned to client if not authorized with user defined credential.

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.auth.enabled | Enable auth interceptor | boolean | false |
| grpc.interceptors.auth.basic | Basic auth credentials as scheme of <user:pass> | []string | [] |
| grpc.interceptors.auth.apiKey | API key | []string | [] |
| grpc.interceptors.auth.ignorePrefix | The paths of prefix that will be ignored by interceptor | []string | [] |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    interceptors:
      auth:
        enabled: true                                      # Optional, default: false
        basic:
          - "user:pass"                                    # Optional, default: []
        ignorePrefix:
          - "/rk/v1"                                       # Optional, default: []
        apiKey:
          - "keys"                                         # Optional, default: []
```

#### Meta
Send application metadata as header to client and [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.meta.enabled | Enable meta interceptor | boolean | false |
| grpc.interceptors.meta.prefix | Header key was formed as X-<Prefix>-XXX | string | RK |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    interceptors:
      meta:
        enabled: true                                      # Optional, default: false
        prefix: "rk"                                       # Optional, default: "rk"
```

#### Tracing
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.tracingTelemetry.enabled | Enable tracing interceptor | boolean | false |
| grpc.interceptors.tracingTelemetry.exporter.file.enabled | Enable file exporter | boolean | false |
| grpc.interceptors.tracingTelemetry.exporter.file.outputPath | Export tracing info to files | string | stdout |
| grpc.interceptors.tracingTelemetry.exporter.jaeger.agent.enabled | Export tracing info to jaeger agent | boolean | false |
| grpc.interceptors.tracingTelemetry.exporter.jaeger.agent.host | As name described | string | localhost |
| grpc.interceptors.tracingTelemetry.exporter.jaeger.agent.port | As name described | int | 6831 |
| grpc.interceptors.tracingTelemetry.exporter.jaeger.collector.enabled | Export tracing info to jaeger collector | boolean | false |
| grpc.interceptors.tracingTelemetry.exporter.jaeger.collector.endpoint | As name described | string | http://localhost:16368/api/trace |
| grpc.interceptors.tracingTelemetry.exporter.jaeger.collector.username | As name described | string | "" |
| grpc.interceptors.tracingTelemetry.exporter.jaeger.collector.password | As name described | string | "" |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    interceptors:
      tracingTelemetry:
        enabled: true                                      # Optional, default: false
        exporter:                                          # Optional, default will create a stdout exporter
          file:
            enabled: true                                  # Optional, default: false
            outputPath: "logs/trace.log"                   # Optional, default: stdout
          jaeger:
            agent:
              enabled: false                               # Optional, default: false
              host: ""                                     # Optional, default: localhost
              port: 0                                      # Optional, default: 6831
            collector:
              enabled: true                                # Optional, default: false
              endpoint: ""                                 # Optional, default: http://localhost:14268/api/traces
              username: ""                                 # Optional, default: ""
              password: ""                                 # Optional, default: ""
```

#### RateLimit
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.rateLimit.enabled | Enable rate limit interceptor | boolean | false |
| grpc.interceptors.rateLimit.algorithm | Provide algorithm, tokenBucket and leakyBucket are available options | string | tokenBucket |
| grpc.interceptors.rateLimit.reqPerSec | Request per second globally | int | 0 |
| grpc.interceptors.rateLimit.paths.path | [gRPC](https://grpc.io/docs/languages/go/) full name | string | "" |
| grpc.interceptors.rateLimit.paths.reqPerSec | Request per second by [gRPC](https://grpc.io/docs/languages/go/) full method name | int | 0 |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    interceptors:
      rateLimit:
        enabled: false                                     # Optional, default: false
        algorithm: "leakyBucket"                           # Optional, default: "tokenBucket", options: [tokenBucket, leakyBucket]
        reqPerSec: 100                                     # Optional, default: 1000000
        paths:
          - path: "/rk/v1/healthy"                         # Optional, default: ""
            reqPerSec: 0                                   # Optional, default: 1000000
```

#### Timeout
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.timeout.enabled | Enable timeout interceptor | boolean | false |
| grpc.interceptors.timeout.timeoutMs | Global timeout in milliseconds. | int | 5000 |
| grpc.interceptors.timeout.paths.path | Full path | string | "" |
| grpc.interceptors.timeout.paths.timeoutMs | Timeout in milliseconds by full path | int | 5000 |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    interceptors:
      timeout:
        enabled: false                                     # Optional, default: false
        timeoutMs: 5000                                    # Optional, default: 5000
        paths:
          - path: "/rk/v1/healthy"                         # Optional, default: ""
            timeoutMs: 1000                                # Optional, default: 5000
```

#### CORS
Middleware for [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.cors.enabled | Enable cors interceptor | boolean | false |
| grpc.interceptors.cors.allowOrigins | Provide allowed origins with wildcard enabled. | []string | * |
| grpc.interceptors.cors.allowMethods | Provide allowed methods returns as response header of OPTIONS request. | []string | All http methods |
| grpc.interceptors.cors.allowHeaders | Provide allowed headers returns as response header of OPTIONS request. | []string | Headers from request |
| grpc.interceptors.cors.allowCredentials | Returns as response header of OPTIONS request. | bool | false |
| grpc.interceptors.cors.exposeHeaders | Provide exposed headers returns as response header of OPTIONS request. | []string | "" |
| grpc.interceptors.cors.maxAge | Provide max age returns as response header of OPTIONS request. | int | 0 |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    interceptors:
      cors:
        enabled: true                                     # Optional, default: false
        allowOrigins:
          - "http://localhost:*"                          # Optional, default: *
        allowCredentials: false                           # Optional, default: false
        allowHeaders: []                                  # Optional, default: []
        allowMethods: []                                  # Optional, default: []
        exposeHeaders: []                                 # Optional, default: []
        maxAge: 0                                         # Optional, default: 0
```

#### JWT
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.jwt.enabled | Enable JWT interceptor | boolean | false |
| grpc.interceptors.jwt.signingKey | Required, Provide signing key. | string | "" |
| grpc.interceptors.jwt.ignorePrefix | Provide ignoring path prefix. | []string | [] |
| grpc.interceptors.jwt.signingKeys | Provide signing keys as scheme of <key>:<value>. | []string | [] |
| grpc.interceptors.jwt.signingAlgo | Provide signing algorithm. | string | HS256 |
| grpc.interceptors.jwt.tokenLookup | Provide token lookup scheme, please see bellow description. | string | "header:Authorization" |
| grpc.interceptors.jwt.authScheme | Provide auth scheme. | string | Bearer |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    interceptors:
      jwt:
        enabled: true                                      # Optional, default: false
        signingKey: "my-secret"                            # Required
        ignorePrefix:                                      # Optional, default: []
          - "/rk/v1/tv"
          - "/sw"
          - "/rk/v1/assets"
        signingKeys:                                       # Optional
          - "key:value"
        signingAlgo: ""                                    # Optional, default: "HS256"
        tokenLookup: "header:<name>"                       # Optional, default: "header:Authorization"
        authScheme: "Bearer"                               # Optional, default: "Bearer"
```

The supported scheme of **tokenLookup** 

```
// Optional. Default value "header:Authorization".
// Possible values:
// - "header:<name>"
// Multiply sources example:
// - "header: Authorization,cookie: myowncookie"
```

#### Secure
Middleware for [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.secure.enabled | Enable secure interceptor | boolean | false |
| grpc.interceptors.secure.xssProtection | X-XSS-Protection header value. | string | "1; mode=block" |
| grpc.interceptors.secure.contentTypeNosniff | X-Content-Type-Options header value. | string | nosniff |
| grpc.interceptors.secure.xFrameOptions | X-Frame-Options header value. | string | SAMEORIGIN |
| grpc.interceptors.secure.hstsMaxAge | Strict-Transport-Security header value. | int | 0 |
| grpc.interceptors.secure.hstsExcludeSubdomains | Excluding subdomains of HSTS. | bool | false |
| grpc.interceptors.secure.hstsPreloadEnabled | Enabling HSTS preload. | bool | false |
| grpc.interceptors.secure.contentSecurityPolicy | Content-Security-Policy header value. | string | "" |
| grpc.interceptors.secure.cspReportOnly | Content-Security-Policy-Report-Only header value. | bool | false |
| grpc.interceptors.secure.referrerPolicy | Referrer-Policy header value. | string | "" |
| grpc.interceptors.secure.ignorePrefix | Ignoring path prefix. | []string | [] |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    interceptors:
      secure:
        enabled: true                                     # Optional, default: false
        xssProtection: ""                                 # Optional, default: "1; mode=block"
        contentTypeNosniff: ""                            # Optional, default: nosniff
        xFrameOptions: ""                                 # Optional, default: SAMEORIGIN
        hstsMaxAge: 0                                     # Optional, default: 0
        hstsExcludeSubdomains: false                      # Optional, default: false
        hstsPreloadEnabled: false                         # Optional, default: false
        contentSecurityPolicy: ""                         # Optional, default: ""
        cspReportOnly: false                              # Optional, default: false
        referrerPolicy: ""                                # Optional, default: ""
        ignorePrefix: []                                  # Optional, default: []
```

#### CSRF
Middleware for [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.csrf.enabled | Enable csrf interceptor | boolean | false |
| grpc.interceptors.csrf.tokenLength | Provide the length of the generated token. | int | 32 |
| grpc.interceptors.csrf.tokenLookup | Provide csrf token lookup rules, please see code comments for details. | string | "header:X-CSRF-Token" |
| grpc.interceptors.csrf.cookieName | Provide name of the CSRF cookie. This cookie will store CSRF token. | string | _csrf |
| grpc.interceptors.csrf.cookieDomain | Domain of the CSRF cookie. | string | "" |
| grpc.interceptors.csrf.cookiePath | Path of the CSRF cookie. | string | "" |
| grpc.interceptors.csrf.cookieMaxAge | Provide max age (in seconds) of the CSRF cookie. | int | 86400 |
| grpc.interceptors.csrf.cookieHttpOnly | Indicates if CSRF cookie is HTTP only. | bool | false |
| grpc.interceptors.csrf.cookieSameSite | Indicates SameSite mode of the CSRF cookie. Options: lax, strict, none, default | string | default |
| grpc.interceptors.csrf.ignorePrefix | Ignoring path prefix. | []string | [] |

```yaml
grpc:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    interceptors:
      csrf:
        enabled: true
        tokenLength: 32                                   # Optional, default: 32
        tokenLookup: "header:X-CSRF-Token"                # Optional, default: "header:X-CSRF-Token"
        cookieName: "_csrf"                               # Optional, default: _csrf
        cookieDomain: ""                                  # Optional, default: ""
        cookiePath: ""                                    # Optional, default: ""
        cookieMaxAge: 86400                               # Optional, default: 86400
        cookieHttpOnly: false                             # Optional, default: false
        cookieSameSite: "default"                         # Optional, default: "default", options: lax, strict, none, default
        ignorePrefix: []                                  # Optional, default: []
```

### Full YAML
```yaml
---
#app:
#  description: "this is description"                      # Optional, default: ""
#  keywords: ["rk", "golang"]                              # Optional, default: []
#  homeUrl: "http://example.com"                           # Optional, default: ""
#  iconUrl: "http://example.com"                           # Optional, default: ""
#  docsUrl: ["http://example.com"]                         # Optional, default: []
#  maintainers: ["rk-dev"]                                 # Optional, default: []
#zapLogger:
#  - name: zap-logger                                      # Required
#    description: "Description of entry"                   # Optional
#    zap:
#      level: info                                         # Optional, default: info, options: [debug, DEBUG, info, INFO, warn, WARN, dpanic, DPANIC, panic, PANIC, fatal, FATAL]
#      development: true                                   # Optional, default: true
#      disableCaller: false                                # Optional, default: false
#      disableStacktrace: true                             # Optional, default: true
#      sampling:
#        initial: 0                                        # Optional, default: 0
#        thereafter: 0                                     # Optional, default: 0
#      encoding: console                                   # Optional, default: "console", options: [console, json]
#      encoderConfig:
#        messageKey: "msg"                                 # Optional, default: "msg"
#        levelKey: "level"                                 # Optional, default: "level"
#        timeKey: "ts"                                     # Optional, default: "ts"
#        nameKey: "logger"                                 # Optional, default: "logger"
#        callerKey: "caller"                               # Optional, default: "caller"
#        functionKey: ""                                   # Optional, default: ""
#        stacktraceKey: "msg"                              # Optional, default: "msg"
#        lineEnding: "\n"                                  # Optional, default: "\n"
#        levelEncoder: "capitalColor"                      # Optional, default: "capitalColor", options: [capital, capitalColor, color, lowercase]
#        timeEncoder: "iso8601"                            # Optional, default: "iso8601", options: [rfc3339nano, RFC3339Nano, rfc3339, RFC3339, iso8601, ISO8601, millis, nanos]
#        durationEncoder: "string"                         # Optional, default: "string", options: [string, nanos, ms]
#        callerEncoder: ""                                 # Optional, default: ""
#        nameEncoder: ""                                   # Optional, default: ""
#        consoleSeparator: ""                              # Optional, default: ""
#      outputPaths: [ "stdout" ]                           # Optional, default: ["stdout"], stdout would be replaced if specified
#      errorOutputPaths: [ "stderr" ]                      # Optional, default: ["stderr"], stderr would be replaced if specified
#      initialFields:                                      # Optional, default: empty map
#        key: "value"
#    lumberjack:
#      filename: "rkapp.log"                               # Optional, default: It uses <processname>-lumberjack.log in os.TempDir() if empty.
#      maxsize: 1024                                       # Optional, default: 1024 (MB)
#      maxage: 7                                           # Optional, default: 7 (days)
#      maxbackups: 3                                       # Optional, default: 3 (days)
#      localtime: true                                     # Optional, default: true
#      compress: true                                      # Optional, default: true
#    loki:
#      enabled: true                                       # Optional, default: false
#      addr: localhost:3100                                # Optional, default: localhost:3100
#      path: /loki/api/v1/push                             # Optional, default: /loki/api/v1/push
#      username: ""                                        # Optional, default: ""
#      password: ""                                        # Optional, default: ""
#      maxBatchWaitMs: 3000                                # Optional, default: 3000
#      maxBatchSize: 1000                                  # Optional, default: 1000
#      insecureSkipVerify: false                           # Optional, default: false
#      labels:                                             # Optional, default: empty map
#        my_label_key: my_label_value
#eventLogger:
#  - name: event-logger                                    # Required
#    encoding: "json"                                      # Optional, default: console, options: [json, console]
#    outputPaths: []                                       # Optional, default: ["stdout"], stdout would be replaced if specified
#    lumberjack:
#      filename: "rkapp.log"                               # Optional, default: It uses <processname>-lumberjack.log in os.TempDir() if empty.
#      maxsize: 1024                                       # Optional, default: 1024 (MB)
#      maxage: 7                                           # Optional, default: 7 (days)
#      maxbackups: 3                                       # Optional, default: 3 (days)
#      localtime: true                                     # Optional, default: true
#      compress: true                                      # Optional, default: true
#    loki:
#      enabled: true                                       # Optional, default: false
#      addr: localhost:3100                                # Optional, default: localhost:3100
#      path: /loki/api/v1/push                             # Optional, default: /loki/api/v1/push
#      username: ""                                        # Optional, default: ""
#      password: ""                                        # Optional, default: ""
#      maxBatchWaitMs: 3000                                # Optional, default: 3000
#      maxBatchSize: 1000                                  # Optional, default: 1000
#      insecureSkipVerify: false                           # Optional, default: false
#      labels:                                             # Optional, default: empty map
#        my_label_key: my_label_value
#cred:
#  - name: "local-cred"                                    # Required
#    description: "Description of entry"                   # Optional
#    provider: "localFs"                                   # Required, etcd, consul, localFs, remoteFs are supported options
#    locale: "*::*::*::*"                                  # Required, default: ""
#    paths:                                                # Optional
#      - "example/boot/full/cred.yaml"
#cert:                                                     # Optional
#  - name: "local-cert"                                    # Required
#    provider: "localFs"                                   # Required, etcd, consul, localFs, remoteFs are supported options
#    locale: "*::*::*::*"                                  # Required, default: ""
#    description: "Description of entry"                   # Optional
#    serverCertPath: "example/boot/full/server.pem"        # Optional, default: "", path of certificate on local FS
#    serverKeyPath: "example/boot/full/server-key.pem"     # Optional, default: "", path of certificate on local FS
#    clientCertPath: "example/boot/full/server.pem"        # Optional, default: "", path of certificate on local FS
#config:
#  - name: rk-main                                         # Required
#    path: "example/boot/full/config.yaml"                 # Required
#    locale: "*::*::*::*"                                  # Required, default: ""
#    description: "Description of entry"                   # Optional
grpc:
  - name: greeter                                          # Required
    enabled: true                                          # Required
    port: 8080                                             # Required
#    description: "greeter server"                         # Optional, default: ""
#    enableReflection: true                                # Optional, default: false
#    enableRkGwOption: true                                # Optional, default: false
#    gwOption:                                             # Optional, default: nil
#      marshal:                                            # Optional, default: nil
#        multiline: false                                  # Optional, default: false
#        emitUnpopulated: false                            # Optional, default: false
#        indent: ""                                        # Optional, default: false
#        allowPartial: false                               # Optional, default: false
#        useProtoNames: false                              # Optional, default: false
#        useEnumNumbers: false                             # Optional, default: false
#      unmarshal:                                          # Optional, default: nil
#        allowPartial: false                               # Optional, default: false
#        discardUnknown: false                             # Optional, default: false
#    noRecvMsgSizeLimit: true                              # Optional, default: false
#    gwMappingFilePaths: []                                # Optional
#    certEntry: "local-cert"                               # Optional, default: "", reference of cert entry declared above
#    sw:
#      enabled: true                                       # Optional, default: false
#      path: "sw"                                          # Optional, default: "sw"
#      jsonPath: ""                                        # Optional
#      headers: ["sw:rk"]                                  # Optional, default: []
#    commonService:
#      enabled: true                                       # Optional, default: false
#    static:
#      enabled: true                                       # Optional, default: false
#      path: "/rk/v1/static"                               # Optional, default: /rk/v1/static
#      sourceType: local                                   # Required, options: pkger, local
#      sourcePath: "."                                     # Required, full path of source directory
#    tv:
#      enabled:  true                                      # Optional, default: false
#    prom:
#      enabled: true                                       # Optional, default: false
#      path: ""                                            # Optional, default: "metrics"
#      pusher:
#        enabled: false                                    # Optional, default: false
#        jobName: "greeter-pusher"                         # Required
#        remoteAddress: "localhost:9091"                   # Required
#        basicAuth: "user:pass"                            # Optional, default: ""
#        intervalMs: 10000                                 # Optional, default: 1000
#        cert:                                             # Optional
#          ref: "local-test"                               # Optional, default: "", reference of cert entry declared above
#    logger:
#      zapLogger: zap-logger                               # Optional, default: logger of STDOUT, reference of logger entry declared above
#      eventLogger: event-logger                           # Optional, default: logger of STDOUT, reference of logger entry declared above
#    interceptors:
#      loggingZap:
#        enabled: true                                     # Optional, default: false
#        zapLoggerEncoding: "json"                         # Optional, default: "console"
#        zapLoggerOutputPaths: ["logs/app.log"]            # Optional, default: ["stdout"]
#        eventLoggerEncoding: "json"                       # Optional, default: "console"
#        eventLoggerOutputPaths: ["logs/event.log"]        # Optional, default: ["stdout"]
#      metricsProm:
#        enabled: true                                     # Optional, default: false
#      auth:
#        enabled: true                                     # Optional, default: false
#        basic:
#          - "user:pass"                                   # Optional, default: []
#        ignorePrefix:
#          - "/rk/v1"                                      # Optional, default: []
#        apiKey:
#          - "keys"                                        # Optional, default: []
#      meta:
#        enabled: true                                     # Optional, default: false
#        prefix: "rk"                                      # Optional, default: "rk"
#      tracingTelemetry:
#        enabled: true                                     # Optional, default: false
#        exporter:                                         # Optional, default will create a stdout exporter
#          file:
#            enabled: true                                 # Optional, default: false
#            outputPath: "logs/trace.log"                  # Optional, default: stdout
#          jaeger:
#            agent:
#              enabled: false                              # Optional, default: false
#              host: ""                                    # Optional, default: localhost
#              port: 0                                     # Optional, default: 6831
#            collector:
#              enabled: true                               # Optional, default: false
#              endpoint: ""                                # Optional, default: http://localhost:14268/api/traces
#              username: ""                                # Optional, default: ""
#              password: ""                                # Optional, default: ""
#      rateLimit:
#        enabled: false                                    # Optional, default: false
#        algorithm: "leakyBucket"                          # Optional, default: "tokenBucket"
#        reqPerSec: 100                                    # Optional, default: 1000000
#        paths:
#          - path: "/rk.api.v1.RkCommonService/Healthy"    # Optional, default: ""
#            reqPerSec: 0                                  # Optional, default: 1000000
#      timeout:
#        enabled: false                                    # Optional, default: false
#        timeoutMs: 5000                                   # Optional, default: 5000
#        paths:
#          - path: "/rk.api.v1.RkCommonService/Healthy"    # Optional, default: ""
#            timeoutMs: 1000                               # Optional, default: 5000
#      jwt:
#        enabled: true                                     # Optional, default: false
#        signingKey: "my-secret"                           # Required
#        ignorePrefix:                                     # Optional, default: []
#          - "/rk/v1/tv"
#          - "/sw"
#          - "/rk/v1/assets"
#        signingKeys:                                      # Optional
#          - "key:value"
#        signingAlgo: ""                                   # Optional, default: "HS256"
#        tokenLookup: "header:<name>"                      # Optional, default: "header:Authorization"
#        authScheme: "Bearer"                              # Optional, default: "Bearer"
#      csrf:
#        enabled: true
#        tokenLength: 32                                   # Optional, default: 32
#        tokenLookup: "header:X-CSRF-Token"                # Optional, default: "header:X-CSRF-Token"
#        cookieName: "_csrf"                               # Optional, default: _csrf
#        cookieDomain: ""                                  # Optional, default: ""
#        cookiePath: ""                                    # Optional, default: ""
#        cookieMaxAge: 86400                               # Optional, default: 86400
#        cookieHttpOnly: false                             # Optional, default: false
#        cookieSameSite: "default"                         # Optional, default: "default", options: lax, strict, none, default
#        ignorePrefix: []                                  # Optional, default: []
#      cors:
#        enabled: true                                     # Optional, default: false
#        allowOrigins:
#          - "http://localhost:*"                          # Optional, default: *
#        allowCredentials: false                           # Optional, default: false
#        allowHeaders: []                                  # Optional, default: []
#        allowMethods: []                                  # Optional, default: []
#        exposeHeaders: []                                 # Optional, default: []
#        maxAge: 0                                         # Optional, default: 0
```

## Development Status: Stable

Released under the [Apache 2.0 License](../LICENSE).