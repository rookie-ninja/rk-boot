// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkboot

import (
	"context"
	"fmt"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"syscall"
	"testing"
	"time"
)

func TestNewBoot_HappyCase_Gin(t *testing.T) {
	config := `
---
gin:
  - name: ut-gin
    port: 8080
    enabled: true
`

	filePath := createFileAtTestTempDir(t, "ut-boot-gin.yaml", config)

	boot := NewBoot(WithBootConfigPath(filePath))
	boot.AddShutdownHookFunc("ut-shutdown", func() {
		// noop
	})

	boot.Bootstrap(context.TODO())

	assert.NotNil(t, boot.GetAppInfoEntry())
	assert.NotNil(t, boot.GetZapLoggerEntry(rkentry.DefaultZapLoggerEntryName))
	assert.NotNil(t, boot.GetZapLoggerEntryDefault())
	assert.NotNil(t, boot.GetEventLoggerEntry(rkentry.DefaultEventLoggerEntryName))
	assert.NotNil(t, boot.GetEventLoggerEntryDefault())
	assert.Nil(t, boot.GetConfigEntry(""))
	assert.Nil(t, boot.GetCertEntry(""))
	assert.Nil(t, boot.GetPromEntry(""))

	go func() {
		boot.WaitForShutdownSig(context.TODO())
	}()

	time.Sleep(1 * time.Second)
	rkentry.GlobalAppCtx.GetShutdownSig() <- syscall.SIGTERM
	time.Sleep(1 * time.Second)
	rkentry.GlobalAppCtx.RemoveEntry("ut-gin")
}

func TestNewBoot_HappyCase_Echo(t *testing.T) {
	config := `
---
echo:
  - name: ut-echo
    port: 8080
    enabled: true
`

	filePath := createFileAtTestTempDir(t, "ut-boot-echo.yaml", config)

	boot := NewBoot(WithBootConfigPath(filePath))
	boot.AddShutdownHookFunc("ut-shutdown", func() {
		// noop
	})

	boot.Bootstrap(context.TODO())

	assert.NotNil(t, boot.GetAppInfoEntry())
	assert.NotNil(t, boot.GetZapLoggerEntry(rkentry.DefaultZapLoggerEntryName))
	assert.NotNil(t, boot.GetZapLoggerEntryDefault())
	assert.NotNil(t, boot.GetEventLoggerEntry(rkentry.DefaultEventLoggerEntryName))
	assert.NotNil(t, boot.GetEventLoggerEntryDefault())
	assert.Nil(t, boot.GetConfigEntry(""))
	assert.Nil(t, boot.GetCertEntry(""))
	assert.Nil(t, boot.GetPromEntry(""))

	go func() {
		boot.WaitForShutdownSig(context.TODO())
	}()

	time.Sleep(1 * time.Second)
	rkentry.GlobalAppCtx.GetShutdownSig() <- syscall.SIGTERM
	time.Sleep(1 * time.Second)
	rkentry.GlobalAppCtx.RemoveEntry("ut-echo")
}

func TestNewBoot_HappyCase_Gf(t *testing.T) {
	config := `
---
gf:
  - name: ut-gf
    port: 8080
    enabled: true
`

	filePath := createFileAtTestTempDir(t, "ut-boot-gf.yaml", config)

	boot := NewBoot(WithBootConfigPath(filePath))
	boot.AddShutdownHookFunc("ut-shutdown", func() {
		// noop
	})

	boot.Bootstrap(context.TODO())

	assert.NotNil(t, boot.GetAppInfoEntry())
	assert.NotNil(t, boot.GetZapLoggerEntry(rkentry.DefaultZapLoggerEntryName))
	assert.NotNil(t, boot.GetZapLoggerEntryDefault())
	assert.NotNil(t, boot.GetEventLoggerEntry(rkentry.DefaultEventLoggerEntryName))
	assert.NotNil(t, boot.GetEventLoggerEntryDefault())
	assert.Nil(t, boot.GetConfigEntry(""))
	assert.Nil(t, boot.GetCertEntry(""))
	assert.Nil(t, boot.GetPromEntry(""))

	go func() {
		boot.WaitForShutdownSig(context.TODO())
	}()

	time.Sleep(1 * time.Second)
	rkentry.GlobalAppCtx.GetShutdownSig() <- syscall.SIGTERM
	time.Sleep(1 * time.Second)
	rkentry.GlobalAppCtx.RemoveEntry("ut-gf")
}

func TestNewBoot_HappyCase_Grpc(t *testing.T) {
	config := `
---
grpc:
  - name: ut-grpc
    port: 8080
    enabled: true
`

	filePath := createFileAtTestTempDir(t, "ut-boot-grpc.yaml", config)

	boot := NewBoot(WithBootConfigPath(filePath))
	fmt.Println(filePath)
	boot.AddShutdownHookFunc("ut-shutdown", func() {
		// noop
	})

	boot.Bootstrap(context.TODO())

	assert.NotNil(t, boot.GetAppInfoEntry())
	assert.NotNil(t, boot.GetZapLoggerEntry(rkentry.DefaultZapLoggerEntryName))
	assert.NotNil(t, boot.GetZapLoggerEntryDefault())
	assert.NotNil(t, boot.GetEventLoggerEntry(rkentry.DefaultEventLoggerEntryName))
	assert.NotNil(t, boot.GetEventLoggerEntryDefault())
	assert.Nil(t, boot.GetConfigEntry(""))
	assert.Nil(t, boot.GetCertEntry(""))
	assert.Nil(t, boot.GetPromEntry(""))

	go func() {
		boot.WaitForShutdownSig(context.TODO())
	}()

	time.Sleep(1 * time.Second)
	rkentry.GlobalAppCtx.GetShutdownSig() <- syscall.SIGTERM
	time.Sleep(1 * time.Second)
	rkentry.GlobalAppCtx.RemoveEntry("ut-grpc")

}

func createFileAtTestTempDir(t *testing.T, filename, content string) string {
	tempDir := path.Join(t.TempDir(), filename)

	assert.Nil(t, ioutil.WriteFile(tempDir, []byte(content), os.ModePerm))
	return tempDir
}
