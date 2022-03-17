# Example

Init [go-redis](https://github.com/go-redis/redis) from YAML config.

This belongs to [rk-boot](https://github.com/rookie-ninja/rk-boot) family. We suggest use this lib from [rk-boot](https://github.com/rookie-ninja/rk-boot).

## Installation
- rk-boot: Bootstrapper base
- rk-gin: Bootstrapper for [gin-gonic/gin](https://github.com/gin-gonic/gin) Web Framework for API
- rk-db/redis: Bootstrapper for [gorm](https://github.com/go-gorm/gorm) of redis

```
go get github.com/rookie-ninja/rk-boot/v2
go get github.com/rookie-ninja/rk-gin/v2
go get github.com/rookie-ninja/rk-db/redis
```

## Quick Start
In the bellow example, we will run Redis locally and implement API of Get/Set of K/V.

- GET /v1/get, get value
- POST /v1/set, set value

### 1.Create boot.yaml
[boot.yaml](boot.yaml)

- Create web server with Gin framework at port 8080
- Create Redis entry which connects Redis at localhost:6379

```yaml
---
gin:
  - name: server
    enabled: true
    port: 8080
redis:
  - name: redis                      # Required
    enabled: true                    # Required
    addrs: ["localhost:6379"]        # Required, One addr is for single, multiple is for cluster
```

### 2.Create main.go

In the main() function, we implement bellow things.

- Register APIs into Gin router.

```go
// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rookie-ninja/rk-boot/v2"
	"github.com/rookie-ninja/rk-db/redis"
	"github.com/rookie-ninja/rk-gin/v2/boot"
	"net/http"
	"time"
)

var redisClient *redis.Client

func main() {
	boot := rkboot.NewBoot()

	boot.Bootstrap(context.TODO())

	// Auto migrate database and init global userDb variable
	redisEntry := rkredis.GetRedisEntry("redis")
	redisClient, _ = redisEntry.GetClient()

	// Register APIs
	ginEntry := rkgin.GetGinEntry("server")
	ginEntry.Router.GET("/v1/get", Get)
	ginEntry.Router.POST("/v1/set", Set)

	boot.WaitForShutdownSig(context.TODO())
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Set(ctx *gin.Context) {
	payload := &KV{}

	if err := ctx.BindJSON(payload); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	cmd := redisClient.Set(ctx.Request.Context(), payload.Key, payload.Value, time.Minute)

	if cmd.Err() != nil {
		ctx.JSON(http.StatusInternalServerError, cmd.Err())
		return
	}

	ctx.Status(http.StatusOK)
}

func Get(ctx *gin.Context) {
	key := ctx.Query("key")

	cmd := redisClient.Get(ctx.Request.Context(), key)

	if cmd.Err() != nil {
		ctx.JSON(http.StatusInternalServerError, cmd.Err())
		return
	}

	payload := &KV{
		Key:   key,
		Value: cmd.Val(),
	}

	ctx.JSON(http.StatusOK, payload)
}
```

### 3.Start server

```
$ go run main.go

2022-01-20T18:59:32.976+0800    INFO    boot/gin_entry.go:596   Bootstrap ginEntry      {"eventId": "8d1ec972-6439-4026-bedf-7e1f62724849", "entryName": "server"}
------------------------------------------------------------------------
endTime=2022-01-20T18:59:32.976769+08:00
startTime=2022-01-20T18:59:32.97669+08:00
elapsedNano=78808
timezone=CST
ids={"eventId":"8d1ec972-6439-4026-bedf-7e1f62724849"}
app={"appName":"rk","appVersion":"","entryName":"server","entryType":"Gin"}
env={"arch":"amd64","az":"*","domain":"*","hostname":"lark.local","localIP":"10.8.0.6","os":"darwin","realm":"*","region":"*"}
payloads={"ginPort":8080}
error={}
counters={}
pairs={}
timing={}
remoteAddr=localhost
operation=Bootstrap
resCode=OK
eventStatus=Ended
EOE
2022-01-20T18:59:32.976+0800    INFO    redis/boot.go:298       Bootstrap redis entry   {"entryName": "redis", "clientType": "Single"}
```

### 4.Validation
#### 4.1 Set value

```
$ curl -X POST "localhost:8080/v1/set" -d '{"key":"my-key","value":"my-value"}'
```

#### 4.2 Get value

```
$ curl -X GET "localhost:8080/v1/get?key=my-key"
{"key":"my-key","value":"my-value"}
```

## YAML Options
User can start multiple [go-redis](https://github.com/go-redis/redis) instances at the same time. Please make sure use different names.

Nearly all the fields were followed fields defined at [Option](https://github.com/go-redis/redis/blob/master/options.go)

```yaml
redis:
  - name: redis                      # Required
    enabled: true                    # Required
    addrs: ["localhost:6379"]        # Required, One addr is for single, multiple is for cluster
#    domain: "*"                     # Optional
#    description: ""                 # Optional
#
#    # For HA
#    mansterName: ""                 # Optional, required when connecting to Sentinel(HA)
#    sentinelPass: ""                # Optional, default: ""
#
#    # For cluster
#    maxRedirects: 3                 # Optional, default: 3
#    readOnly: false                 # Optional, default: false
#    routeByLatency: false           # Optional, default: false
#    routeRandomly: false            # Optional, default: false
#
#    # Common options
#    db: 0                           # Optional, default: 0
#    user: ""                        # Optional, default: ""
#    pass: ""                        # Optional, default: ""
#    maxRetries: 3                   # Optional, default: 3
#    minRetryBackoffMs: 8            # Optional, default: 8
#    maxRetryBackoffMs: 512          # Optional, default: 512
#    dialTimeoutMs: 5000             # Optional, default: 5000 (5 seconds)
#    readTimeoutMs: 3000             # Optional, default: 3000 (3 seconds)
#    writeTimeoutMs: 1               # Optional, default: 3000 (3 seconds)
#    poolFIFO: false                 # Optional, default: false
#    poolSize: 10                    # Optional, default: 10
#    minIdleConn: 0                  # Optional, default: 0
#    maxConnAgeMs: 0                 # Optional, default: no aged connection
#    poolTimeoutMs: 1300             # Optional, default: 1300 (1.3 seconds)
#    idleTimeoutMs: 1                # Optional, default: 5 minutes
#    idleCheckFrequencyMs: 1         # Optional, default: 1 minutes
#
#    # For logger
#    loggerEntry: ""                 # Optional, default: default logger with STDOUT
```

### Usage of domain

```
RK use <domain> to distinguish different environment.
Variable of <locale> could be composed as form of <domain>
- domain: Stands for different environment, like dev, test, prod and so on, users can define it by themselves.
          Environment variable: DOMAIN
          Eg: prod
          Wildcard: supported

How it works?
Firstly, get environment variable named as  DOMAIN.
Secondly, compare every element in locale variable and environment variable.
If variables in locale represented as wildcard(*), we will ignore comparison step.

Example:
# let's assuming we are going to define DB address which is different based on environment.
# Then, user can distinguish DB address based on locale.
# We recommend to include locale with wildcard.
---
DB:
  - name: redis-default
    domain: "*"
    addr: "192.0.0.1:6379"
  - name: redis-in-test
    domain: "test"
    addr: "192.0.0.1:6379"
  - name: redis-in-prod
    domain: "prod"
    addr: "176.0.0.1:6379"
```