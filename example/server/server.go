// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-boot/example/proto"
	"google.golang.org/grpc"
	"time"
)

func main() {
	boot := rk_boot.NewBoot(rk_boot.WithBootConfigPath("example/configs/boot.yaml"))

	// register gRpc
	boot.GetGRpcEntry("greeter").AddRegFuncs(registerGreeter)
	boot.GetGRpcEntry("greeter").AddGWRegFuncs(proto.RegisterGreeterHandlerFromEndpoint)

	boot.Bootstrap()
	boot.Quitter(5 * time.Second)
}

func registerGreeter(server *grpc.Server) {
	proto.RegisterGreeterServer(server, &GreeterServer{})
}

type GreeterServer struct{}

func (server *GreeterServer) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{
		Message: "hello",
	}, nil
}
