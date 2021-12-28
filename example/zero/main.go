// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-boot/zero"
	"github.com/tal-tech/go-zero/rest"
	"net/http"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample rk-demo server.
// @termsOfService http://swagger.io/terms/

// @securityDefinitions.basic BasicAuth

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	// Create a new boot instance.
	boot := rkboot.NewBoot()

	// Register handler
	entry := rkbootzero.GetZeroEntry("greeter")
	entry.Server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/v1/hello",
		Handler: hello,
	})

	// Bootstrap
	boot.Bootstrap(context.TODO())

	boot.WaitForShutdownSig(context.TODO())
}

// @Summary Hello
// @Id 1
// @Tags Hello
// @version 1.0
// @produce application/json
// @Success 200 string string
// @Router /v1/hello [get]
func hello(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("hello"))
}
