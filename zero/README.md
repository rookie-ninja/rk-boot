# rk-boot/zero

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
    - [5.1 GoZero server](#51-gozero-server)
    - [5.2 Swagger & TV & Prometheus Client](#52-swagger--tv--prometheus-client)
    - [5.3 Logging](#53-logging)
    - [5.4 Meta](#54-meta)
    - [5.5 Send request](#55-send-request)
    - [5.6 RPC logs](#56-rpc-logs)
- [YAML Options](#yaml-options)
  - [go-zero](#go-zero)
  - [CommonService](#commonservice)
  - [Swagger](#swagger)
  - [Prometheus Client](#prometheus-client)
  - [TV](#tv)
  - [Middlewares](#middlewares)
    - [Log](#log)
    - [Metrics (prometheus)](#metrics-prometheus)
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
- [Development Status: Testing](#development-status-testing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Overview
Middlewares & bootstrapper designed for [zeromicro/go-zero](https://github.com/zeromicro/go-zero) web framework.

With rk-boot/zero, user can start [zeromicro/go-zero](https://github.com/zeromicro/go-zero) server via boot.yaml file easily.

**Dependency**: [rk-zero](https://github.com/rookie-ninja/rk-zero)

## Architecture
![image](https://github.com/rookie-ninja/rk-zero/blob/master/docs/img/zero-arch.png?raw=true)

## Features
### Supported bootstrap
| Bootstrap | Description |
| --- | --- |
| YAML based | Start [zeromicro/go-zero](https://github.com/zeromicro/go-zero) microservice from YAML |
| Code based | Start [zeromicro/go-zero](https://github.com/zeromicro/go-zero) microservice from code |

### Supported instances
All instances could be configured via YAML or Code.

**User can enable anyone of those as needed! No mandatory binding!**

| Instance | Description |
| --- | --- |
| rest.Server | Compatible with original [zeromicro/go-zero](https://github.com/zeromicro/go-zero) service functionalities |
| Config | Configure [spf13/viper](https://github.com/spf13/viper) as config instance and reference it from YAML |
| Logger | Configure [uber-go/zap](https://github.com/uber-go/zap) logger configuration and reference it from YAML |
| EventLogger | Configure logging of RPC with [rk-query](https://github.com/rookie-ninja/rk-query) and reference it from YAML |
| Credential | Fetch credentials from remote datastore like ETCD. |
| Cert | Fetch TLS/SSL certificates from remote datastore like ETCD and start microservice. |
| Prometheus | Start prometheus client at client side and push metrics to pushgateway as needed. |
| Swagger | Builtin swagger UI handler. |
| CommonService | List of common APIs. |
| TV | A Web UI shows microservice and environment information. |

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
`go get github.com/rookie-ninja/rk-boot/zero`

## Quick Start
In the bellow example, we will start microservice with bellow functionality and middlewares enabled via YAML.

- [zeromicro/go-zero](https://github.com/zeromicro/go-zero) server
- Swagger UI
- CommonService
- TV
- Prometheus Metrics (middleware)
- Logging (middleware)
- Meta (middleware)

### 1.Create boot.yaml
Since we are going to generate swagger config files with [swag](https://github.com/swaggo/swag)，the generated config files will be in docs/ folder by default.
[zero.sw.jsonPath: "docs"] needs to be specified in order to make server read swagger config file for user defined API.

```yaml
---
zero:
  - name: greeter                     # Required
    port: 8080                        # Required
    enabled: true                     # Required
    tv:
      enabled: true                   # Optional, default: false
    prom:
      enabled: true                   # Optional, default: false
    sw:
      enabled: true                   # Optional, default: false
      jsonPath: "docs"                # Optional, default: ""
    commonService:
      enabled: true                   # Optional, default: false
    interceptors:
      loggingZap:
        enabled: true                 # Optional, default: false
      metricsProm:
        enabled: true                 # Optional, default: false
      meta:
        enabled: true                 # Optional, default: false
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
	"encoding/json"
	"fmt"
	"github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-boot/zero"
	"github.com/tal-tech/go-zero/rest"
	"net/http"
)

// @title RK Swagger for Zero
// @version 1.0
// @description This is a greeter service with rk-boot.

// Application entrance.
func main() {
	// Create a new boot instance.
	boot := rkboot.NewBoot()

	// Register handler
	zeroEntry := rkbootzero.GetZeroEntry("greeter")
	zeroEntry.Server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/v1/greeter",
		Handler: Greeter,
	})

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
func Greeter(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(&GreeterResponse{
		Message: fmt.Sprintf("Hello %s!", request.URL.Query().Get("name")),
	})

	writer.Write(bytes)
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
#### 5.1 GoZero server
Try to test GoZero Service with [curl](https://curl.se/)

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
2022-01-09T00:46:59.941+0800    INFO    boot/zero_entry.go:1168 Bootstrap zeroEntry     {"eventId": "e333316b-9ed9-406f-8f91-c63e60c80b04", "entryName": "greeter"}
------------------------------------------------------------------------
endTime=2022-01-09T00:46:59.94303+08:00
startTime=2022-01-09T00:46:59.941778+08:00
elapsedNano=1252080
timezone=CST
ids={"eventId":"e333316b-9ed9-406f-8f91-c63e60c80b04"}
app={"appName":"rk","appVersion":"","entryName":"greeter","entryType":"ZeroEntry"}
env={"arch":"amd64","az":"*","domain":"*","hostname":"lark.local","localIP":"10.8.0.2","os":"darwin","realm":"*","region":"*"}
payloads={"commonServiceEnabled":true,"commonServicePathPrefix":"/rk/v1/","promEnabled":true,"promPath":"/metrics","promPort":8080,"swEnabled":true,"swPath":"/sw/","tvEnabled":true,"tvPath":"/rk/v1/tv/","zeroPort":8080}
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
By default, we will send back some metadata to client including gateway with headers.

```shell script
$ curl -vs localhost:8080/rk/v1/healthy
...
< HTTP/1.1 200 OK
< Content-Type: application/json
< X-Request-Id: f068670c-56dc-42ed-9368-14961acd8433
< X-Rk-App-Name: rk
< X-Rk-App-Unix-Time: 2021-12-29T01:50:35.791623+08:00
< X-Rk-App-Version: 
< X-Rk-Received-Time: 2021-12-29T01:50:35.791623+08:00
< X-Trace-Id: 3a84ce72b53ce43468a5e81e8dceadcb
< Date: Tue, 28 Dec 2021 17:50:35 GMT
...
```

#### 5.5 Send request
We registered /v1/greeter API in [zeromicro/go-zero](https://github.com/zeromicro/go-zero) server and let's validate it!

```shell script
$ curl "localhost:8080/v1/greeter?name=rk-dev"
{"Message":"Hello rk-dev!"}
```

#### 5.6 RPC logs
Bellow logs would be printed in stdout.

```
------------------------------------------------------------------------
endTime=2022-01-09T00:49:08.693105+08:00
startTime=2022-01-09T00:49:08.693026+08:00
elapsedNano=79193
timezone=CST
ids={"eventId":"112de027-2e2c-4f0c-8f9b-9874567aa18f","requestId":"112de027-2e2c-4f0c-8f9b-9874567aa18f","traceId":"1ca745c23eb7896795d88d468385bc47"}
app={"appName":"rk","appVersion":"","entryName":"greeter","entryType":"ZeroEntry"}
env={"arch":"amd64","az":"*","domain":"*","hostname":"lark.local","localIP":"10.8.0.2","os":"darwin","realm":"*","region":"*"}
payloads={"apiMethod":"GET","apiPath":"/v1/greeter","apiProtocol":"HTTP/1.1","apiQuery":"name=rk-dev","userAgent":"curl/7.64.1"}
error={}
counters={}
pairs={}
timing={}
remoteAddr=localhost:64704
operation=/v1/greeter
resCode=200
eventStatus=Ended
EOE
```

## YAML Options
User can start multiple [zeromicro/go-zero](https://github.com/zeromicro/go-zero) instances at the same time. Please make sure use different port and name.

### go-zero
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| zero.name | The name of zero server | string | N/A |
| zero.port | The port of zero server | integer | nil, server won't start |
| zero.enabled | Enable zero entry or not | bool | false |
| zero.description | Description of zero entry. | string | "" |
| zero.cert.ref | Reference of cert entry declared in [cert entry](https://github.com/rookie-ninja/rk-entry#certentry) | string | "" |
| zero.logger.zapLogger.ref | Reference of zapLoggerEntry declared in [zapLoggerEntry](https://github.com/rookie-ninja/rk-entry#zaploggerentry) | string | "" |
| zero.logger.eventLogger.ref | Reference of eventLoggerEntry declared in [eventLoggerEntry](https://github.com/rookie-ninja/rk-entry#eventloggerentry) | string | "" |

```yaml
zero:
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
| /rk/v1/sys | Get OS stat. |
| /rk/v1/tv | Get HTML page of /tv. |

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| zero.commonService.enabled | Enable embedded common service | boolean | false |

```yaml
zero:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    commonService:
      enabled: true                                        # Optional, default: false
```

### Swagger
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| zero.sw.enabled | Enable swagger service over zero server | boolean | false |
| zero.sw.path | The path access swagger service from web | string | /sw |
| zero.sw.jsonPath | Where the swagger.json files are stored locally | string | "" |
| zero.sw.headers | Headers would be sent to caller as scheme of [key:value] | []string | [] |

```yaml
zero:
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
| zero.prom.enabled | Enable prometheus | boolean | false |
| zero.prom.path | Path of prometheus | string | /metrics |
| zero.prom.pusher.enabled | Enable prometheus pusher | bool | false |
| zero.prom.pusher.jobName | Job name would be attached as label while pushing to remote pushgateway | string | "" |
| zero.prom.pusher.remoteAddress | PushGateWay address, could be form of http://x.x.x.x or x.x.x.x | string | "" |
| zero.prom.pusher.intervalMs | Push interval in milliseconds | string | 1000 |
| zero.prom.pusher.basicAuth | Basic auth used to interact with remote pushgateway, form of [user:pass] | string | "" |
| zero.prom.pusher.cert.ref | Reference of rkentry.CertEntry | string | "" |

```yaml
zero:
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
| zero.tv.enabled | Enable RK TV | boolean | false |

```yaml
zero:
  - name: greeter                                          # Required
    port: 8080                                             # Required
    enabled: true                                          # Required
    tv:
      enabled: true                                        # Optional, default: false
```

### Middlewares
#### Log
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| zero.interceptors.loggingZap.enabled | Enable log interceptor | boolean | false |
| zero.interceptors.loggingZap.zapLoggerEncoding | json or console | string | console |
| zero.interceptors.loggingZap.zapLoggerOutputPaths | Output paths | []string | stdout |
| zero.interceptors.loggingZap.eventLoggerEncoding | json or console | string | console |
| zero.interceptors.loggingZap.eventLoggerOutputPaths | Output paths | []string | false |

```yaml
zero:
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
endTime=2021-11-01T23:31:01.706614+08:00
startTime=2021-11-01T23:31:01.706335+08:00
elapsedNano=278966
timezone=CST
ids={"eventId":"61cae46e-ea98-47b5-8a39-1090d015e09a","requestId":"61cae46e-ea98-47b5-8a39-1090d015e09a"}
app={"appName":"rk-zero","appVersion":"master-e4538d7","entryName":"greeter","entryType":"ZeroEntry"}
env={"arch":"amd64","az":"*","domain":"*","hostname":"lark.local","localIP":"192.168.1.104","os":"darwin","realm":"*","region":"*"}
payloads={"apiMethod":"GET","apiPath":"/rk/v1/healthy","apiProtocol":"HTTP/1.1","apiQuery":"","userAgent":"curl/7.64.1"}
error={}
counters={}
pairs={}
timing={}
remoteAddr=localhost:54376
operation=/rk/v1/healthy
resCode=200
eventStatus=Ended
EOE
```

#### Metrics (prometheus)
[zero.prom.enabled: true] is necessary since middleware needs to prometheus client instance in server.

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| zero.interceptors.metricsProm.enabled | Enable metrics interceptor | boolean | false |

```yaml
zero:
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
| zero.interceptors.auth.enabled | Enable auth interceptor | boolean | false |
| zero.interceptors.auth.basic | Basic auth credentials as scheme of <user:pass> | []string | [] |
| zero.interceptors.auth.apiKey | API key auth | []string | [] |
| zero.interceptors.auth.ignorePrefix | The paths of prefix that will be ignored by interceptor | []string | [] |

```yaml
zero:
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
| zero.interceptors.meta.enabled | Enable meta interceptor | boolean | false |
| zero.interceptors.meta.prefix | Header key was formed as X-<Prefix>-XXX | string | RK |

```yaml
zero:
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
| zero.interceptors.tracingTelemetry.enabled | Enable tracing interceptor | boolean | false |
| zero.interceptors.tracingTelemetry.exporter.file.enabled | Enable file exporter | boolean | RK |
| zero.interceptors.tracingTelemetry.exporter.file.outputPath | Export tracing info to files | string | stdout |
| zero.interceptors.tracingTelemetry.exporter.jaeger.agent.enabled | Export tracing info to jaeger agent | boolean | false |
| zero.interceptors.tracingTelemetry.exporter.jaeger.agent.host | As name described | string | localhost |
| zero.interceptors.tracingTelemetry.exporter.jaeger.agent.port | As name described | int | 6831 |
| zero.interceptors.tracingTelemetry.exporter.jaeger.collector.enabled | Export tracing info to jaeger collector | boolean | false |
| zero.interceptors.tracingTelemetry.exporter.jaeger.collector.endpoint | As name described | string | http://localhost:16368/api/trace |
| zero.interceptors.tracingTelemetry.exporter.jaeger.collector.username | As name described | string | "" |
| zero.interceptors.tracingTelemetry.exporter.jaeger.collector.password | As name described | string | "" |

```yaml
zero:
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
| zero.interceptors.rateLimit.enabled | Enable rate limit interceptor | boolean | false |
| zero.interceptors.rateLimit.algorithm | Provide algorithm, tokenBucket and leakyBucket are available options | string | tokenBucket |
| zero.interceptors.rateLimit.reqPerSec | Request per second globally | int | 0 |
| zero.interceptors.rateLimit.paths.path | Full path | string | "" |
| zero.interceptors.rateLimit.paths.reqPerSec | Request per second by full path | int | 0 |

```yaml
zero:
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
| zero.interceptors.timeout.enabled | Enable timeout interceptor | boolean | false |
| zero.interceptors.timeout.timeoutMs | Global timeout in milliseconds. | int | 5000 |
| zero.interceptors.timeout.paths.path | Full path | string | "" |
| zero.interceptors.timeout.paths.timeoutMs | Timeout in milliseconds by full path | int | 5000 |

```yaml
zero:
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
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| zero.interceptors.cors.enabled | Enable cors interceptor | boolean | false |
| zero.interceptors.cors.allowOrigins | Provide allowed origins with wildcard enabled. | []string | * |
| zero.interceptors.cors.allowMethods | Provide allowed methods returns as response header of OPTIONS request. | []string | All http methods |
| zero.interceptors.cors.allowHeaders | Provide allowed headers returns as response header of OPTIONS request. | []string | Headers from request |
| zero.interceptors.cors.allowCredentials | Returns as response header of OPTIONS request. | bool | false |
| zero.interceptors.cors.exposeHeaders | Provide exposed headers returns as response header of OPTIONS request. | []string | "" |
| zero.interceptors.cors.maxAge | Provide max age returns as response header of OPTIONS request. | int | 0 |

```yaml
---
zero:
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
| zero.interceptors.jwt.enabled | Enable JWT interceptor | boolean | false |
| zero.interceptors.jwt.signingKey | Required, Provide signing key. | string | "" |
| zero.interceptors.jwt.ignorePrefix | Provide ignoring path prefix. | []string | [] |
| zero.interceptors.jwt.signingKeys | Provide signing keys as scheme of <key>:<value>. | []string | [] |
| zero.interceptors.jwt.signingAlgo | Provide signing algorithm. | string | HS256 |
| zero.interceptors.jwt.tokenLookup | Provide token lookup scheme, please see bellow description. | string | "header:Authorization" |
| zero.interceptors.jwt.authScheme | Provide auth scheme. | string | Bearer |

```yaml
zero:
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
| zero.interceptors.secure.enabled | Enable secure interceptor | boolean | false |
| zero.interceptors.secure.xssProtection | X-XSS-Protection header value. | string | "1; mode=block" |
| zero.interceptors.secure.contentTypeNosniff | X-Content-Type-Options header value. | string | nosniff |
| zero.interceptors.secure.xFrameOptions | X-Frame-Options header value. | string | SAMEORIGIN |
| zero.interceptors.secure.hstsMaxAge | Strict-Transport-Security header value. | int | 0 |
| zero.interceptors.secure.hstsExcludeSubdomains | Excluding subdomains of HSTS. | bool | false |
| zero.interceptors.secure.hstsPreloadEnabled | Enabling HSTS preload. | bool | false |
| zero.interceptors.secure.contentSecurityPolicy | Content-Security-Policy header value. | string | "" |
| zero.interceptors.secure.cspReportOnly | Content-Security-Policy-Report-Only header value. | bool | false |
| zero.interceptors.secure.referrerPolicy | Referrer-Policy header value. | string | "" |
| zero.interceptors.secure.ignorePrefix | Ignoring path prefix. | []string | [] |

```yaml
zero:
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
| zero.interceptors.csrf.enabled | Enable csrf interceptor | boolean | false |
| zero.interceptors.csrf.tokenLength | Provide the length of the generated token. | int | 32 |
| zero.interceptors.csrf.tokenLookup | Provide csrf token lookup rules, please see code comments for details. | string | "header:X-CSRF-Token" |
| zero.interceptors.csrf.cookieName | Provide name of the CSRF cookie. This cookie will store CSRF token. | string | _csrf |
| zero.interceptors.csrf.cookieDomain | Domain of the CSRF cookie. | string | "" |
| zero.interceptors.csrf.cookiePath | Path of the CSRF cookie. | string | "" |
| zero.interceptors.csrf.cookieMaxAge | Provide max age (in seconds) of the CSRF cookie. | int | 86400 |
| zero.interceptors.csrf.cookieHttpOnly | Indicates if CSRF cookie is HTTP only. | bool | false |
| zero.interceptors.csrf.cookieSameSite | Indicates SameSite mode of the CSRF cookie. Options: lax, strict, none, default | string | default |
| zero.interceptors.csrf.ignorePrefix | Ignoring path prefix. | []string | [] |

```yaml
zero:
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
#    description: "Description of entry"                   # Optional
#    locale: "*::*::*::*"                                  # Optional, default: *::*::*::*
#    paths:                                                # Optional
#      - "example/boot/full/cred.yaml"
#cert:
#  - name: "local-cert"                                    # Required
#    provider: "localFs"                                   # Required, etcd, consul, localFs, remoteFs are supported options
#    description: "Description of entry"                   # Optional
#    locale: "*::*::*::*"                                  # Optional, default: *::*::*::*
#    serverCertPath: "example/boot/full/server.pem"        # Optional, default: "", path of certificate on local FS
#    serverKeyPath: "example/boot/full/server-key.pem"     # Optional, default: "", path of certificate on local FS
#    clientCertPath: "example/client.pem"                  # Optional, default: "", path of certificate on local FS
#    clientKeyPath: "example/client.pem"                   # Optional, default: "", path of certificate on local FS
#config:
#  - name: rk-main                                         # Required
#    path: "example/boot/full/config.yaml"                 # Required
#    locale: "*::*::*::*"                                  # Required, default: *::*::*::*
#    description: "Description of entry"                   # Optional
zero:
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

## Development Status: Testing

Released under the [Apache 2.0 License](../LICENSE).

