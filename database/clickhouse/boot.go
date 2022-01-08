// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkbootclickhouse

import (
	_ "github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-db/clickhouse"
	"gorm.io/gorm"
)

// GetClickHouseEntry return rkclickhouse.ClickHouseEntry with name.
//
// The entryName was specified in boot.yaml file as bellow.
//
// user is the name of entry as function parameter.
//
// clickHouse:
//  - name: user                        # Required
//    enabled: true                     # Required
//    locale: "*::*::*::*"              # Required
func GetClickHouseEntry(entryName string) *rkclickhouse.ClickHouseEntry {
	return rkclickhouse.GetClickHouseEntry(entryName)
}

// GetGormDb return gorm.DB instance with entryName and database name.
//
// rkclickhouse will init gorm.DB by reading boot.yaml file and validate the connectivity for us.
//
// The entryName and dataBaseName was specified in boot.yaml file as bellow.
//
// entryName: user
// dataBaseName: usermeta
//
// clickHouse:
//  - name: user                        # Required
//    enabled: true                     # Required
//    locale: "*::*::*::*"              # Required
//    addr: "localhost:9000"            # Optional, default: localhost:9000
//    user: default                     # Optional, default: default
//    pass: ""                          # Optional, default: ""
//    database:
//      - name: usermeta                # Required
//        autoCreate: true              # Optional, default: false
func GetGormDb(entryName, dataBaseName string) *gorm.DB {
	if entry := GetClickHouseEntry(entryName); entry != nil {
		return entry.GetDB(dataBaseName)
	}

	return nil
}
