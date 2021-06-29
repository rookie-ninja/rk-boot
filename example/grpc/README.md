# rk-grpc
Interceptor & bootstrapper designed for grpc. Currently, supports bellow functionalities.

| Name | Description |
| ---- | ---- |
| Start with YAML | Start service with YAML config. |
| Start with code | Start service from code. |
| GRPC Service | GRPC service defined with protocol buffer. |
| GRPC Gateway Service | GRPC Gateway service with new port. |
| Swagger Service | Swagger UI with same port as GRPC Gateway. |
| Common Service | List of common API available on GRPC, GRPC Gateway and swagger. |
| TV Service | A Web UI shows application and environment information. |
| Metrics interceptor | Collect RPC metrics and export as prometheus client with same port of GRPC gateway. |
| Log interceptor | Log every RPC requests as event with rk-query. |
| Trace interceptor | Collect RPC trace and export it to stdout, file or jaeger. |
| Panic interceptor | Recover from panic for RPC requests and log it. |
| Meta interceptor | Send application metadata as header to client and GRPC Gateway. |
| Auth interceptor | Support [Basic Auth], [Bearer Token] and [API Key] authrization types. |

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Installation](#installation)
- [YAML Config](#yaml-config)
  - [GRPC Service](#grpc-service)
  - [Common Service](#common-service)
  - [GRPC Gateway Service](#grpc-gateway-service)
    - [Swagger Service](#swagger-service)
    - [Prom Client](#prom-client)
    - [TV Service](#tv-service)
  - [Interceptors](#interceptors)
    - [Log](#log)
    - [Metrics](#metrics)
    - [Auth](#auth)
    - [Meta](#meta)
    - [Tracing](#tracing)
  - [Development Status: Stable](#development-status-stable)
  - [Appendix](#appendix)
  - [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation
`go get -u github.com/rookie-ninja/rk-grpc`

## YAML Config
Available configuration
User can start multiple grpc servers at the same time. Please make sure use different port and name.

### GRPC Service
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.name | The name of grpc server | string | N/A |
| grpc.port | The port of grpc server | integer | nil, server won't start |
| grpc.description | Description of grpc entry. | string | "" |
| grpc.reflection | Enable grpc server reflection | boolean | false |
| grpc.cert.ref | Reference of cert entry declared in [cert entry](https://github.com/rookie-ninja/rk-entry#certentry) | string | "" |
| grpc.logger.zapLogger.ref | Reference of zapLoggerEntry declared in [zapLoggerEntry](https://github.com/rookie-ninja/rk-entry#zaploggerentry) | string | "" |
| grpc.logger.eventLogger.ref | Reference of eventLoggerEntry declared in [eventLoggerEntry](https://github.com/rookie-ninja/rk-entry#eventloggerentry) | string | "" |

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

### GRPC Gateway Service
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.gw.enabled | Enable gateway service over gRpc server | boolean | false |
| grpc.gw.port | The port of gRpc gateway | integer | 0 |
| grpc.gw.rkServerOption | Enable RK style gateway server options. [detail](boot/gw_server_options.go) | false |
| grpc.gw.gwMappingFilePaths | The grpc gateway mapping file path. [example](boot/api/v1/gw_mapping.yaml) | string array | [] |
| grpc.gw.cert.ref | Reference of cert entry declared in cert section | string | "" |

#### Swagger Service
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.gw.sw.enabled | Enable swagger service over gRpc server | boolean | false |
| grpc.gw.sw.path | The path access swagger service from web | string | /sw |
| grpc.gw.sw.jsonPath | Where the swagger.json files are stored locally | string | "" |
| grpc.gw.sw.headers | Headers would be sent to caller as scheme of [key:value] | []string | [] |

#### Prom Client
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.gw.prom.enabled | Enable prometheus | boolean | false |
| grpc.gw.prom.path | Path of prometheus | string | /metrics |
| grpc.gw.prom.pusher.enabled | Enable prometheus pusher | bool | false |
| grpc.gw.prom.pusher.jobName | Job name would be attached as label while pushing to remote pushgateway | string | "" |
| grpc.gw.prom.pusher.remoteAddress | PushGateWay address, could be form of http://x.x.x.x or x.x.x.x | string | "" |
| grpc.gw.prom.pusher.intervalMs | Push interval in milliseconds | string | 1000 |
| grpc.gw.prom.pusher.basicAuth | Basic auth used to interact with remote pushgateway, form of [user:pass] | string | "" |
| grpc.gw.prom.pusher.cert.ref | Reference of rkentry.CertEntry | string | "" |

#### TV Service
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.gw.tv.enabled | Enable RK TV | boolean | false |

### Interceptors
#### Log
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.loggingZap.enabled | Enable log interceptor | boolean | false |

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
app={"appName":"rk-boot","appVersion":"master-xxx","entryName":"greeter","entryType":"GrpcEntry"}
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

#### Auth
Enable the server side auth. codes.Unauthenticated would be returned to client if not authorized with user defined credential.

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.auth.enabled | Enable auth interceptor | boolean | false |
| grpc.interceptors.auth.basic | Basic auth credentials as scheme of <user:pass> | []string | [] |
| grpc.interceptors.auth.bearer | Bearer auth tokens | []string | [] |
| grpc.interceptors.auth.api | API key | []string | [] |

#### Meta
Send application metadata as header to client and GRPC Gateway.

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.meta.enabled | Enable meta interceptor | boolean | false |
| grpc.interceptors.meta.prefix | Header key was formed as X-<Prefix>-XXX | string | RK |

#### Tracing
| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.interceptors.tracingTelemetry.enabled | Enable tracing interceptor | boolean | false |
| grpc.interceptors.exporter.file.enabled | Enable file exporter | boolean | RK |
| grpc.interceptors.exporter.file.outputPath | Export tracing info to files | string | stdout |
| grpc.interceptors.exporter.jaeger.enabled | Export tracing info jaeger | boolean | false |
| grpc.interceptors.exporter.jaeger.collectorEndpoint | As name described | string | localhost:16368/api/trace |
| grpc.interceptors.exporter.jaeger.collectorUsername | As name described | string | "" |
| grpc.interceptors.exporter.jaeger.collectorPassword | As name described | string | "" |

### Development Status: Stable

### Appendix
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

