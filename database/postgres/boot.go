// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkbootpostgres

import (
	_ "github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-db/postgres"
	"gorm.io/gorm"
)

// GetPostgresEntry return rkpostgres.PostgresEntry with name.
//
// The entryName was specified in boot.yaml file as bellow.
//
// user is the name of entry as function parameter.
//
// postgres:
//  - name: user                        # Required
//    enabled: true                     # Required
//    locale: "*::*::*::*"              # Required
func GetPostgresEntry(entryName string) *rkpostgres.PostgresEntry {
	return rkpostgres.GetPostgresEntry(entryName)
}

// GetGormDb return gorm.DB instance with entryName and database name.
//
// rkpostgres will init gorm.DB by reading boot.yaml file and validate the connectivity for us.
//
// The entryName and dataBaseName was specified in boot.yaml file as bellow.
//
// entryName: user
// dataBaseName: user-meta
//
// postgres:
//  - name: user                        # Required
//    enabled: true                     # Required
//    locale: "*::*::*::*"              # Required
//    addr: "localhost:5432"            # Optional, default: localhost:5432
//    user: postgres                    # Optional, default: postgres
//    pass: pass                        # Optional, default: pass
//    database:
//      - name: user-meta               # Required
//        autoCreate: true              # Optional, default: false
func GetGormDb(entryName, dataBaseName string) *gorm.DB {
	if entry := GetPostgresEntry(entryName); entry != nil {
		return entry.GetDB(dataBaseName)
	}

	return nil
}
