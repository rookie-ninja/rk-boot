// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-example/api/gen/v1"
	"google.golang.org/grpc"
)

func main() {
	boot := rkboot.NewBoot(rkboot.WithBootConfigPath("example/grpc/boot.yaml"))

	// register grpc
	boot.GetGrpcEntry("greeter").AddGrpcRegFuncs(registerGreeter)
	boot.GetGrpcEntry("greeter").AddGwRegFuncs(hello.RegisterGreeterHandlerFromEndpoint)

	// Bootstrap
	boot.Bootstrap(context.TODO())

	// Wait for shutdown sig
	boot.WaitForShutdownSig(context.TODO())
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
