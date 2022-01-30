# Example

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Quick Start](#quick-start)
  - [Installation](#installation)
  - [Start SQL Server with docker](#start-sql-server-with-docker)
  - [1.Create boot.yaml](#1create-bootyaml)
  - [2.Create main.go](#2create-maingo)
  - [3.Start server](#3start-server)
  - [4.Validation](#4validation)
    - [4.1 Create user](#41-create-user)
    - [4.1 Update user](#41-update-user)
    - [4.1 List users](#41-list-users)
    - [4.1 Get user](#41-get-user)
    - [4.1 Delete user](#41-delete-user)
- [YAML Options](#yaml-options)
  - [Usage of locale](#usage-of-locale)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Quick Start
In the bellow example, we will run SQL Server locally and implement API of Create/List/Get/Update/Delete for User model with Gin.

- GET /v1/user, List users
- GET /v1/user/:id, Get user
- PUT /v1/user, Create user
- POST /v1/user/:id, Update user
- DELETE /v1/user/:id, Delete user

### Installation
In this example, we will start a web service with gin. As a result, bellow dependencies should be installed.

```
go get github.com/rookie-ninja/rk-boot/gin
go get github.com/rookie-ninja/rk-boot/database/sqlserver
```

### Start SQL Server with docker
We are going to use [sa:MyStrongPass$] as credential of SQL Server.

```
$ docker run --name sql-server -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=MyStrongPass$" -p 1433:1433 -d mcr.microsoft.com/mssql/server
```

### 1.Create boot.yaml
- Create web server with Gin framework at port 8080
- Create SQL Server entry which connects SQL Server at localhost:1433

```yaml
---
gin:
  - name: user-service
    port: 8080
    enabled: true
sqlServer:
  - name: user                        # Required
    enabled: true                     # Required
    locale: "*::*::*::*"              # Required
    addr: "localhost:1433"            # Optional, default: localhost:1433
    user: sa                          # Optional, default: sa
    pass: MyStrongPass$               # Optional, default: pass
    database:
      - name: user-meta               # Required
        autoCreate: true              # Optional, default: false
#        dryRun: true                 # Optional, default: false
#        params: []                   # Optional, default: []
#    logger:
#      zapLogger: zap                   # Optional, default: default logger with STDOUT
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
	"github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-boot/database/sqlserver"
	"github.com/rookie-ninja/rk-boot/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

var userDb *gorm.DB

func main() {
	boot := rkboot.NewBoot()

	boot.Bootstrap(context.TODO())

	// Auto migrate database and init global userDb variable
	userDb = rkbootsqlserver.GetGormDb("user", "user-meta")

	if !userDb.DryRun {
		userDb.AutoMigrate(&User{})
	}

	// Register APIs
	ginEntry := rkbootgin.GetGinEntry("user-service")
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
	CreatedAt time.Time      `yaml:"-" json:"-"`
	UpdatedAt time.Time      `yaml:"-" json:"-"`
	DeletedAt gorm.DeletedAt `yaml:"-" json:"-" gorm:"index"`
}

type User struct {
	Base
	Id   int    `yaml:"id" json:"id" gorm:"primaryKey"`
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
	res := userDb.Where("id = ?", uid).Find(user)

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, res.Error)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func CreateUser(ctx *gin.Context) {
	user := &User{
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
		Name: ctx.Query("name"),
	}

	res := userDb.Where("id = ?", uid).Updates(user)

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, res.Error)
		return
	}

	if res.RowsAffected < 1 {
		ctx.JSON(http.StatusNotFound, "user not found")
		return
	}

	// get user
	userDb.Where("id = ?", uid).Find(user)

	ctx.JSON(http.StatusOK, user)
}

func DeleteUser(ctx *gin.Context) {
	uid, _ := strconv.Atoi(ctx.Param("id"))
	res := userDb.Delete(&User{
		Id: uid,
	})

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, res.Error)
		return
	}

	if res.RowsAffected < 1 {
		ctx.JSON(http.StatusNotFound, "user not found")
		return
	}

	ctx.String(http.StatusOK, "success")
}
```

### 3.Start server

```
$ go run main.go

2022-01-09T04:41:17.010+0800    INFO    Bootstrap sqlServer entry       {"entryName": "user", "sqlServerUser": "sa", "sqlServerAddr": "localhost:1433"}
2022-01-09T04:41:17.010+0800    INFO    creating database user-meta if not exists
2022-01-09T04:41:17.492+0800    INFO    /Users/dongxuny/go/pkg/mod/github.com/rookie-ninja/rk-db/sqlserver@v0.0.1/boot.go:404 SLOW SQL >= 200ms
[366.743ms] [rows:0] 
IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'user-meta')
BEGIN
  CREATE DATABASE [user-meta];
END;

2022-01-09T04:41:17.492+0800    INFO    creating successs or database user-meta exists
2022-01-09T04:41:17.492+0800    INFO    connecting to database user-meta
2022-01-09T04:41:17.510+0800    INFO    connecting to database user-meta success
2022-01-09T04:41:17.510+0800    INFO    boot/gin_entry.go:913   Bootstrap ginEntry      {"eventId": "88026b54-8782-4ab2-99db-0698739db073", "entryName": "user-service"}
------------------------------------------------------------------------
endTime=2022-01-09T04:41:17.510765+08:00
startTime=2022-01-09T04:41:17.510717+08:00
elapsedNano=48442
timezone=CST
ids={"eventId":"88026b54-8782-4ab2-99db-0698739db073"}
app={"appName":"rk","appVersion":"","entryName":"user-service","entryType":"GinEntry"}
env={"arch":"amd64","az":"*","domain":"*","hostname":"lark.local","localIP":"10.8.0.2","os":"darwin","realm":"*","region":"*"}
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
```

### 4.Validation
#### 4.1 Create user
Create a user with name of rk-dev.

```shell script
$ curl -X PUT "localhost:8080/v1/user?name=rk-dev"
{"id":1,"name":"rk-dev"}
```

#### 4.1 Update user
Update user name to rk-dev-updated.

```shell script
$ curl -X POST "localhost:8080/v1/user/1?name=rk-dev-updated"
{"id":1,"name":"rk-dev-updated"}
```

#### 4.1 List users
List users.

```shell script
$ curl -X GET localhost:8080/v1/user
[{"id":1,"name":"rk-dev-updated"}]
```

#### 4.1 Get user
Get user with id=1.

```shell script
$ curl -X GET localhost:8080/v1/user/1
{"id":1,"name":"rk-dev-updated"}
```

#### 4.1 Delete user

```shell script
$ curl -X DELETE localhost:8080/v1/user/1
success
```

![image](img/sqlserver.png)

## YAML Options
User can start multiple [gorm](https://github.com/go-gorm/gorm) instances at the same time. Please make sure use different names.

| name | Required | description | type | default value |
| ------ | ------ | ------ | ------ | ------ |
| sqlServer.name | Required | The name of entry | string | SqlServer |
| sqlServer.enabled | Required | Enable entry or not | bool | false |
| sqlServer.locale | Required | See locale description bellow | string | "" |
| sqlServer.description | Optional | Description of echo entry. | string | "" |
| sqlServer.user | Optional | SQL Server username | string | sa |
| sqlServer.pass | Optional | SQL Server password | string | pass |
| sqlServer.addr | Optional | SQL Server remote address | string | localhost:1433 |
| sqlServer.database.name | Required | Name of database | string | "" |
| sqlServer.database.autoCreate | Optional | Create DB if missing | bool | false |
| sqlServer.database.dryRun | Optional | Run gorm.DB with dry run mode | bool | false |
| sqlServer.database.params | Optional | Connection params | []string | [] |
| sqlServer.logger.zapLogger | Optional | Reference of zap logger entry name | string | "" |

### Usage of locale

```
RK use <realm>::<region>::<az>::<domain> to distinguish different environment.
Variable of <locale> could be composed as form of <realm>::<region>::<az>::<domain>
- realm: It could be a company, department and so on, like RK-Corp.
         Environment variable: REALM
         Eg: RK-Corp
         Wildcard: supported

- region: Please see AWS web site: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html
          Environment variable: REGION
          Eg: us-east
          Wildcard: supported

- az: Availability zone, please see AWS web site for details: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html
      Environment variable: AZ
      Eg: us-east-1
      Wildcard: supported

- domain: Stands for different environment, like dev, test, prod and so on, users can define it by themselves.
          Environment variable: DOMAIN
          Eg: prod
          Wildcard: supported

How it works?
First, we will split locale with "::" and extract realm, region, az and domain.
Second, get environment variable named as REALM, REGION, AZ and DOMAIN.
Finally, compare every element in locale variable and environment variable.
If variables in locale represented as wildcard(*), we will ignore comparison step.

Example:
# let's assuming we are going to define DB address which is different based on environment.
# Then, user can distinguish DB address based on locale.
# We recommend to include locale with wildcard.
---
DB:
  - name: redis-default
    locale: "*::*::*::*"
    addr: "192.0.0.1:6379"
  - name: redis-in-test
    locale: "*::*::*::test"
    addr: "192.0.0.1:6379"
  - name: redis-in-prod
    locale: "*::*::*::prod"
    addr: "176.0.0.1:6379"
```
