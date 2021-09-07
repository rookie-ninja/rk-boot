// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkboot

import (
	"context"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"syscall"
	"testing"
	"time"
)

func TestNewBoot_HappyCase(t *testing.T) {
	config := `
---
gin:
  - name: ut-gin
    port: 8080
`

	filePath := createFileAtTestTempDir(t, config)

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
	assert.NotNil(t, boot.GetGinEntry("ut-gin"))
	assert.Nil(t, boot.GetGrpcEntry(""))
	assert.Nil(t, boot.GetPromEntry(""))

	go func() {
		time.Sleep(3 * time.Second)
		rkentry.GlobalAppCtx.GetShutdownSig() <- syscall.SIGTERM
	}()

	boot.WaitForShutdownSig(context.TODO())
}

func createFileAtTestTempDir(t *testing.T, content string) string {
	tempDir := path.Join(t.TempDir(), "ut-boot.yaml")
	assert.Nil(t, ioutil.WriteFile(tempDir, []byte(content), os.ModePerm))
	return tempDir
}
