// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkbootsqlite

import (
	_ "github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-db/sqlite"
	"gorm.io/gorm"
)

// GetSqliteEntry return rksqlite.SqliteEntry with name.
//
// The entryName was specified in boot.yaml file as bellow.
//
// user is the name of entry as function parameter.
//
// sqlite:
//  - name: user                        # Required
//    enabled: true                     # Required
//    locale: "*::*::*::*"              # Required
func GetSqliteEntry(entryName string) *rksqlite.SqliteEntry {
	return rksqlite.GetSqliteEntry(entryName)
}

// GetGormDb return gorm.DB instance with entryName and database name.
//
// rksqlite will init gorm.DB by reading boot.yaml file and validate the connectivity for us.
//
// The entryName and dataBaseName was specified in boot.yaml file as bellow.
//
// entryName: user
// dataBaseName: user-meta
//
// sqlite:
//  - name: user                        # Required
//    enabled: true                     # Required
//    locale: "*::*::*::*"              # Required
//    database:
//      - name: user-meta               # Required
//        inMemory: true                # Optional, default: false
func GetGormDb(entryName, dataBaseName string) *gorm.DB {
	if entry := GetSqliteEntry(entryName); entry != nil {
		return entry.GetDB(dataBaseName)
	}

	return nil
}
