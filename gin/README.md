# rk-boot/gin
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
  - [1.Create boot.yaml](#1create-bootyaml)
  - [2.Create main.go](#2create-maingo)
  - [3.Generate swagger config file](#3generate-swagger-config-file)
  - [4.Start server](#4start-server)
  - [5.Validation](#5validation)
- [YAML Options](#yaml-options)
  - [Gin](#gin)
  - [CommonService](#commonservice)
  - [Swagger UI](#swagger-ui)
  - [Prometheus Client](#prometheus-client)
  - [TV](#tv)
  - [Static file handler](#static-file-handler)
  - [Middlewares](#middlewares)
  - [Full YAML](#full-yaml)
- [Development Status: Stable](#development-status-stable)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Overview
Middlewares & bootstrapper designed for [gin-gonic/gin](https://github.com/gin-gonic/gin) web framework, see [official docs](https://rkdev.info/docs/bootstrapper/user-guide/gin-golang/)

With rk-boot/gin, user can start [gin-gonic/gin](https://github.com/gin-gonic/gin) server via boot.yaml file easily.

**Dependency**: [rk-gin](https://github.com/rookie-ninja/rk-gin)

## Architecture
![image](https://raw.githubusercontent.com/rookie-ninja/rk-gin/master/docs/img/gin-arch.png)

## Features
### Supported bootstrap
| Bootstrap | Description |
| --- | --- |
| YAML based | Start [gin-gonic/gin](https://github.com/gin-gonic/gin) microservice from YAML |
| Code based | Start [gin-gonic/gin](https://github.com/gin-gonic/gin) microservice from code |

### Supported instances
All instances could be configured via YAML or Code.

**User can enable anyone of those as needed! No mandatory binding!**

| Instance | Description |
| --- | --- |
| gin.Router | Compatible with original [gin-gonic/gin](https://github.com/gin-gonic/gin) service functionalities |
| Config | Configure [spf13/viper](https://github.com/spf13/viper) as config instance and reference it from YAML |
| Logger | Configure [uber-go/zap](https://github.com/uber-go/zap) logger configuration and reference it from YAML |
| EventLogger | Configure logging of RPC with [rk-query](https://github.com/rookie-ninja/rk-query) and reference it from YAML |
| Credential | Fetch credentials from remote datastore like ETCD. |
| Cert | Fetch TLS/SSL certificates from remote datastore like ETCD and start microservice. |
| Prometheus | Start prometheus client at client side and push metrics to pushgateway as needed. |
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
| Meta | Send microservice metadata as header to client. |
| Auth | Support [Basic Auth] and [API Key] authorization types. |
| RateLimit | Limiting RPC rate globally or per path. |
| Timeout | Timing out request by configuration. |
| Gzip | Compress and Decompress message body based on request header with gzip format . |
| CORS | Server side CORS validation. |
| JWT | Server side JWT validation. |
| Secure | Server side secure validation. |
| CSRF | Server side csrf validation. |

## Installation
`go get github.com/rookie-ninja/rk-boot/gin`

## Quick Start
In the bellow example, we will start microservice with bellow functionality and middlewares enabled via YAML.

- [gin-gonic/gin](https://github.com/gin-gonic/gin) server
- Swagger UI
- CommonService
- TV
- Prometheus Metrics (middleware)
- Logging (middleware)
- Meta (middleware)

### 1.Create boot.yaml
Since we are going to generate swagger config files with [swag](https://github.com/swaggo/swag)，the generated config files will be in docs/ folder by default.
[gin.sw.jsonPath: "docs"] needs to be specified in order to make server read swagger config file for user defined API.

```yaml
---
gin:
  - name: greeter                     # Required
    port: 8080                        # Required
    enabled: true                     # Required
    tv:
      enabled: true                   # Optional, default: false
    prom:
      enabled: true                   # Optional, default: false
    sw:                               # Optional
      enabled: true                   # Optional, default: false
      jsonPath: "docs"                # Optional, default: ""
    commonService:                    # Optional
      enabled: true                   # Optional, default: false
    interceptors:
      loggingZap:
        enabled: true
      metricsProm:
        enabled: true
      meta:
        enabled: true
```

### 2.Create main.go
Since we are going to generate swagger config files with [swag](https://github.com/swaggo/swag)，some comments needs to be added as bellow.

```go
// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-boot/gin"
	"net/http"
)

// @title RK Swagger for Gin
// @version 1.0
// @description This is a greeter service with rk-boot.

// Application entrance.
func main() {
	// Create a new boot instance.
	boot := rkboot.NewBoot()

	// Register handler
	ginEntry := rkbootgin.GetGinEntry("greeter")
	ginEntry.Router.GET("/v1/greeter", Greeter)

	// Bootstrap
	boot.Bootstrap(context.Background())

	// Wait for shutdown sig
	boot.WaitForShutdownSig(context.Background())
}

// @Summary Greeter service
// @Id 1
// @version 1.0
// @produce application/json
// @Param name query string true "Input name"
// @Success 200 {object} GreeterResponse
// @Router /v1/greeter [get]
func Greeter(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, &GreeterResponse{
		Message: fmt.Sprintf("Hello %s!", ctx.Query("name")),
	})
}

// Response.
type GreeterResponse struct {
	Message string
}
```

### 3.Generate swagger config file
Files would be generated as bellow.

```
$ swag init

$ tree
.
├── boot.yaml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
└── main.go
```

### 4.Start server
```go
$ go run main.go
```

### 5.Validation
#### 5.1 Gin server
Try to test Gin Service with [curl](https://curl.se/)

```shell script
# Curl to common service
$ curl localhost:8080/rk/v1/healthy
{"healthy":true}
```

#### 5.2 Swagger & TV & Prometheus Client
- Swagger UI: [http://localhost:8080/sw](http://localhost:8080/sw)
- TV: [http://localhost:8080/rk/v1/tv](http://localhost:8080/rk/v1/tv)
- Prometheus Client: [http://localhost:8080/metrics](http://localhost:8080/metrics)

#### 5.3 Logging
By default, we enable zap logger and event logger with encoding type of [console]. Encoding type of [json] is also supported.

```shell script
2022-01-08T19:55:09.178+0800    INFO    boot/gin_entry.go:913   Bootstrap ginEntry      {"eventId": "61f1b5f8-f9d1-4161-8e5a-7b927b0de78e", "entryName": "greeter"}
------------------------------------------------------------------------
endTime=2022-01-08T19:55:09.179436+08:00
startTime=2022-01-08T19:55:09.178152+08:00
elapsedNano=1284874
timezone=CST
ids={"eventId":"61f1b5f8-f9d1-4161-8e5a-7b927b0de78e"}
app={"appName":"rk","appVersion":"","entryName":"greeter","entryType":"GinEntry"}
env={"arch":"amd64","az":"*","domain":"*","hostname":"lark.local","localIP":"10.8.0.2","os":"darwin","realm":"*","region":"*"}
payloads={"commonServiceEnabled":true,"commonServicePathPrefix":"/rk/v1/","ginPort":8080,"promEnabled":true,"promPath":"/metrics","promPort":8080,"swEnabled":true,"swPath":"/sw/","tvEnabled":true,"tvPath":"/rk/v1/tv/"}
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

#### 5.4 Meta
Please refer [documentation](https://rkdev.info/docs/bootstrapper/user-guide/gin-golang/basic/middleware-meta/) for details of configuration.

By default, we will send back some metadata to client including gateway with headers.

```shell script
$ curl -vs "localhost:8080/v1/greeter?name=rk-dev"
...
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< X-Request-Id: 65a06e34-a9d8-40f4-b144-5fe8aca715f6
< X-Rk-App-Name: rk
< X-Rk-App-Unix-Time: 2022-01-08T19:55:53.583409+08:00
< X-Rk-Received-Time: 2022-01-08T19:55:53.583409+08:00
< Date: Sat, 08 Jan 2022 11:55:53 GMT
...
{"Message":"Hello rk-dev!"}
```

#### 5.5 Send request
We registered /v1/greeter API in [gin-gonic/gin](https://github.com/gin-gonic/gin) server and let's validate it!

```shell script
$ curl "localhost:8080/v1/greeter?name=rk-dev"
{"Message":"Hello rk-dev!"}
```

#### 5.6 RPC logs
Bellow logs would be printed in stdout.

```
------------------------------------------------------------------------
endTime=2022-01-08T19:55:53.583461+08:00
startTime=2022-01-08T19:55:53.5834+08:00
elapsedNano=60919
timezone=CST
ids={"eventId":"65a06e34-a9d8-40f4-b144-5fe8aca715f6","requestId":"65a06e34-a9d8-40f4-b144-5fe8aca715f6"}
app={"appName":"rk","appVersion":"","entryName":"greeter","entryType":"GinEntry"}
env={"arch":"amd64","az":"*","domain":"*","hostname":"lark.local","localIP":"10.8.0.2","os":"darwin","realm":"*","region":"*"}
payloads={"apiMethod":"GET","apiPath":"/v1/greeter","apiProtocol":"HTTP/1.1","apiQuery":"name=rk-dev","userAgent":"curl/7.64.1"}
error={}
counters={}
pairs={}
timing={}
remoteAddr=localhost:57925
operation=/v1/greeter
resCode=200
eventStatus=Ended
EOE
```

#### 5.7 RPC prometheus metrics
Access [http://localhost:8080/metrics](http://localhost:8080/metrics)

## YAML Options
User can start multiple [gin-gonic/gin](https://github.com/gin-gonic/gin) instances at the same time. Please make sure use different port and name.

### Gin
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.name | Required, The name of gin server | string | N/A |
| gin.port | Required, The port of gin server | integer | nil, server won't start |
| gin.enabled | Optional, Enable Gin entry or not | bool | false |
| gin.description | Optional, Description of gin entry. | string | "" |
| gin.cert.ref | Optional, Reference of cert entry declared in [cert entry](https://github.com/rookie-ninja/rk-entry#certentry) | string | "" |
| gin.logger.zapLogger.ref | Optional, Reference of zapLoggerEntry declared in [zapLoggerEntry](https://github.com/rookie-ninja/rk-entry#zaploggerentry) | string | "" |
| gin.logger.eventLogger.ref | Optional, Reference of eventLoggerEntry declared in [eventLoggerEntry](https://github.com/rookie-ninja/rk-entry#eventloggerentry) | string | "" |

```yaml
gin:
  - name: greeter                                         # Required
    port: 8080                                            # Required
    enabled: true                                         # Required
    description: "greeter server"                         # Optional, default: ""
    cert:
      ref: "local-cert"                                   # Optional, default: "", reference of cert entry declared above
    logger:
      zapLogger:
        ref: zap-logger                                   # Optional, default: logger of STDOUT, reference of logger entry declared above
      eventLogger:
        ref: event-logger                                 # Optional, default: logger of STDOUT, reference of logger entry declared above
```

### CommonService
| Path | Description |
| ---- | ---- |
| /rk/v1/apis | List APIs in current GinEntry. |
| /rk/v1/certs | List CertEntry. |
| /rk/v1/configs | List ConfigEntry. |
| /rk/v1/deps | List dependencies related application, entire contents of go.mod file would be returned. |
| /rk/v1/entries | List all Entries. |
| /rk/v1/gc | Trigger GC |
| /rk/v1/healthy | Get application healthy status. |
| /rk/v1/info | Get application and process info. |
| /rk/v1/license | Get license related application, entire contents of LICENSE file would be returned. |
| /rk/v1/logs | List logger related entries. |
| /rk/v1/git | Get git information. |
| /rk/v1/readme | Get contents of README file. |
| /rk/v1/req | List prometheus metrics of requests. |
| /rk/v1/sys | Get OS stat. |
| /rk/v1/tv | Get HTML page of /tv. |

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.commonService.enabled | Optional, Enable embedded common service | boolean | false |

```yaml
gin:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    commonService:
      enabled: true                                        # Optional, default: false
```

### Swagger UI
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.sw.enabled | Optional, Enable swagger service over gin server | boolean | false |
| gin.sw.path | Optional, The path access swagger service from web | string | /sw |
| gin.sw.jsonPath | Optional, Where the swagger.json files are stored locally | string | "" |
| gin.sw.headers | Optional, Headers would be sent to caller as scheme of [key:value] | []string | [] |

```yaml
gin:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    sw:
      enabled: true                                        # Optional, default: false
      path: "sw"                                           # Optional, default: "sw"
      jsonPath: ""                                         # Optional, default: ""
      headers: ["sw:rk"]                                   # Optional, default: []
```

### Prometheus Client
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.prom.enabled | Optional, Enable prometheus | boolean | false |
| gin.prom.path | Optional, Path of prometheus | string | /metrics |
| gin.prom.pusher.enabled | Optional, Enable prometheus pusher | bool | false |
| gin.prom.pusher.jobName | Optional, Job name would be attached as label while pushing to remote pushgateway | string | "" |
| gin.prom.pusher.remoteAddress | Optional, PushGateWay address, could be form of http://x.x.x.x or x.x.x.x | string | "" |
| gin.prom.pusher.intervalMs | Optional, Push interval in milliseconds | string | 1000 |
| gin.prom.pusher.basicAuth | Optional, Basic auth used to interact with remote pushgateway, form of [user:pass] | string | "" |
| gin.prom.pusher.cert.ref | Optional, Reference of rkentry.CertEntry | string | "" |

```yaml
gin:
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

### TV
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.tv.enabled | Optional, Enable RK TV | boolean | false |

```yaml
gin:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    tv:
      enabled: true                                        # Optional, default: false
```

### Static file handler
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.static.enabled | Optional, Enable static file handler | boolean | false |
| gin.static.path | Optional, path of static file handler | string | /rk/v1/static |
| gin.static.sourceType | Required, local and pkger supported | string | "" |
| gin.static.sourcePath | Required, full path of source directory | string | "" |

- About [pkger](https://github.com/markbates/pkger)
User can use pkger command line tool to embed static files into .go files.

Please use sourcePath like: github.com/rookie-ninja/rk-gin:/boot/assets

```yaml
gin:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    static:
      enabled: true                                        # Optional, default: false
      path: "/rk/v1/static"                                # Optional, default: /rk/v1/static
      sourceType: local                                    # Required, options: pkger, local
      sourcePath: "."                                      # Required, full path of source directory
```

### Middlewares
#### Log
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.interceptors.loggingZap.enabled | Enable log interceptor | boolean | false |
| gin.interceptors.loggingZap.zapLoggerEncoding | json or console | string | console |
| gin.interceptors.loggingZap.zapLoggerOutputPaths | Output paths | []string | stdout |
| gin.interceptors.loggingZap.eventLoggerEncoding | json or console | string | console |
| gin.interceptors.loggingZap.eventLoggerOutputPaths | Output paths | []string | false |

```yaml
gin:
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
endTime=2021-06-25T01:30:45.144023+08:00
startTime=2021-06-25T01:30:45.143767+08:00
elapsedNano=255948
timezone=CST
ids={"eventId":"3332e575-43d8-4bfe-84dd-45b5fc5fb104","requestId":"3332e575-43d8-4bfe-84dd-45b5fc5fb104","traceId":"65b9aa7a9705268bba492fdf4a0e5652"}
app={"appName":"rk-gin","appVersion":"master-xxx","entryName":"greeter","entryType":"GinEntry"}
env={"arch":"amd64","az":"*","domain":"*","hostname":"lark.local","localIP":"10.8.0.2","os":"darwin","realm":"*","region":"*"}
payloads={"apiMethod":"GET","apiPath":"/rk/v1/healthy","apiProtocol":"HTTP/1.1","apiQuery":"","userAgent":"curl/7.64.1"}
error={}
counters={}
pairs={}
timing={}
remoteAddr=localhost:60718
operation=/rk/v1/healthy
resCode=200
eventStatus=Ended
EOE
```

#### Metrics (prometheus)
[gin.prom.enabled: true] is necessary since middleware needs to prometheus client instance in server.

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.interceptors.metricsProm.enabled | Enable metrics interceptor | boolean | false |

```yaml
gin:
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
| gin.interceptors.auth.enabled | Enable auth interceptor | boolean | false |
| gin.interceptors.auth.basic | Basic auth credentials as scheme of <user:pass> | []string | [] |
| gin.interceptors.auth.apiKey | API key auth | []string | [] |
| gin.interceptors.auth.ignorePrefix | The paths of prefix that will be ignored by interceptor | []string | [] |

```yaml
gin:
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
Send application metadata as header to client.

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.interceptors.meta.enabled | Enable meta interceptor | boolean | false |
| gin.interceptors.meta.prefix | Header key was formed as X-<Prefix>-XXX | string | RK |

```yaml
gin:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    interceptors:
      meta:
        enabled: true                                      # Optional, default: false
        prefix: "rk"                                       # Optional, default: "rk"
```

#### Tracing (open telemetry)
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.interceptors.tracingTelemetry.enabled | Enable tracing interceptor | boolean | false |
| gin.interceptors.tracingTelemetry.exporter.file.enabled | Enable file exporter | boolean | RK |
| gin.interceptors.tracingTelemetry.exporter.file.outputPath | Export tracing info to files | string | stdout |
| gin.interceptors.tracingTelemetry.exporter.jaeger.agent.enabled | Export tracing info to jaeger agent | boolean | false |
| gin.interceptors.tracingTelemetry.exporter.jaeger.agent.host | As name described | string | localhost |
| gin.interceptors.tracingTelemetry.exporter.jaeger.agent.port | As name described | int | 6831 |
| gin.interceptors.tracingTelemetry.exporter.jaeger.collector.enabled | Export tracing info to jaeger collector | boolean | false |
| gin.interceptors.tracingTelemetry.exporter.jaeger.collector.endpoint | As name described | string | http://localhost:16368/api/trace |
| gin.interceptors.tracingTelemetry.exporter.jaeger.collector.username | As name described | string | "" |
| gin.interceptors.tracingTelemetry.exporter.jaeger.collector.password | As name described | string | "" |

```yaml
gin:
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
| gin.interceptors.rateLimit.enabled | Enable rate limit interceptor | boolean | false |
| gin.interceptors.rateLimit.algorithm | Provide algorithm, tokenBucket and leakyBucket are available options | string | tokenBucket |
| gin.interceptors.rateLimit.reqPerSec | Request per second globally | int | 0 |
| gin.interceptors.rateLimit.paths.path | Full path | string | "" |
| gin.interceptors.rateLimit.paths.reqPerSec | Request per second by full path | int | 0 |

```yaml
gin:
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
| gin.interceptors.timeout.enabled | Enable timeout interceptor | boolean | false |
| gin.interceptors.timeout.timeoutMs | Global timeout in milliseconds. | int | 5000 |
| gin.interceptors.timeout.paths.path | Full path | string | "" |
| gin.interceptors.timeout.paths.timeoutMs | Timeout in milliseconds by full path | int | 5000 |

```yaml
gin:
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

#### Gzip
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.interceptors.gzip.enabled | Enable gzip interceptor | boolean | false |
| gin.interceptors.gzip.level | Provide level of compression, options are noCompression, bestSpeed, bestCompression, defaultCompression, huffmanOnly. | string | defaultCompression |

```yaml
gin:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    interceptors:
      gzip:
        enabled: true                                      # Optional, default: false
        level: defaultCompression                          # Optional, options: [noCompression, bestSpeed， bestCompression, defaultCompression, huffmanOnly]
```

#### CORS
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.interceptors.cors.enabled | Enable cors interceptor | boolean | false |
| gin.interceptors.cors.allowOrigins | Provide allowed origins with wildcard enabled. | []string | * |
| gin.interceptors.cors.allowMethods | Provide allowed methods returns as response header of OPTIONS request. | []string | All http methods |
| gin.interceptors.cors.allowHeaders | Provide allowed headers returns as response header of OPTIONS request. | []string | Headers from request |
| gin.interceptors.cors.allowCredentials | Returns as response header of OPTIONS request. | bool | false |
| gin.interceptors.cors.exposeHeaders | Provide exposed headers returns as response header of OPTIONS request. | []string | "" |
| gin.interceptors.cors.maxAge | Provide max age returns as response header of OPTIONS request. | int | 0 |

```yaml
---
gin:
  - name: greeter                     # Required
    port: 8080                        # Required
    enabled: true                     # Required
    interceptors:
      cors:
        enabled: true                 # Optional, default: false
        allowOrigins:
          - "http://localhost:*"      # Optional, default: *
```

#### JWT
> rk-gin using github.com/golang-jwt/jwt/v4, please beware of version compatibility.

In order to make swagger UI and RK tv work under JWT without JWT token, we need to ignore prefixes of paths as bellow.

```yaml
jwt:
  ...
  ignorePrefix:
   - "/rk/v1/tv"
   - "/sw"
   - "/rk/v1/assets"
```

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.interceptors.jwt.enabled | Enable JWT interceptor | boolean | false |
| gin.interceptors.jwt.signingKey | Required, Provide signing key. | string | "" |
| gin.interceptors.jwt.ignorePrefix | Provide ignoring path prefix. | []string | [] |
| gin.interceptors.jwt.signingKeys | Provide signing keys as scheme of <key>:<value>. | []string | [] |
| gin.interceptors.jwt.signingAlgo | Provide signing algorithm. | string | HS256 |
| gin.interceptors.jwt.tokenLookup | Provide token lookup scheme, please see bellow description. | string | "header:Authorization" |
| gin.interceptors.jwt.authScheme | Provide auth scheme. | string | Bearer |

```yaml
gin:
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
// - "query:<name>"
// - "param:<name>"
// - "cookie:<name>"
// - "form:<name>"
// Multiply sources example:
// - "header: Authorization,cookie: myowncookie"
```

#### Secure
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.interceptors.secure.enabled | Enable secure interceptor | boolean | false |
| gin.interceptors.secure.xssProtection | X-XSS-Protection header value. | string | "1; mode=block" |
| gin.interceptors.secure.contentTypeNosniff | X-Content-Type-Options header value. | string | nosniff |
| gin.interceptors.secure.xFrameOptions | X-Frame-Options header value. | string | SAMEORIGIN |
| gin.interceptors.secure.hstsMaxAge | Strict-Transport-Security header value. | int | 0 |
| gin.interceptors.secure.hstsExcludeSubdomains | Excluding subdomains of HSTS. | bool | false |
| gin.interceptors.secure.hstsPreloadEnabled | Enabling HSTS preload. | bool | false |
| gin.interceptors.secure.contentSecurityPolicy | Content-Security-Policy header value. | string | "" |
| gin.interceptors.secure.cspReportOnly | Content-Security-Policy-Report-Only header value. | bool | false |
| gin.interceptors.secure.referrerPolicy | Referrer-Policy header value. | string | "" |
| gin.interceptors.secure.ignorePrefix | Ignoring path prefix. | []string | [] |

```yaml
gin:
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
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| gin.interceptors.csrf.enabled | Enable csrf interceptor | boolean | false |
| gin.interceptors.csrf.tokenLength | Provide the length of the generated token. | int | 32 |
| gin.interceptors.csrf.tokenLookup | Provide csrf token lookup rules, please see code comments for details. | string | "header:X-CSRF-Token" |
| gin.interceptors.csrf.cookieName | Provide name of the CSRF cookie. This cookie will store CSRF token. | string | _csrf |
| gin.interceptors.csrf.cookieDomain | Domain of the CSRF cookie. | string | "" |
| gin.interceptors.csrf.cookiePath | Path of the CSRF cookie. | string | "" |
| gin.interceptors.csrf.cookieMaxAge | Provide max age (in seconds) of the CSRF cookie. | int | 86400 |
| gin.interceptors.csrf.cookieHttpOnly | Indicates if CSRF cookie is HTTP only. | bool | false |
| gin.interceptors.csrf.cookieSameSite | Indicates SameSite mode of the CSRF cookie. Options: lax, strict, none, default | string | default |
| gin.interceptors.csrf.ignorePrefix | Ignoring path prefix. | []string | [] |

```yaml
gin:
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
#eventLogger:
#  - name: event-logger                                    # Required
#    description: "Description of entry"                   # Optional
#cred:
#  - name: "local-cred"                                    # Required
#    provider: "localFs"                                   # Required, etcd, consul, localFs, remoteFs are supported options
#    locale: "*::*::*::*"                                  # Required, default: *::*::*::*
#    description: "Description of entry"                   # Optional
#    paths:                                                # Optional
#      - "example/boot/full/cred.yaml"
#cert:
#  - name: "local-cert"                                    # Required
#    provider: "localFs"                                   # Required, etcd, consul, localFs, remoteFs are supported options
#    locale: "*::*::*::*"                                  # Required, default: *::*::*::*
#    description: "Description of entry"                   # Optional
#    serverCertPath: "example/boot/full/server.pem"        # Optional, default: "", path of certificate on local FS
#    serverKeyPath: "example/boot/full/server-key.pem"     # Optional, default: "", path of certificate on local FS
#    clientCertPath: "example/client.pem"                  # Optional, default: "", path of certificate on local FS
#    clientKeyPath: "example/client.pem"                   # Optional, default: "", path of certificate on local FS
#config:
#  - name: rk-main                                         # Required
#    path: "example/boot/full/config.yaml"                 # Required
#    locale: "*::*::*::*"                                  # Required, default: *::*::*::*
#    description: "Description of entry"                   # Optional
gin:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
#    description: "greeter server"                         # Optional, default: ""
#    cert:
#      ref: "local-cert"                                   # Optional, default: "", reference of cert entry declared above
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
#      zapLogger:
#        ref: zap-logger                                   # Optional, default: logger of STDOUT, reference of logger entry declared above
#      eventLogger:
#        ref: event-logger                                 # Optional, default: logger of STDOUT, reference of logger entry declared above
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
#          - path: "/rk/v1/healthy"                        # Optional, default: ""
#            reqPerSec: 0                                  # Optional, default: 1000000
#      timeout:
#        enabled: false                                    # Optional, default: false
#        timeoutMs: 5000                                   # Optional, default: 5000
#        paths:
#          - path: "/rk/v1/healthy"                        # Optional, default: ""
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
#      secure:
#        enabled: true                                     # Optional, default: false
#        xssProtection: ""                                 # Optional, default: "1; mode=block"
#        contentTypeNosniff: ""                            # Optional, default: nosniff
#        xFrameOptions: ""                                 # Optional, default: SAMEORIGIN
#        hstsMaxAge: 0                                     # Optional, default: 0
#        hstsExcludeSubdomains: false                      # Optional, default: false
#        hstsPreloadEnabled: false                         # Optional, default: false
#        contentSecurityPolicy: ""                         # Optional, default: ""
#        cspReportOnly: false                              # Optional, default: false
#        referrerPolicy: ""                                # Optional, default: ""
#        ignorePrefix: []                                  # Optional, default: []
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
#      gzip:
#        enabled: true
#        level: bestSpeed                                  # Optional, options: [noCompression, bestSpeed， bestCompression, defaultCompression, huffmanOnly]
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

