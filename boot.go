// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

// Package rkboot is bootstrapper for rk style application
package rkboot

import (
	"context"
	"embed"
	"fmt"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-entry/entry"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
	"runtime/debug"
)

// Boot is a structure for bootstrapping rk style application
type Boot struct {
	BootConfigPath string `yaml:"bootConfigPath" json:"bootConfigPath"`
	bootConfigRaw  []byte `yaml:"-" json:"-"`
	EventId        string `yaml:"eventId" json:"eventId"`
}

// BootOption is used as options while bootstrapping from code
type BootOption func(*Boot)

// WithBootConfigPath provide boot config yaml file.
func WithBootConfigPath(filePath string) BootOption {
	return func(boot *Boot) {
		boot.BootConfigPath = filePath
	}
}

// WithBootConfigString provide boot config as string.
func WithBootConfigString(bootConfigStr string) BootOption {
	return func(boot *Boot) {
		if len(bootConfigStr) > 0 {
			boot.bootConfigRaw = []byte(bootConfigStr)
		}
	}
}

// WithBootConfigBytes provide boot config as string.
func WithBootConfigBytes(bootConfigBytes []byte) BootOption {
	return func(boot *Boot) {
		if len(bootConfigBytes) > 0 {
			boot.bootConfigRaw = bootConfigBytes
		}
	}
}

// WithBootConfigEmbedFs provide boot config as file in embed.FS.
func WithBootConfigEmbedFs(fs embed.FS, filePath string) BootOption {
	return func(boot *Boot) {
		bytes, err := fs.ReadFile(filePath)
		if err != nil {
			rkcommon.ShutdownWithError(err)
		}

		boot.bootConfigRaw = bytes
	}
}

// NewBoot create a bootstrapper.
func NewBoot(opts ...BootOption) *Boot {
	defer syncLog("N/A")

	boot := &Boot{
		EventId: rkcommon.GenerateRequestId(),
	}

	for i := range opts {
		opts[i](boot)
	}

	if len(boot.bootConfigRaw) > 0 {
		wd, err := os.Getwd()
		if err != nil {
			rkcommon.ShutdownWithError(err)
		}

		boot.BootConfigPath = path.Join(wd, "boot-gen.yaml")
		err = ioutil.WriteFile(boot.BootConfigPath, boot.bootConfigRaw, os.ModePerm)
		if err != nil {
			rkcommon.ShutdownWithError(err)
		}
	}

	if len(boot.BootConfigPath) < 1 {
		boot.BootConfigPath = "boot.yaml"
	}

	// Register and bootstrap internal entries with boot config.
	rkentry.RegisterInternalEntriesFromConfig(boot.BootConfigPath)

	// Register external entries.
	regFuncList := rkentry.ListEntryRegFunc()
	for i := range regFuncList {
		regFuncList[i](boot.BootConfigPath)
	}

	return boot
}

// Bootstrap entries in rkentry.GlobalAppCtx including bellow:
//
// Internal entries:
// 1: rkentry.AppInfoEntry
// 2: rkentry.ConfigEntry
// 3: rkentry.ZapLoggerEntry
// 4: rkentry.EventLoggerEntry
// 5: rkentry.CertEntry
//
// External entries:
// User defined entries
func (boot *Boot) Bootstrap(ctx context.Context) {
	defer syncLog(boot.EventId)

	ctx = context.WithValue(ctx, "eventId", boot.EventId)

	// Bootstrap external entries
	for _, entry := range rkentry.GlobalAppCtx.ListEntries() {
		entry.Bootstrap(ctx)
	}
}

// WaitForShutdownSig wait for shutdown signal.
// 1: Call shutdown hook function added by user.
// 2: Call interrupt function of entries in rkentry.GlobalAppCtx.
func (boot *Boot) WaitForShutdownSig(ctx context.Context) {
	rkentry.GlobalAppCtx.WaitForShutdownSig()

	// Call shutdown hook function
	for _, f := range rkentry.GlobalAppCtx.ListShutdownHooks() {
		f()
	}

	// Call interrupt
	boot.interrupt(ctx)
}

// AddShutdownHookFunc add shutdown hook function
func (boot *Boot) AddShutdownHookFunc(name string, f rkentry.ShutdownHook) {
	rkentry.GlobalAppCtx.AddShutdownHook(name, f)
}

// Interrupt entries in rkentry.GlobalAppCtx including bellow:
//
// Internal entries:
// 1: rkentry.AppInfoEntry
// 2: rkentry.ConfigEntry
// 3: rkentry.ZapLoggerEntry
// 4: rkentry.EventLoggerEntry
// 5: rkentry.CertEntry
//
// External entries:
// User defined entries
func (boot *Boot) interrupt(ctx context.Context) {
	defer syncLog(boot.EventId)

	ctx = context.WithValue(ctx, "eventId", boot.EventId)

	// Interrupt external entries
	boot.interruptHelper(ctx, rkentry.GlobalAppCtx.ListEntries())

	// Interrupt internal entries
	rkentry.GlobalAppCtx.GetAppInfoEntry().Interrupt(ctx)
	boot.interruptHelper(ctx, rkentry.GlobalAppCtx.ListConfigEntriesRaw())
	boot.interruptHelper(ctx, rkentry.GlobalAppCtx.ListCertEntriesRaw())
	boot.interruptHelper(ctx, rkentry.GlobalAppCtx.ListEventLoggerEntriesRaw())
	boot.interruptHelper(ctx, rkentry.GlobalAppCtx.ListZapLoggerEntriesRaw())
}

// Helper function which all interrupt() function for each entry.
func (boot *Boot) interruptHelper(ctx context.Context, m map[string]rkentry.Entry) {
	for _, entry := range m {
		entry.Interrupt(ctx)
	}
}

// GetAppInfoEntry returns rkentry.AppInfoEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetAppInfoEntry() *rkentry.AppInfoEntry {
	return rkentry.GlobalAppCtx.GetAppInfoEntry()
}

// GetZapLoggerEntry returns rkentry.ZapLoggerEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetZapLoggerEntry(name string) *rkentry.ZapLoggerEntry {
	return rkentry.GlobalAppCtx.GetZapLoggerEntry(name)
}

// GetZapLoggerEntryDefault returns default rkentry.ZapLoggerEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetZapLoggerEntryDefault() *rkentry.ZapLoggerEntry {
	return rkentry.GlobalAppCtx.GetZapLoggerEntryDefault()
}

// GetEventLoggerEntry returns rkentry.EventLoggerEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetEventLoggerEntry(name string) *rkentry.EventLoggerEntry {
	return rkentry.GlobalAppCtx.GetEventLoggerEntry(name)
}

// GetEventLoggerEntryDefault returns default rkentry.EventLoggerEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetEventLoggerEntryDefault() *rkentry.EventLoggerEntry {
	return rkentry.GlobalAppCtx.GetEventLoggerEntryDefault()
}

// GetConfigEntry returns rkentry.ConfigEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetConfigEntry(name string) *rkentry.ConfigEntry {
	return rkentry.GlobalAppCtx.GetConfigEntry(name)
}

// GetCertEntry returns rkentry.CertEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetCertEntry(name string) *rkentry.CertEntry {
	return rkentry.GlobalAppCtx.GetCertEntry(name)
}

// GetEntry returns rkentry.Entry interface which user needs to convert by himself.
func (boot *Boot) GetEntry(name string) interface{} {
	return rkentry.GlobalAppCtx.GetEntry(name)
}

// sync logs
func syncLog(eventId string) {
	if r := recover(); r != nil {
		stackTrace := fmt.Sprintf("Panic occured, shutting down... \n%s", string(debug.Stack()))
		for _, v := range rkentry.GlobalAppCtx.ListZapLoggerEntries() {
			if v == rkentry.GlobalAppCtx.GetZapLoggerEntryDefault() {
				continue
			}
			if v.Logger != nil {
				v.Logger.Error(stackTrace,
					zap.String("eventId", eventId),
					zap.Any("RootCause", r))
			}
			v.Sync()
		}

		for _, v := range rkentry.GlobalAppCtx.ListEventLoggerEntries() {
			v.Sync()
		}

		panic(r)
	}
}
