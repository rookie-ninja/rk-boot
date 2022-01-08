// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkbootsqlserver

import (
	_ "github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-db/sqlserver"
	"gorm.io/gorm"
)

// GetSqlserverEntry return rksqlserver.SqlServerEntry with name.
//
// The entryName was specified in boot.yaml file as bellow.
//
// user is the name of entry as function parameter.
//
// sqlite:
//  - name: user                        # Required
//    enabled: true                     # Required
//    locale: "*::*::*::*"              # Required
func GetSqlserverEntry(entryName string) *rksqlserver.SqlServerEntry {
	return rksqlserver.GetSqlServerEntry(entryName)
}

// GetGormDb return gorm.DB instance with entryName and database name.
//
// rksqlserver will init gorm.DB by reading boot.yaml file and validate the connectivity for us.
//
// The entryName and dataBaseName was specified in boot.yaml file as bellow.
//
// entryName: user
// dataBaseName: user-meta
//
// sqlServer:
//  - name: user                        # Required
//    enabled: true                     # Required
//    locale: "*::*::*::*"              # Required
//    addr: "localhost:1433"            # Optional, default: localhost:1433
//    user: sa                          # Optional, default: sa
//    pass: pass                        # Optional, default: pass
//    database:
//      - name: user-meta               # Required
//        autoCreate: true              # Optional, default: false
func GetGormDb(entryName, dataBaseName string) *gorm.DB {
	if entry := GetSqlserverEntry(entryName); entry != nil {
		return entry.GetDB(dataBaseName)
	}

	return nil
}
