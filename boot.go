// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

// Package rkboot is bootstrapper for rk style application
package rkboot

import (
	"context"
	"github.com/rookie-ninja/rk-echo/boot"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-gin/boot"
	"github.com/rookie-ninja/rk-grpc/boot"
	"github.com/rookie-ninja/rk-prom"
)

// Boot is a structure for bootstrapping rk style application
type Boot struct {
	BootConfigPath string `yaml:"bootConfigPath" json:"bootConfigPath"`
}

// BootOption is used as options while bootstrapping from code
type BootOption func(*Boot)

// WithBootConfigPath provide boot config yaml file.
func WithBootConfigPath(filePath string) BootOption {
	return func(boot *Boot) {
		boot.BootConfigPath = filePath
	}
}

// NewBoot create a bootstrapper.
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

// GetGinEntry returns rkgin.GinEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetGinEntry(name string) *rkgin.GinEntry {
	entryRaw := rkentry.GlobalAppCtx.GetEntry(name)

	if entry, ok := entryRaw.(*rkgin.GinEntry); ok {
		return entry
	}

	return nil
}

// GetGrpcEntry returns rkgrpc.GrpcEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetGrpcEntry(name string) *rkgrpc.GrpcEntry {
	entryRaw := rkentry.GlobalAppCtx.GetEntry(name)

	if entry, ok := entryRaw.(*rkgrpc.GrpcEntry); ok {
		return entry
	}

	return nil
}

// GetEchoEntry returns rkecho.EchoEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetEchoEntry(name string) *rkecho.EchoEntry {
	entryRaw := rkentry.GlobalAppCtx.GetEntry(name)

	if entry, ok := entryRaw.(*rkecho.EchoEntry); ok {
		return entry
	}

	return nil
}

// GetPromEntry returns rkprom.PromEntry from rkentry.GlobalAppCtx.
func (boot *Boot) GetPromEntry(name string) *rkprom.PromEntry {
	entryRaw := rkentry.GlobalAppCtx.GetEntry(name)

	if entry, ok := entryRaw.(*rkprom.PromEntry); ok {
		return entry
	}

	return nil
}
