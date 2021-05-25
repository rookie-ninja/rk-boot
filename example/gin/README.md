# rk-gin
Interceptor & bootstrapper designed for gin framework.
Currently, supports bellow interceptors

- auth
- logging
- metrics
- panic

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Installation](#installation)
- [Quick Start](#quick-start)
  - [Start Gin server from YAML config](#start-gin-server-from-yaml-config)
  - [Start Gin server from code](#start-gin-server-from-code)
  - [Logging Interceptor](#logging-interceptor)
  - [Metrics interceptor](#metrics-interceptor)
  - [Panic interceptor](#panic-interceptor)
  - [Auth interceptor](#auth-interceptor)
  - [Common Service](#common-service)
  - [TV Service](#tv-service)
  - [Development Status: Stable](#development-status-stable)
  - [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation
`go get -u github.com/rookie-ninja/rk-gin`

## Quick Start
Bootstrapper can be used with YAML config.

### Start Gin server from YAML config
User can access common service with localhost:8080/sw
```yaml
---
rk: # NOT required
  appName: rk-example-entry           # Optional, default: "rkApp"
gin:
  - name: greeter                     # Required
    port: 8080                        # Required
    sw:                               # Optional
      enabled: true                   # Optional, default: false
    commonService:                    # Optional
      enabled: true                   # Optional, default: false
    interceptors:                     
      loggingZap:
        enabled: true                 # Optional, default: false
```

```go
func bootFromConfig() {
	// Bootstrap basic entries from boot config.
	rkentry.RegisterInternalEntriesFromConfig("example/boot/boot.yaml")

	// Bootstrap gin entry from boot config
	res := rkgin.RegisterGinEntriesWithConfig("example/boot/boot.yaml")

	// Bootstrap gin entry
	go res["greeter"].Bootstrap(context.Background())

	// Wait for shutdown signal
	rkentry.GlobalAppCtx.WaitForShutdownSig()

	// Interrupt gin entry
	res["greeter"].Interrupt(context.Background())
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

Interceptors can be used with chain.

### Start Gin server from code

```go
func bootFromCode() {
	// Create gin entry
	entry := rkgin.RegisterGinEntry(
		rkgin.WithNameGin("greeter"),
		rkgin.WithPortGin(8080),
		rkgin.WithCommonServiceEntryGin(rkgin.NewCommonServiceEntry()),
		rkgin.WithInterceptorsGin(rkginlog.LoggingZapInterceptor([]rkginlog.Option{}...)))

	// Start server
	go entry.Bootstrap(context.Background())

	// Wait for shutdown sig
	rkentry.GlobalAppCtx.WaitForShutdownSig()

	// Interrupt server
	entry.Interrupt(context.Background())
}
```

### Logging Interceptor
Logging interceptor uses [zap logger](https://github.com/uber-go/zap) and [rk-query](https://github.com/rookie-ninja/rk-query) logs every request.
[rk-prom](https://github.com/rookie-ninja/rk-prom) also used for prometheus metrics.

```go
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(
		rkginbasic.BasicInterceptor(),
		rkginlog.LoggingZapInterceptor(
			rkginlog.WithEventFactory(rkquery.NewEventFactory()),
			rkginlog.WithLogger(rklogger.StdoutLogger)),
		rkginpanic.PanicInterceptor(),
	)

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world")
	})
	router.Run(":8080")
```

Output: 
```log
------------------------------------------------------------------------
endTime=2021-05-15T00:38:18.028488+08:00
startTime=2021-05-15T00:38:18.028409+08:00
elapsedNano=79322
hostname=lark.local
timing={}
counter={}
pair={}
error={}
field={"apiMethod":"GET","apiPath":"/hello","apiProtocol":"HTTP/1.1","apiQuery":"","appName":"rkApp","appVersion":"v0.0.0","az":"unknown","domain":"unknown","elapsedNano":79322,"endTime":"2021-05-15T00:38:18.028488+08:00","entryName":"rkentry","entryType":"entry","incomingRequestIds":[],"localIp":"192.168.101.5","outgoingRequestIds":[],"realm":"unknown","region":"unknown","remoteIp":"localhost","remotePort":"50622","resCode":200,"startTime":"2021-05-15T00:38:18.028409+08:00","userAgent":"curl/7.64.1"}
remoteAddr=localhost:50622
appName=unknown
appVersion=unknown
entryName=rkentry
entryType=entry
locale=unknown
operation=GET:/hello
eventStatus=Ended
resCode=200
timezone=CST
os=darwin
arch=amd64
EOE
```

### Metrics interceptor
```go
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(
		rkginbasic.BasicInterceptor(),
        rkginmetrics.MetricsPromInterceptor(),
		rkginpanic.PanicInterceptor(),
	)

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world")
	})
	router.Run(":8080")
```


### Panic interceptor
```go
func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(
		rkginlog.LoggingZapInterceptor(
			rkginlog.WithEventFactory(rkquery.NewEventFactory()),
			rkginlog.WithLogger(rklogger.StdoutLogger)),
		rkginpanic.PanicInterceptor())

	router.GET("/hello", func(ctx *gin.Context) {
		panic(errors.New(""))
	})
	router.Run(":8080")
}
```
Output
```log
------------------------------------------------------------------------
endTime=2021-05-15T00:41:09.804115+08:00
startTime=2021-05-15T00:41:09.803402+08:00
elapsedNano=713047
hostname=lark.local
timing={}
counter={}
pair={}
error={}
field={"apiMethod":"GET","apiPath":"/hello","apiProtocol":"HTTP/1.1","apiQuery":"","appName":"rkApp","appVersion":"v0.0.0","az":"unknown","domain":"unknown","elapsedNano":721167,"endTime":"2021-05-15T00:41:09.804123+08:00","entryName":"rkentry","entryType":"entry","incomingRequestIds":[],"localIp":"10.8.0.2","outgoingRequestIds":[],"realm":"unknown","region":"unknown","remoteIp":"localhost","remotePort":"51347","resCode":500,"startTime":"2021-05-15T00:41:09.803402+08:00","userAgent":"curl/7.64.1"}
remoteAddr=localhost:51347
appName=unknown
appVersion=unknown
entryName=rkentry
entryType=entry
locale=unknown
operation=GET:/hello
eventStatus=Ended
resCode=500
timezone=CST
os=darwin
arch=amd64
EOE
```

### Auth interceptor
```go
func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(
		rkginlog.LoggingZapInterceptor(
			rkginlog.WithEventFactory(rkquery.NewEventFactory()),
			rkginlog.WithLogger(rklogger.StdoutLogger)),
		rkginauth.BasicAuthInterceptor(gin.Accounts{"user": "pass"}, "realm"),
		rkginpanic.PanicInterceptor())

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world")
	})
	router.Run(":8080")
}
```
Output
```log
------------------------------------------------------------------------
endTime=2021-05-15T00:48:12.436157+08:00
startTime=2021-05-15T00:48:12.435914+08:00
elapsedNano=243532
hostname=lark.local
timing={}
counter={}
pair={}
error={}
field={"apiMethod":"GET","apiPath":"/hello","apiProtocol":"HTTP/1.1","apiQuery":"","appName":"rkApp","appVersion":"v0.0.0","az":"unknown","domain":"unknown","elapsedNano":243532,"endTime":"2021-05-15T00:48:12.436157+08:00","entryName":"rkentry","entryType":"entry","incomingRequestIds":[],"localIp":"10.8.0.2","outgoingRequestIds":[],"realm":"unknown","region":"unknown","remoteIp":"localhost","remotePort":"53164","resCode":200,"startTime":"2021-05-15T00:48:12.435914+08:00","userAgent":"curl/7.64.1"}
remoteAddr=localhost:53164
appName=unknown
appVersion=unknown
entryName=rkentry
entryType=entry
locale=unknown
operation=GET:/hello
eventStatus=Ended
resCode=200
timezone=CST
os=darwin
arch=amd64
EOE
```

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

### Contributing
We encourage and support an active, healthy community of contributors &mdash;
including you! Details are in the [contribution guide](CONTRIBUTING.md) and
the [code of conduct](CODE_OF_CONDUCT.md). The rk maintainers keep an eye on
issues and pull requests, but you can also report any negative conduct to
dongxuny@gmail.com. That email list is a private, safe space; even the zap
maintainers don't have access, so don't hesitate to hold us to a high
standard.

Released under the [MIT License](LICENSE).

