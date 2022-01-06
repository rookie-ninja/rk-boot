// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkbootmysql

import (
	_ "github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-db/mysql"
	"gorm.io/gorm"
)

// GetMySqlEntry return rkmysql.MySqlEntry with name.
//
// The entryName was specified in boot.yaml file as bellow.
//
// user is the name of entry as function parameter.
//
// mySql:
//  - name: user                        # Required
//    enabled: true                     # Required
//    locale: "*::*::*::*"              # Required
//    addr: "localhost:3306"            # Optional, default: localhost:3306
func GetMySqlEntry(entryName string) *rkmysql.MySqlEntry {
	return rkmysql.GetMySqlEntry(entryName)
}

// GetGormDb return gorm.DB instance with entryName and database name.
//
// rkmysql will init gorm.DB by reading boot.yaml file and validate the connectivity for us.
//
// The entryName and dataBaseName was specified in boot.yaml file as bellow.
//
// entryName: user
// dataBaseName: user-meta
//
// mySql:
//  - name: user                        # Required
//    enabled: true                     # Required
//    locale: "*::*::*::*"              # Required
//    addr: "localhost:3306"            # Optional, default: localhost:3306
//    database:
//      - name: user-meta               # Required
//        autoCreate: true              # Optional, default: false
func GetGormDb(entryName, dataBaseName string) *gorm.DB {
	if entry := GetMySqlEntry(entryName); entry != nil {
		return entry.GetDB(dataBaseName)
	}

	return nil
}
