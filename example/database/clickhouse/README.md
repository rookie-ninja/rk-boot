# Example

Init [gorm](https://github.com/go-gorm/gorm) from YAML config.

This belongs to [rk-boot](https://github.com/rookie-ninja/rk-boot) family. We suggest use this lib from [rk-boot](https://github.com/rookie-ninja/rk-boot).

## Installation
- rk-boot: Bootstrapper base
- rk-gin: Bootstrapper for [gin-gonic/gin](https://github.com/gin-gonic/gin) Web Framework for API
- rk-db/clickhouse: Bootstrapper for [gorm](https://github.com/go-gorm/gorm) of clickhouse

```
go get github.com/rookie-ninja/rk-boot/v2
go get github.com/rookie-ninja/rk-gin/v2
go get github.com/rookie-ninja/rk-db/clickhouse
```

## Quick Start
In the bellow example, we will run ClickHouse locally and implement API of Create/List/Get/Update/Delete for User model with Gin.

- GET /v1/user, List users
- GET /v1/user/:id, Get user
- PUT /v1/user, Create user
- POST /v1/user/:id, Update user
- DELETE /v1/user/:id, Delete user

### 1.Create boot.yaml
[boot.yaml](boot.yaml)

- Create web server with Gin framework at port 8080
- Create ClickHouse entry which connects ClickHouse at localhost:9000

```yaml
---
gin:
  - name: user-service
    port: 8080
    enabled: true
clickhouse:
  - name: user-db                          # Required
    enabled: true                          # Required
    locale: "*::*::*::*"                   # Required
    addr: "localhost:9000"                 # Optional, default: localhost:9000
    user: default                          # Optional, default: default
    pass: ""                               # Optional, default: ""
    database:
      - name: user                         # Required
        autoCreate: true                   # Optional, default: false
#        dryRun: false                     # Optional, default: false
#        params: []                        # Optional, default: []
#    loggerEntry: ""                       # Optional, default: default logger with STDOUT
```

### 2.Create main.go

In the main() function, we implement two things.

- Add User{} as auto migrate option which will create table in DB automatically if missing.
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
	"github.com/rookie-ninja/rk-boot/v2"
	"github.com/rookie-ninja/rk-db/clickhouse"
	"github.com/rookie-ninja/rk-gin/v2/boot"
	"github.com/rs/xid"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var userDb *gorm.DB

func main() {
	boot := rkboot.NewBoot()

	boot.Bootstrap(context.TODO())

	// Auto migrate database and init global userDb variable
	clickHouseEntry := rkclickhouse.GetClickHouseEntry("user-db")
	userDb = clickHouseEntry.GetDB("user")
	if !userDb.DryRun {
		userDb.AutoMigrate(&User{})
	}

	// Register APIs
	ginEntry := rkgin.GetGinEntry("user-service")
	ginEntry.Router.GET("/v1/user", ListUsers)
	ginEntry.Router.GET("/v1/user/:id", GetUser)
	ginEntry.Router.PUT("/v1/user", CreateUser)
	ginEntry.Router.POST("/v1/user/:id", UpdateUser)
	ginEntry.Router.DELETE("/v1/user/:id", DeleteUser)

	boot.WaitForShutdownSig(context.TODO())
}

// *************************************
// *************** Model ***************
// *************************************

type Base struct {
	CreatedAt time.Time `yaml:"-" json:"-"`
	UpdatedAt time.Time `yaml:"-" json:"-"`
}

type User struct {
	Base
	Id   string `yaml:"id" json:"id"`
	Name string `yaml:"name" json:"name"`
}

func ListUsers(ctx *gin.Context) {
	userList := make([]*User, 0)
	res := userDb.Find(&userList)

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, res.Error)
		return
	}
	ctx.JSON(http.StatusOK, userList)
}

func GetUser(ctx *gin.Context) {
	uid := ctx.Param("id")
	user := &User{}
	res := userDb.Find(user, "id = ?", uid)

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, res.Error)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func CreateUser(ctx *gin.Context) {
	user := &User{
		Id:   xid.New().String(),
		Name: ctx.Query("name"),
	}

	res := userDb.Create(user)

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, res.Error)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func UpdateUser(ctx *gin.Context) {
	uid := ctx.Param("id")
	user := &User{
		Id:   uid,
		Name: ctx.Query("name"),
	}

	res := userDb.Where("id = ?", uid).Updates(user)

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, res.Error)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func DeleteUser(ctx *gin.Context) {
	uid := ctx.Param("id")

	res := userDb.Delete(&User{}, "id = ?", uid)

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, res.Error)
		return
	}

	ctx.String(http.StatusOK, "success")
}
```

### 3.Start server

```
$ go run main.go

2022-01-07T03:11:18.538+0800    INFO    boot/gin_entry.go:913   Bootstrap ginEntry      {"eventId": "181b17a7-591f-419a-95cc-2cda7efc61f2", "entryName": "user-service"}
------------------------------------------------------------------------
endTime=2022-01-07T03:11:18.53883+08:00
startTime=2022-01-07T03:11:18.538741+08:00
elapsedNano=88391
timezone=CST
ids={"eventId":"181b17a7-591f-419a-95cc-2cda7efc61f2"}
app={"appName":"rk","appVersion":"","entryName":"user-service","entryType":"GinEntry"}
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
2022-01-07T03:11:18.538+0800    INFO    Bootstrap ClickHouse entry      {"entryName": "user-db", "clickHouseUser": "default", "clickHouseAddr": "localhost:9000"}
2022-01-07T03:11:18.538+0800    INFO    creating database user if not exists
2022-01-07T03:11:18.556+0800    INFO    creating successs or database user exists
2022-01-07T03:11:18.556+0800    INFO    connecting to database user
2022-01-07T03:11:18.567+0800    INFO    connecting to database user success
```

### 4.Validation
#### 4.1 Create user
Create a user with name of rk-dev.

```shell script
$ curl -X PUT "localhost:8080/v1/user?name=rk-dev"
{"id":"c7bjufjd0cvqfaenpqjg","name":"rk-dev"}
```

#### 4.1 Update user
Update user name to rk-dev-updated.

```shell script
$ curl -X POST "localhost:8080/v1/user/c7bjufjd0cvqfaenpqjg?name=rk-dev-updated"
{"id":"c7bjufjd0cvqfaenpqjg","name":"rk-dev-updated"}
```

#### 4.1 List users
List users.

```shell script
$ curl -X GET localhost:8080/v1/user
[{"id":"c7bjufjd0cvqfaenpqjg","name":"rk-dev-updated"}]
```

#### 4.1 Get user
Get user with id=c7bjtobd0cvqfaenpqj0.

```shell script
$ curl -X GET localhost:8080/v1/user/c7bjufjd0cvqfaenpqjg
{"id":"c7bjufjd0cvqfaenpqjg","name":"rk-dev-updated"}
```

#### 4.1 Delete user

```shell script
$ curl -X DELETE localhost:8080/v1/user/c7bjufjd0cvqfaenpqjg
success
```

## YAML Options
User can start multiple [gorm](https://github.com/go-gorm/gorm) instances at the same time. Please make sure use different names.

| name                           | Required | description                        | type     | default value  |
|--------------------------------|----------|------------------------------------|----------|----------------|
| clickhouse.name                | Required | The name of entry                  | string   | ClickHouse     |
| clickhouse.enabled             | Required | Enable entry or not                | bool     | false          |
| clickhouse.domain              | Required | See domain description bellow      | string   | ""             |
| clickhouse.description         | Optional | Description of echo entry.         | string   | ""             |
| clickhouse.user                | Optional | ClickHouse username                | string   | root           |
| clickhouse.pass                | Optional | ClickHouse password                | string   | pass           |
| clickhouse.addr                | Optional | ClickHouse remote address          | string   | localhost:9000 |
| clickhouse.database.name       | Required | Name of database                   | string   | ""             |
| clickhouse.database.autoCreate | Optional | Create DB if missing               | bool     | false          |
| clickhouse.database.dryRun     | Optional | Run gorm.DB with dry run mode      | bool     | false          |
| clickhouse.database.params     | Optional | Connection params                  | []string | [""]           |
| clickhouse.loggerEntry         | Optional | Reference of zap logger entry name | string   | ""             |

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