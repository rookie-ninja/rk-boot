<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [rk-grpc](#rk-grpc)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
    - [Start gRpc server from YAML config](#start-grpc-server-from-yaml-config)
    - [Server side interceptor](#server-side-interceptor)
    - [Client side interceptor](#client-side-interceptor)
    - [Common Services](#common-services)
    - [Development Status: Stable](#development-status-stable)
    - [Appendix](#appendix)
    - [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# rk-grpc
Interceptor & bootstrapper designed for grpc. Currently, supports bellow interceptors

- logging & metrics
- auth
- panic
- bootstrapper

## Installation
`go get -u github.com/rookie-ninja/rk-grpc`

## Quick Start
Bootstrapper can be used with YAML config

### Start gRpc server from YAML config
User can access common service with localhost:8080/sw
```yaml
---
grpc:
  - name: greeter
    port: 1949
    enableCommonService: true
    gw:
      enabled: true
      port: 8080
      enableTV: true
      sw:
        enabled: true
        path: sw
    loggingInterceptor:
      enabled: true
      enableLogging: true
      enableMetrics: true
      enablePayloadLogging: true
```

```go
package main

import (
	"github.com/rookie-ninja/rk-grpc/boot"
	"github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
)

func main() {
	fac := rk_query.NewEventFactory()
	boot := rk_grpc.NewGRpcEntries("example/boot/boot.yaml", fac, rk_logger.StdoutLogger)
	boot["greeter"].Bootstrap(fac.CreateEvent())
}
```

Available configuration
User can start multiple servers at the same time

| name | description | type | default value |
| ------ | ------ | ------ | ------ |
| grpc.name | The name of gRpc server | string |
| grpc.port | The port of gRpc server | integer |
| grpc.enableCommonService | Enable embedded common service | true, false |
| grpc.tls.enabled | enable tls or not | boolean | false | 
| grpc.tls.user.enabled | enable user provided CA file? | boolean | false |
| grpc.tls.user.certFile | cert file path | string | empty string |
| grpc.tls.user.keyFile | key file path | string | empty string | 
| grpc.tls.auth.enabled | server will generate CA files | string | false |
| grpc.tls.auth.certOutput | cert file output path | string | current working directory | 
| grpc.gw.enabled | Enable gateway service over gRpc server | true, false |
| grpc.gw.port | The port of gRpc gateway | true, false |
| grpc.gw.insecure | Run gateway with insecure mode | true, false |
| grpc.gw.enableCommonService | Enable embedded common service | true, false |
| grpc.gw.enableTV | Enable RK TV | true, false |
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

### Server side interceptor

Example:
```go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/rookie-ninja/rk-grpc/example/proto"
	"github.com/rookie-ninja/rk-grpc/interceptor/context"
	"github.com/rookie-ninja/rk-grpc/interceptor/log/zap"
	"github.com/rookie-ninja/rk-grpc/interceptor/panic"
	"github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
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

	// create event factory
	factory := rk_query.NewEventFactory()

	// create server interceptor
	opt := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			rk_grpc_log.UnaryServerInterceptor(
				rk_grpc_log.WithEventFactory(factory),
				rk_grpc_log.WithLogger(rk_logger.StdoutLogger)),
			rk_grpc_panic.UnaryServerInterceptor(rk_grpc_panic.PanicToStderr)),
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
	event := rk_grpc_ctx.GetEvent(ctx)
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
	rk_grpc_ctx.AddToOutgoingMD(ctx, "key", "1", "2")
	// add request id
	rk_grpc_ctx.AddRequestIdToOutgoingMD(ctx)

	rk_grpc_ctx.GetLogger(ctx).Info("this is info message")

	return &proto.HelloResponse{
		Message: "hello",
	}, nil
}
```
Output
```
------------------------------------------------------------------------
end_time=2020-11-06T01:17:50.710002+08:00
start_time=2020-11-06T01:17:49.708046+08:00
time=1001
hostname=JEREMYYIN-MB0
event_id=["bb69e3d7-0a9f-4621-8987-7a468366be1c","37448bca-1b3f-4e51-8abb-1573dfcaaaa1"]
timing={"sleep.count":1,"sleep.elapsed_ms":1001}
counter={"ctr":1}
pair={"key":"value"}
error={"std-err":1}
field={"api.role":"unary_server","api.service":"Greeter","api.verb":"SayHello","app_version":"latest","az":"unknown","deadline":"2020-11-06T01:17:54+08:00","domain":"unknown","elapsed_ms":1001,"end_time":"2020-11-06T01:17:50.710002+08:00","incoming_request_id":["bb69e3d7-0a9f-4621-8987-7a468366be1c"],"key":"value","local.IP":"10.8.0.2","outgoing_request_id":["37448bca-1b3f-4e51-8abb-1573dfcaaaa1"],"realm":"unknown","region":"unknown","remote.IP":"localhost","remote.net_type":"tcp","remote.port":"61086","request_payload":"{\"name\":\"name\"}","res_code":"OK","response_payload":"{\"message\":\"hello\"}","start_time":"2020-11-06T01:17:49.708046+08:00"}
remote_addr=localhost
app_name=Unknown
operation=SayHello
event_status=Ended
history=s-sleep:1604596669708,e-sleep:1001,end:1
timezone=CST
os=darwin
arch=amd64
EOE
```

### Client side interceptor

Example:
```go
package main

import (
	"context"
	"encoding/json"
	"github.com/rookie-ninja/rk-grpc/example/proto"
	"github.com/rookie-ninja/rk-grpc/interceptor/context"
	"github.com/rookie-ninja/rk-grpc/interceptor/log/zap"
	"github.com/rookie-ninja/rk-grpc/interceptor/retry"
	"github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	// create event factory
	factory := rk_query.NewEventFactory()

	// create client interceptor
	opt := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			rk_grpc_log.UnaryClientInterceptor(
				rk_grpc_log.WithEventFactory(factory),
				rk_grpc_log.WithLogger(rk_logger.StdoutLogger)),
			rk_grpc_retry.UnaryClientInterceptor()),
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
	ctx, cancel := context.WithTimeout(rk_grpc_ctx.NewContext(), 5*time.Second)
	defer cancel()

	// add metadata
	rk_grpc_ctx.AddToOutgoingMD(ctx, "key", "1", "2")
	// add request id
	rk_grpc_ctx.AddRequestIdToOutgoingMD(ctx)

	// call server
	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: "name"})

	rk_grpc_ctx.GetLogger(ctx).Info("This is info message")

	// print incoming metadata
	bytes, _ := json.Marshal(rk_grpc_ctx.GetIncomingMD(ctx))
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
end_time=2020-11-06T01:17:50.710937+08:00
start_time=2020-11-06T01:17:49.706934+08:00
time=1004
hostname=JEREMYYIN-MB0
event_id=["37448bca-1b3f-4e51-8abb-1573dfcaaaa1","bb69e3d7-0a9f-4621-8987-7a468366be1c"]
timing={}
counter={"rk_max_retries":0}
pair={}
error={}
field={"api.role":"unary_client","api.service":"Greeter","api.verb":"SayHello","app_version":"latest","az":"unknown","deadline":"2020-11-06T01:17:54+08:00","domain":"unknown","elapsed_ms":1004,"end_time":"2020-11-06T01:17:50.710942+08:00","incoming_request_id":["37448bca-1b3f-4e51-8abb-1573dfcaaaa1"],"local.IP":"10.8.0.2","outgoing_request_id":["bb69e3d7-0a9f-4621-8987-7a468366be1c"],"realm":"unknown","region":"unknown","remote.IP":"localhost","remote.port":"8080","request_payload":"{\"name\":\"name\"}","res_code":"OK","response_payload":"{\"message\":\"hello\"}","start_time":"2020-11-06T01:17:49.706934+08:00"}
remote_addr=localhost
app_name=Unknown
operation=SayHello
event_status=Ended
timezone=CST
os=darwin
arch=amd64
EOE
```

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
Use bellow command to rebuild proto files
cd to ./boot folder
- protoc -I. --go_out=plugins=grpc:. --go_opt=paths=source_relative api/v1/*.proto
- protoc -I. --openapiv2_out=logtostderr=true,json_names_for_fields=false,grpc_api_configuration=api/v1/gw_mapping.yaml:. api/v1/*.proto
- protoc -I. --grpc-gateway_out=logtostderr=true,paths=source_relative,grpc_api_configuration=api/v1/gw_mapping.yaml:. api/v1/*.proto

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

