// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkboot

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/rookie-ninja/rk-entry/v2/entry"
	"github.com/stretchr/testify/assert"
	"syscall"
	"testing"
	"time"
)

func TestNewBoot_WithRawCase(t *testing.T) {
	config := `
---
myEntry:
  name: ut
  enabled: true
`

	triggerBefore := false
	triggerAfter := false

	boot := NewBoot(WithBootConfigRaw([]byte(config)))
	boot.AddShutdownHookFunc("ut-shutdown", func() {
		// noop
	})
	boot.AddHookFuncBeforeBootstrap("myEntry", "ut", func(ctx context.Context) {
		triggerBefore = true
	})
	boot.AddHookFuncAfterBootstrap("myEntry", "ut", func(ctx context.Context) {
		triggerAfter = true
	})

	boot.Bootstrap(context.TODO())

	assert.Len(t, rkentry.GlobalAppCtx.ListEntriesByType("myEntry"), 1)
	assert.True(t, triggerBefore)
	assert.True(t, triggerAfter)

	go func() {
		boot.WaitForShutdownSig(context.TODO())
	}()

	time.Sleep(1 * time.Second)
	rkentry.GlobalAppCtx.GetShutdownSig() <- syscall.SIGTERM
	time.Sleep(1 * time.Second)
	rkentry.GlobalAppCtx.RemoveEntry(rkentry.GlobalAppCtx.GetEntry("myEntry", "ut"))
}

//go:embed testdata/boot.yaml
var embedFS embed.FS

func TestNewBoot_WithEmbedCase(t *testing.T) {
	boot := NewBoot(WithBootConfigPath("testdata/boot.yaml", &embedFS))
	myEntry := rkentry.GlobalAppCtx.GetEntry("myEntry", "ut")
	assert.NotNil(t, myEntry)

	boot.Bootstrap(context.TODO())
	boot.interrupt(context.TODO())

	rkentry.GlobalAppCtx.RemoveEntry(rkentry.GlobalAppCtx.GetEntry("myEntry", "ut"))
}

func assertPanic(t *testing.T) {
	if r := recover(); r != nil {
		fmt.Println("adsfadfafd")
		// Expect panic to be called with non nil error
		assert.True(t, true)
	} else {
		// This should never be called in case of a bug
		assert.True(t, false)
	}
}

// Register entry, must be in init() function since we need to register entry at beginning
func init() {
	rkentry.RegisterUserEntryRegFunc(RegisterMyEntriesFromConfig)
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
func RegisterMyEntriesFromConfig(raw []byte) map[string]rkentry.Entry {
	res := make(map[string]rkentry.Entry)

	// 1: decode config map into boot config struct
	config := &BootConfig{}
	rkentry.UnmarshalBootYAML(raw, config)

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
