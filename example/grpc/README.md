# rk-grpc
Interceptor & bootstrapper designed for grpc. Currently, supports bellow interceptors

- logging
- metrics
- auth
- panic

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Installation](#installation)
- [Quick Start](#quick-start)
  - [Start gRpc server from YAML config](#start-grpc-server-from-yaml-config)
  - [Logging interceptor](#logging-interceptor)
  - [Client side interceptor](#client-side-interceptor)
  - [Common Services](#common-services)
  - [TV Service](#tv-service)
  - [Development Status: Stable](#development-status-stable)
  - [Appendix](#appendix)
  - [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation
`go get -u github.com/rookie-ninja/rk-grpc`

## Quick Start
Bootstrapper can be used with YAML config

### Start gRpc server from YAML config
User can access common service with localhost:8080/sw
```yaml
---
rk: # NOT required
  appName: rk-example-entry           # Optional, default: "rkApp"
grpc:
  - name: greeter
    port: 1949
    commonService:
      enabled: true
    gw:
      enabled: true
      port: 8080
      tv:
        enabled: true
      sw:
        enabled: true
      prom:
        enabled: true
    interceptors:
      loggingZap:
        enabled: true
      metricsProm:
        enabled: true
```

```go
func bootFromConfig() {
	// Bootstrap basic entries from boot config.
	rkentry.RegisterInternalEntriesFromConfig("example/boot/boot.yaml")

	// Bootstrap grpc entry from boot config
	res := rkgrpc.RegisterGrpcEntriesWithConfig("example/boot/boot.yaml")

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

### Logging interceptor

Example:
```go
// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-grpc/example/interceptor/proto"
	"github.com/rookie-ninja/rk-grpc/interceptor/auth/basic_auth"
	"github.com/rookie-ninja/rk-grpc/interceptor/basic"
	"github.com/rookie-ninja/rk-grpc/interceptor/context"
	"github.com/rookie-ninja/rk-grpc/interceptor/log/zap"
	"github.com/rookie-ninja/rk-grpc/interceptor/metrics/prom"
	"github.com/rookie-ninja/rk-grpc/interceptor/panic"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	// create listener
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	entryName := "example-entry-server-name"
	entryType := "example-entry-client"

	// create server interceptor
	// basic interceptor
	basicInter := rkgrpcbasic.UnaryServerInterceptor(
		rkgrpcbasic.WithEntryNameAndType(entryName, entryType))

	// logging interceptor
	logInter := rkgrpclog.UnaryServerInterceptor(
		rkgrpclog.WithEntryNameAndType(entryName, entryType),
		rkgrpclog.WithZapLoggerEntry(rkentry.GlobalAppCtx.GetZapLoggerEntryDefault()),
		rkgrpclog.WithEventLoggerEntry(rkentry.GlobalAppCtx.GetEventLoggerEntryDefault()))

	// prometheus metrics interceptor
	metricsInter := rkgrpcmetrics.UnaryServerInterceptor(
		rkgrpcmetrics.WithEntryNameAndType(entryName, entryType),
		rkgrpcmetrics.WithRegisterer(prometheus.NewRegistry()))

	// basic auth interceptor
	basicAuthInter := rkgrpcbasicauth.UnaryServerInterceptor(
		rkgrpcbasicauth.WithEntryNameAndType(entryName, entryType),
		rkgrpcbasicauth.WithCredential("user:name"))

	// panic interceptor
	panicInter := rkgrpcpanic.UnaryServerInterceptor()

	opt := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			basicInter,
			logInter,
			metricsInter,
			basicAuthInter,
			panicInter),
	}

	// create server
	s := grpc.NewServer(opt...)
	proto.RegisterGreeterServer(s, &GreeterServer{})

	// serving
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type GreeterServer struct{}

func (server *GreeterServer) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloResponse, error) {
	event := rkgrpcctx.GetEvent(ctx)
	// add fields
	event.AddFields(zap.String("key", "value"))
	// add error
	event.AddErr(errors.New(""))
	// add pair
	event.AddPair("key", "value")
	// set counter
	event.SetCounter("ctr", 1)
	// timer
	event.StartTimer("sleep")
	time.Sleep(1 * time.Second)
	event.EndTimer("sleep")
	// add to metadata
	rkgrpcctx.AddToOutgoingMD(ctx, "key", "1", "2")
	// add request id
	rkgrpcctx.AddRequestIdToOutgoingMD(ctx)

	// print incoming metadata
	bytes, _ := json.Marshal(rkgrpcctx.GetIncomingMD(ctx))
	println(string(bytes))

	// print with logger to check whether id was printed
	rkgrpcctx.GetZapLogger(ctx).Info("this is info message")

	return &proto.HelloResponse{
		Message: "hello",
	}, nil
}
```
Output
```
------------------------------------------------------------------------
endTime=2021-05-25T02:08:03.348494+08:00
startTime=2021-05-25T02:08:03.348264+08:00
elapsedNano=229639
hostname=lark.local
eventId=fakeId
timing={}
counter={}
pair={}
error={}
field={"appName":"rkApp","appVersion":"v0.0.0","az":"unknown","deadline":"2021-05-25T02:08:08+08:00","domain":"unknown","elapsedNano":229639,"endTime":"2021-05-25T02:08:03.348494+08:00","grpcMethod":"SayHello","grpcService":"Greeter","grpcType":"unaryServer","gwMethod":"unknown","gwPath":"unknown","incomingRequestId":["76bec784-de98-4e52-98a0-ee106bac6600"],"localIp":"10.8.0.6","outgoingRequestId":[],"realm":"unknown","region":"unknown","remoteIp":"localhost","remoteNetType":"tcp","remotePort":"58222","resCode":"Unauthenticated","startTime":"2021-05-25T02:08:03.348264+08:00"}
remoteAddr=localhost
appName=unknown
appVersion=unknown
locale=unknown
operation=SayHello
eventStatus=Ended
timezone=CST
os=darwin
arch=amd64
EOE
```

### Client side interceptor

Example:
```go
// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-grpc/example/interceptor/proto"
	"github.com/rookie-ninja/rk-grpc/interceptor/basic"
	"github.com/rookie-ninja/rk-grpc/interceptor/context"
	"github.com/rookie-ninja/rk-grpc/interceptor/log/zap"
	"github.com/rookie-ninja/rk-grpc/interceptor/metrics/prom"
	"github.com/rookie-ninja/rk-grpc/interceptor/panic"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	entryName := "example-entry-client-name"
	entryType := "example-entry-client"

	// create server interceptor
	basicInter := rkgrpcbasic.UnaryClientInterceptor(
		rkgrpcbasic.WithEntryNameAndType(entryName, entryType))

	logInter := rkgrpclog.UnaryClientInterceptor(
		rkgrpclog.WithEntryNameAndType(entryName, entryType),
		rkgrpclog.WithZapLoggerEntry(rkentry.GlobalAppCtx.GetZapLoggerEntryDefault()),
		rkgrpclog.WithEventLoggerEntry(rkentry.GlobalAppCtx.GetEventLoggerEntryDefault()))

	metricsInter := rkgrpcmetrics.UnaryClientInterceptor(
		rkgrpcmetrics.WithEntryNameAndType(entryName, entryType),
		rkgrpcmetrics.WithRegisterer(prometheus.NewRegistry()))

	panicInter := rkgrpcpanic.UnaryClientInterceptor()

	// create client interceptor
	opt := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			basicInter,
			logInter,
			metricsInter,
			panicInter,
		),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}

	// Set up a connection to the server.
	conn, err := grpc.DialContext(context.Background(), "localhost:8080", opt...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// create grpc client
	c := proto.NewGreeterClient(conn)
	// create with rk context
	ctx, cancel := context.WithTimeout(rkgrpcctx.NewContext(), 5*time.Second)
	defer cancel()

	// add metadata
	rkgrpcctx.AddToOutgoingMD(ctx, "key", "1", "2")
	// add request id
	rkgrpcctx.AddRequestIdToOutgoingMD(ctx)

	// call server
	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: "name"})

	rkgrpcctx.GetZapLogger(ctx).Info("This is info message")

	// print incoming metadata
	bytes, _ := json.Marshal(rkgrpcctx.GetIncomingMD(ctx))
	println(string(bytes))

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
```
Output 
```
------------------------------------------------------------------------
endTime=2021-05-25T02:08:03.349978+08:00
startTime=2021-05-25T02:08:03.347594+08:00
elapsedNano=2384154
hostname=lark.local
eventId=fakeId
timing={}
counter={}
pair={}
error={}
field={"appName":"rkApp","appVersion":"v0.0.0","az":"unknown","deadline":"2021-05-25T02:08:08+08:00","domain":"unknown","elapsedNano":2384154,"endTime":"2021-05-25T02:08:03.34998+08:00","grpcMethod":"SayHello","grpcService":"","grpcType":"unaryClient","incomingRequestId":[],"localIp":"10.8.0.6","outgoingRequestId":["76bec784-de98-4e52-98a0-ee106bac6600"],"realm":"unknown","region":"unknown","remoteIp":"localhost","remotePort":"8080","resCode":"Unauthenticated","startTime":"2021-05-25T02:08:03.347594+08:00"}
remoteAddr=localhost
appName=unknown
appVersion=unknown
locale=unknown
operation=SayHello
eventStatus=Ended
timezone=CST
os=darwin
arch=amd64
EOE

```

### Common Services
User can start multiple servers at the same time

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

