// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-boot/example/grpc/api/v1"
	"google.golang.org/grpc"
	"time"
)

func main() {
	boot := rk_boot.NewBoot(rk_boot.WithBootConfigPath("example/grpc/configs/boot.yaml"))

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
