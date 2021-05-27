// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rkboot

import (
	"context"
	"github.com/rookie-ninja/rk-entry/entry"
	rkgin "github.com/rookie-ninja/rk-gin/boot"
	rkgrpc "github.com/rookie-ninja/rk-grpc/boot"
	rkprom "github.com/rookie-ninja/rk-prom"
)

type Boot struct {
	BootConfigPath string `yaml:"bootConfigPath" json:"bootConfigPath"`
}

type BootOption func(*Boot)

// Provide boot config yaml file.
func WithBootConfigPath(filePath string) BootOption {
	return func(boot *Boot) {
		boot.BootConfigPath = filePath
	}
}

// Create a bootstrapper.
func NewBoot(opts ...BootOption) *Boot {
	boot := &Boot{}

	for i := range opts {
		opts[i](boot)
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
	// Bootstrap external entries
	for _, entry := range rkentry.GlobalAppCtx.ListEntries() {
		entry.Bootstrap(ctx)
	}
}

// Wait for shutdown signal
func (boot *Boot) WaitForShutdownSig() {
	rkentry.GlobalAppCtx.WaitForShutdownSig()
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
func (boot *Boot) Interrupt(ctx context.Context) {
	// Interrupt internal entries
	rkentry.GlobalAppCtx.GetAppInfoEntry().Interrupt(ctx)
	boot.interruptHelper(ctx, rkentry.GlobalAppCtx.ListZapLoggerEntriesRaw())
	boot.interruptHelper(ctx, rkentry.GlobalAppCtx.ListEventLoggerEntriesRaw())
	boot.interruptHelper(ctx, rkentry.GlobalAppCtx.ListConfigEntriesRaw())
	boot.interruptHelper(ctx, rkentry.GlobalAppCtx.ListCertEntriesRaw())

	// Interrupt external entries
	boot.interruptHelper(ctx, rkentry.GlobalAppCtx.ListEntries())
}

// Helper function which all interrupt() function for each entry.
func (boot *Boot) interruptHelper(ctx context.Context, m map[string]rkentry.Entry) {
	for _, entry := range m {
		entry.Interrupt(ctx)
	}
}

// Get rkentry.AppInfoEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetAppInfoEntry() *rkentry.AppInfoEntry {
	return rkentry.GlobalAppCtx.GetAppInfoEntry()
}

// Get rkentry.ZapLoggerEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetZapLoggerEntry(name string) *rkentry.ZapLoggerEntry {
	return rkentry.GlobalAppCtx.GetZapLoggerEntry(name)
}

// Get default rkentry.ZapLoggerEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetZapLoggerEntryDefault() *rkentry.ZapLoggerEntry {
	return rkentry.GlobalAppCtx.GetZapLoggerEntryDefault()
}

// Get rkentry.EventLoggerEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetEventLoggerEntry(name string) *rkentry.EventLoggerEntry {
	return rkentry.GlobalAppCtx.GetEventLoggerEntry(name)
}

// Get default rkentry.EventLoggerEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetEventLoggerEntryDefault() *rkentry.EventLoggerEntry {
	return rkentry.GlobalAppCtx.GetEventLoggerEntryDefault()
}

// Get rkentry.ConfigEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetConfigEntry(name string) *rkentry.ConfigEntry {
	return rkentry.GlobalAppCtx.GetConfigEntry(name)
}

// Get rkentry.CertEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetCertEntry(name string) *rkentry.CertEntry {
	return rkentry.GlobalAppCtx.GetCertEntry(name)
}

// Get rkgin.GinEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetGinEntry(name string) *rkgin.GinEntry {
	entryRaw := rkentry.GlobalAppCtx.GetEntry(name)

	if entry, ok := entryRaw.(*rkgin.GinEntry); ok {
		return entry
	}

	return nil
}

// Get rkgrpc.GrpcEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetGrpcEntry(name string) *rkgrpc.GrpcEntry {
	entryRaw := rkentry.GlobalAppCtx.GetEntry(name)

	if entry, ok := entryRaw.(*rkgrpc.GrpcEntry); ok {
		return entry
	}

	return nil
}

// Get rkprom.PromEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetPromEntry(name string) *rkprom.PromEntry {
	entryRaw := rkentry.GlobalAppCtx.GetEntry(name)

	if entry, ok := entryRaw.(*rkprom.PromEntry); ok {
		return entry
	}

	return nil
}
