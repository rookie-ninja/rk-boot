# Example
Full documentations:
- [rkdev.info](https://rkdev.info/docs/bootstrapper/user-guide/echo-golang/)
- [rk-echo](https://github.com/rookie-ninja/rk-echo)

Interceptor & bootstrapper designed for [labstack/echo](https://github.com/labstack/echo) web framework. 

![image](docs/img/echo-arch.png)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Installation](#installation)
- [Quick start](#quick-start)
  - [1.Create boot.yaml](#1create-bootyaml)
  - [2.Create main.go](#2create-maingo)
  - [3.Start server](#3start-server)
  - [4.Validation](#4validation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation
`go get github.com/rookie-ninja/rk-boot/echo`

## Quick start
### 1.Create boot.yaml
```yaml
---
echo:
  - name: greeter       # Required, Name of echo entry
    port: 8080          # Required, Port of echo entry
    enabled: true       # Required, Enable echo entry
```

### 2.Create main.go
```go
// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-boot/echo"
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
	boot := rkboot.NewBoot(rkboot.WithBootConfigPath("example/echo/boot.yaml"))

	// Register handler
	entry := rkbootecho.GetEchoEntry("greeter")
	entry.Echo.GET("/v1/hello", hello)

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
func hello(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "hello")
}

```

### 3.Start server

```go
$ go run main.go
```

### 4.Validation
```shell script
$ curl -X GET localhost:8080/v1/greeter
"hello"
```
