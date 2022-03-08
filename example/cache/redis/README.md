# Example
Init cache with [go-redis/cache](https://github.com/go-redis/cache/v8) or local memory cache from YAML config.

This belongs to [rk-boot](https://github.com/rookie-ninja/rk-boot) family. We suggest use this lib from [rk-boot](https://github.com/rookie-ninja/rk-boot).

## Installation
- rk-boot: Bootstrapper base
- rk-gin: Bootstrapper for [gin-gonic/gin](https://github.com/gin-gonic/gin) Web Framework for API
- rk-cache/redis: Bootstrapper for [go-redis/cache](https://github.com/go-redis/cache/v8) of cache

```
go get github.com/rookie-ninja/rk-boot/v2
go get github.com/rookie-ninja/rk-gin/v2
go get github.com/rookie-ninja/rk-cache/redis
```

## Quick Start
In the bellow example, we will run Redis locally and implement API of Get/Set of K/V.

- GET /v1/get, get value
- POST /v1/set, set value

### 1.Create boot.yaml
[boot.yaml](example/boot.yaml)

- Create web server with Gin framework at port 8080
- Create Redis entry which connects Redis at localhost:6379 as cache

```yaml
---
gin:
  - name: cache-service
    port: 8080
    enabled: true
cache:
  - name: redis-cache
    enabled: true
    local:
      enabled: false
    redis:
      enabled: true
```

### 2.Create main.go

In the main() function, we implement bellow things.

- Register APIs into Gin router.

```go
package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rookie-ninja/rk-boot/v2"
	"github.com/rookie-ninja/rk-cache/redis"
	"github.com/rookie-ninja/rk-gin/v2/boot"
	"net/http"
)

var cacheEntry *rkcache.CacheEntry

func main() {
	boot := rkboot.NewBoot()

	boot.Bootstrap(context.TODO())

	// assign cache
	cacheEntry = rkcache.GetCacheEntry("redis-cache")

	// assign router
	ginEntry := rkgin.GetGinEntry("cache-service")
	ginEntry.Router.GET("/v1/get", Get)
	ginEntry.Router.GET("/v1/set", Set)

	boot.WaitForShutdownSig(context.TODO())
}

func Get(ctx *gin.Context) {
	val := ""
	resp := cacheEntry.GetFromCache(&rkcache.CacheReq{
		Key:   "demo-key",
		Value: &val,
	})

	if resp.Error != nil || !resp.Success {
		ctx.JSON(http.StatusInternalServerError, resp.Error)
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{
		"value": val,
	})
}

func Set(ctx *gin.Context) {
	val, ok := ctx.GetQuery("value")
	if !ok {
		ctx.JSON(http.StatusBadRequest, "No value found")
	}

	cacheEntry.AddToCache(&rkcache.CacheReq{
		Key:   "demo-key",
		Value: val,
	})

	ctx.JSON(http.StatusOK, map[string]string{
		"value": val,
	})
}
```

### 3.Start server

```shell
$ go run main.go

2022-03-09T00:10:50.325+0800    INFO    redis/entry.go:163      Bootstrap CacheRedisEntry       {"entryName": "redis-cache", "localCache": false, "redisCache": true}
2022-03-09T00:10:50.325+0800    INFO    redis@v1.0.1/boot.go:253        Bootstrap redisEntry    {"eventId": "2e47b54a-a46c-4c23-ae51-f105bdbc5836", "entryName": "redis-cache", "entryType": "RedisEntry", "clientType": "Single"}
2022-03-09T00:10:50.325+0800    INFO    redis@v1.0.1/boot.go:261        Ping redis at [localhost:6379]
2022-03-09T00:10:50.330+0800    INFO    redis@v1.0.1/boot.go:267        Ping redis at [localhost:6379] success
2022-03-09T00:10:50.330+0800    INFO    boot/gin_entry.go:624   Bootstrap GinEntry      {"eventId": "2e47b54a-a46c-4c23-ae51-f105bdbc5836", "entryName": "cache-service", "entryType": "GinEntry"}
------------------------------------------------------------------------
endTime=2022-03-09T00:10:50.33088+08:00
startTime=2022-03-09T00:10:50.330832+08:00
elapsedNano=47527
timezone=CST
ids={"eventId":"2e47b54a-a46c-4c23-ae51-f105bdbc5836"}
app={"appName":"","appVersion":"","entryName":"cache-service","entryType":"GinEntry"}
env={"arch":"amd64","az":"*","domain":"*","hostname":"lark.local","localIP":"10.8.0.2","os":"darwin","realm":"*","region":"*"}
payloads={"ginPort":8080}
counters={}
pairs={}
timing={}
remoteAddr=localhost
operation=Bootstrap
resCode=OK
eventStatus=Ended
EOE
```

### 4.Validation
#### 4.1 Set value

```shell
$ curl "localhost:8080/v1/set?value=my-value"
{"value":"my-value"}
```

#### 4.2 Get value

```shell
$ curl localhost:8080/v1/get
{"value":"my-value"}
```

#### 4.3 Validate Redis
The key of cache will be encoded with MD5 and value will be base64 encoded.

![image](docs/img/redis.png)


