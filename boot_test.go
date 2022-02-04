// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkboot

import (
	"context"
	"encoding/json"
	rkcommon "github.com/rookie-ninja/rk-common/common"
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
myEntry:
  name: ut
  enabled: true
`

	filePath := createFileAtTestTempDir(t, "ut-boot.yaml", config)

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
	assert.NotNil(t, boot.GetEntry("ut"))

	go func() {
		boot.WaitForShutdownSig(context.TODO())
	}()

	time.Sleep(1 * time.Second)
	rkentry.GlobalAppCtx.GetShutdownSig() <- syscall.SIGTERM
	time.Sleep(1 * time.Second)
	rkentry.GlobalAppCtx.RemoveEntry("ut")
}

func TestNewBoot_WithEmbedCase(t *testing.T) {
	config := `
---
myEntry:
  name: ut
  enabled: true
`
	// provide boot config as string
	boot := NewBoot(WithBootConfigString(config))
	wd, _ := os.Getwd()
	bytes, _ := ioutil.ReadFile(path.Join(wd, "boot-gen.yaml"))
	assert.Equal(t, config, string(bytes))

	boot.Bootstrap(context.TODO())
	boot.interrupt(context.TODO())
	rkentry.GlobalAppCtx.RemoveEntry("ut")
	os.Remove(path.Join(wd, "boot-gen.yaml"))

	// provide boot config as byte
	boot = NewBoot(WithBootConfigBytes([]byte(config)))
	wd, _ = os.Getwd()
	bytes, _ = ioutil.ReadFile(path.Join(wd, "boot-gen.yaml"))
	assert.Equal(t, config, string(bytes))

	boot.Bootstrap(context.TODO())
	boot.interrupt(context.TODO())
	rkentry.GlobalAppCtx.RemoveEntry("ut")
	os.Remove(path.Join(wd, "boot-gen.yaml"))
}

func TestNewBoot_Panic(t *testing.T) {
	defer assertPanic(t)
	defer rkentry.GlobalAppCtx.RemoveEntry("ut")

	config := `
---
zapLogger:
  - name: zap
myEntry:
  name: ut
  enabled: true
  shouldPanic: true
`

	filePath := createFileAtTestTempDir(t, "ut-boot.yaml", config)

	boot := NewBoot(WithBootConfigPath(filePath))
	boot.Bootstrap(context.TODO())
}

func TestNewBoot_EmptyConfig(t *testing.T) {
	defer assertPanic(t)

	NewBoot()
}

func TestBoot_GetEntry(t *testing.T) {
	config := `
---
gin:
  - name: ut-gin
    port: 8080
    enabled: true
`

	filePath := createFileAtTestTempDir(t, "ut-boot-gin.yaml", config)

	boot := NewBoot(WithBootConfigPath(filePath))
	boot.Bootstrap(context.TODO())

	assert.Nil(t, boot.GetEntry(rkentry.GlobalAppCtx.GetAppInfoEntry().GetName()))
}

func createFileAtTestTempDir(t *testing.T, filename, content string) string {
	tempDir := path.Join(t.TempDir(), filename)

	assert.Nil(t, ioutil.WriteFile(tempDir, []byte(content), os.ModePerm))
	return tempDir
}

func assertPanic(t *testing.T) {
	if r := recover(); r != nil {
		// Expect panic to be called with non nil error
		assert.True(t, true)
	} else {
		// This should never be called in case of a bug
		assert.True(t, false)
	}
}

// Register entry, must be in init() function since we need to register entry at beginning
func init() {
	rkentry.RegisterEntryRegFunc(RegisterMyEntriesFromConfig)
}

// A struct which is for unmarshalled YAML
type BootConfig struct {
	MyEntry struct {
		Enabled     bool   `yaml:"enabled" json:"enabled"`
		Name        string `yaml:"name" json:"name"`
		Description string `yaml:"description" json:"description"`
		ShouldPanic bool   `yaml:"shouldPanic" json:"shouldPanic"`
	} `yaml:"myEntry" json:"myEntry"`
}

// An implementation of:
// type EntryRegFunc func(string) map[string]rkentry.Entry
func RegisterMyEntriesFromConfig(configFilePath string) map[string]rkentry.Entry {
	res := make(map[string]rkentry.Entry)

	// 1: decode config map into boot config struct
	config := &BootConfig{}
	rkcommon.UnmarshalBootConfig(configFilePath, config)

	// 3: construct entry
	if config.MyEntry.Enabled {
		entry := RegisterMyEntry(
			WithName(config.MyEntry.Name),
			WithDescription(config.MyEntry.Description),
			WithPanic(config.MyEntry.ShouldPanic))
		res[entry.GetName()] = entry
	}

	return res
}

func RegisterMyEntry(opts ...MyEntryOption) *MyEntry {
	entry := &MyEntry{
		EntryName:        "default",
		EntryType:        "myEntry",
		EntryDescription: "Please contact maintainers to add description of this entry.",
	}

	for i := range opts {
		opts[i](entry)
	}

	if len(entry.EntryName) < 1 {
		entry.EntryName = "my-default"
	}

	if len(entry.EntryDescription) < 1 {
		entry.EntryDescription = "Please contact maintainers to add description of this entry."
	}

	rkentry.GlobalAppCtx.AddEntry(entry)

	return entry
}

type MyEntryOption func(*MyEntry)

func WithName(name string) MyEntryOption {
	return func(entry *MyEntry) {
		entry.EntryName = name
	}
}

func WithDescription(description string) MyEntryOption {
	return func(entry *MyEntry) {
		entry.EntryDescription = description
	}
}

func WithPanic(val bool) MyEntryOption {
	return func(entry *MyEntry) {
		entry.shouldPanic = val
	}
}

type MyEntry struct {
	EntryName        string `json:"entryName" yaml:"entryName"`
	EntryType        string `json:"entryType" yaml:"entryType"`
	EntryDescription string `json:"entryDescription" yaml:"entryDescription"`
	shouldPanic      bool
}

func (entry *MyEntry) Bootstrap(context.Context) {
	if entry.shouldPanic {
		panic("expected panic")
	}
}

func (entry *MyEntry) Interrupt(context.Context) {}

func (entry *MyEntry) GetName() string {
	return entry.EntryName
}

func (entry *MyEntry) GetDescription() string {
	return entry.EntryDescription
}

func (entry *MyEntry) GetType() string {
	return entry.EntryType
}

func (entry *MyEntry) String() string {
	bytes, _ := json.Marshal(entry)
	return string(bytes)
}
