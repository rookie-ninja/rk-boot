# Example
Full documentations:
- [rkdev.info](https://rkdev.info/docs/bootstrapper/user-guide/grpc-golang/)
- [rk-grpc](https://github.com/rookie-ninja/rk-grpc)

Interceptor & bootstrapper designed for [gRPC](https://grpc.io/docs/languages/go/) and [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).

![image](docs/img/grpc-arch.png)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Installation](#installation)
- [Quick start](#quick-start)
  - [1.Prepare .proto files](#1prepare-proto-files)
  - [2.Generate .pb.go files with buf](#2generate-pbgo-files-with-buf)
  - [4.Create boot.yaml](#4create-bootyaml)
  - [5.Create main.go](#5create-maingo)
  - [6.Start server](#6start-server)
  - [7.Validation](#7validation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation
`go get github.com/rookie-ninja/rk-boot/grpc`

## Quick start
### 1.Prepare .proto files
- api/v1/greeter.proto

```protobuf
syntax = "proto3";

package api.v1;

option go_package = "api/v1/hello";

service Greeter {
  rpc Hello (HelloRequest) returns (HelloResponse) {}
}

message HelloRequest {}

message HelloResponse {
  string my_message = 1;
}
```

- api/v1/gw_mapping.yaml

```yaml
type: google.api.Service
config_version: 3

# Please refer google.api.Http in third-party/googleapis/google/api/http.proto file for details.
http:
  rules:
    - selector: api.v1.Greeter.Hello
      get: /v1/hello
```

- buf.yaml

```yaml
version: v1beta1
name: github.com/rk-dev/rk-boot
build:
  roots:
    - api
    - third-party/googleapis
```

- buf.gen.yaml

```yaml
version: v1beta1
plugins:
  - name: go
    out: api/gen
    opt:
     - paths=source_relative
  - name: go-grpc
    out: api/gen
    opt:
      - paths=source_relative
      - require_unimplemented_servers=false
  - name: grpc-gateway
    out: api/gen
    opt:
      - paths=source_relative
      - grpc_api_configuration=api/v1/gw_mapping.yaml
      - allow_repeated_fields_in_body=true
      - generate_unbound_methods=true
  - name: openapiv2
    out: api/gen
    opt:
      - grpc_api_configuration=api/v1/gw_mapping.yaml
      - allow_repeated_fields_in_body=true
```

### 2.Generate .pb.go files with [buf](https://docs.buf.build/introduction)
```
$ buf generate --path api/v1
```

### 4.Create boot.yaml
```yaml
grpc:
  - name: greeter                                      # Required
    enabled: true                                      # Required
    port: 8080                                         # Optional
    enableRkGwOption: true                             # Optional, default: false
```

### 5.Create main.go
```go
// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-boot/grpc"
	"github.com/rookie-ninja/rk-demo/api/gen/v1"
	"google.golang.org/grpc"
)

func main() {
	boot := rkboot.NewBoot()

	// register grpc
	entry := rkbootgrpc.GetGrpcEntry("greeter")
	entry.AddRegFuncGrpc(registerGreeter)
	entry.AddRegFuncGw(hello.RegisterGreeterHandlerFromEndpoint)

	// Bootstrap
	boot.Bootstrap(context.TODO())

	// Wait for shutdown sig
	boot.WaitForShutdownSig(context.TODO())
}

func registerGreeter(server *grpc.Server) {
	hello.RegisterGreeterServer(server, &GreeterServer{})
}

//GreeterServer GreeterServer struct
type GreeterServer struct{}

// Hello response with hello message
func (server *GreeterServer) Hello(_ context.Context, _ *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{
		MyMessage: "hello!",
	}, nil
}
```

### 6.Start server

```go
$ go run main.go
```

### 7.Validation
```shell script
$ curl -X GET localhost:8080/v1/greeter
{"myMessage":"hello!"}
```
