// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"github.com/rookie-ninja/rk-boot/example/grpc/api/gen/v1"
	"github.com/rookie-ninja/rk-grpc/interceptor/basic"
	"github.com/rookie-ninja/rk-grpc/interceptor/context"
	"github.com/rookie-ninja/rk-grpc/interceptor/log/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
	"path"
	"time"
)

func main() {
	wd, _ := os.Getwd()

	cred, err := credentials.NewClientTLSFromFile(path.Join(wd, "example/grpc/certs/server.pem"), "")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// create client interceptor
	opt := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			rkgrpcbasic.UnaryClientInterceptor(
				rkgrpcbasic.WithEntryNameAndType("fake-entry-name", "fake-entry")),
			rkgrpclog.UnaryClientInterceptor()),
		grpc.WithTransportCredentials(cred),
		grpc.WithBlock(),
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	// Set up a connection to the server.
	conn, err := grpc.DialContext(ctx, ":1949", opt...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// create grpc client
	c := hello.NewGreeterClient(conn)
	// create with rk context
	ctx, cancel := context.WithTimeout(rkgrpcctx.NewContext(), 5*time.Second)
	defer cancel()

	// add metadata
	rkgrpcctx.AddToOutgoingMD(ctx, "key", "1", "2")
	// add request id
	rkgrpcctx.AddRequestIdToOutgoingMD(ctx)

	// call server
	r, err := c.SayHello(ctx, &hello.HelloRequest{Name: "rk-dev"})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
