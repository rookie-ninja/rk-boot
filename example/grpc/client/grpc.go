// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"encoding/json"
	"github.com/rookie-ninja/rk-boot/example/grpc/api/v1"
	"github.com/rookie-ninja/rk-grpc/interceptor/context"
	"github.com/rookie-ninja/rk-grpc/interceptor/log/zap"
	"github.com/rookie-ninja/rk-grpc/interceptor/retry"
	"github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
	"path"
	"time"
)

func main() {
	// create event factory
	factory := rk_query.NewEventFactory(rk_query.WithAppName("app"))

	wd, _ := os.Getwd()

	creds, err := credentials.NewClientTLSFromFile(path.Join(wd, "example/grpc/server/cert/ca.pem"), "")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// create client interceptor
	opt := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			rk_grpc_log.UnaryClientInterceptor(
				rk_grpc_log.WithEventFactory(factory),
				rk_grpc_log.WithLogger(rk_logger.StdoutLogger)),
			rk_grpc_retry.UnaryClientInterceptor()),
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
	}

	// Set up a connection to the server.
	conn, err := grpc.DialContext(context.Background(), "localhost:8080", opt...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// create grpc client
	c := hello_v1.NewGreeterClient(conn)
	// create with rk context
	ctx, cancel := context.WithTimeout(rk_grpc_ctx.NewContext(), 5*time.Second)
	defer cancel()

	// add metadata
	rk_grpc_ctx.AddToOutgoingMD(ctx, "key", "1", "2")
	// add request id
	rk_grpc_ctx.AddRequestIdToOutgoingMD(ctx)

	// call server
	r, err := c.SayHello(ctx, &hello_v1.HelloRequest{Name: "name"})

	// print incoming metadata
	bytes, _ := json.Marshal(rk_grpc_ctx.GetIncomingMD(ctx))
	println(string(bytes))

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
