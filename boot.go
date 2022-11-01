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
	"github.com/rookie-ninja/rk-entry/v2/entry"
	"github.com/rookie-ninja/rk-entry/v2/middleware"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"runtime/debug"
)

type hookFuncM map[string]map[string]func(ctx context.Context)

func newHookFuncM() hookFuncM {
	return map[string]map[string]func(ctx context.Context){}
}

func (m hookFuncM) addFunc(entryType, entryName string, f func(ctx context.Context)) {
	if _, ok := m[entryType]; !ok {
		m[entryType] = make(map[string]func(ctx context.Context))
	}

	m[entryType][entryName] = f
}

func (m hookFuncM) getFunc(entryType, entryName string) func(ctx context.Context) {
	if inner, ok := m[entryType]; ok {
		if f, ok := inner[entryName]; ok {
			return f
		}
	}

	return func(ctx context.Context) {}
}

// Boot is a structure for bootstrapping rk style application
type Boot struct {
	bootConfigPath string    `yaml:"-" json:"-"`
	embedFS        *embed.FS `yaml:"-" json:"-"`
	bootConfigRaw  []byte    `yaml:"-" json:"-"`
	beforeHookF    hookFuncM `yaml:"-" json:"-"`
	afterHookF     hookFuncM `yaml:"-" json:"-"`
	EventId        string    `yaml:"-" json:"-"`
	pluginEntries  map[string]map[string]rkentry.Entry
	userEntries    map[string]map[string]rkentry.Entry
	webEntries     map[string]map[string]rkentry.Entry
}

// BootOption is used as options while bootstrapping from code
type BootOption func(*Boot)

// WithBootConfigPath provide boot config yaml file.
func WithBootConfigPath(filePath string, fs *embed.FS) BootOption {
	return func(boot *Boot) {
		boot.bootConfigPath = filePath
		boot.embedFS = fs
	}
}

// WithBootConfigRaw provide boot config as string.
func WithBootConfigRaw(raw []byte) BootOption {
	return func(boot *Boot) {
		if len(raw) > 0 {
			boot.bootConfigRaw = raw
		}
	}
}

// NewBoot create a bootstrapper.
func NewBoot(opts ...BootOption) *Boot {
	defer syncLog("N/A")

	boot := &Boot{
		EventId:       rkmid.GenerateRequestId(),
		beforeHookF:   newHookFuncM(),
		afterHookF:    newHookFuncM(),
		pluginEntries: map[string]map[string]rkentry.Entry{},
		userEntries:   map[string]map[string]rkentry.Entry{},
		webEntries:    map[string]map[string]rkentry.Entry{},
	}

	for i := range opts {
		opts[i](boot)
	}

	raw := boot.readYAML()

	// Register entries need to pre-build.
	rkentry.BootstrapBuiltInEntryFromYAML(raw)

	for _, f := range rkentry.ListPluginEntryRegFunc() {
		for _, v := range f(raw) {
			if boot.pluginEntries[v.GetType()] == nil {
				boot.pluginEntries[v.GetType()] = make(map[string]rkentry.Entry)
			}
			boot.pluginEntries[v.GetType()][v.GetName()] = v
		}
	}

	for _, f := range rkentry.ListUserEntryRegFunc() {
		for _, v := range f(raw) {
			if boot.userEntries[v.GetType()] == nil {
				boot.userEntries[v.GetType()] = make(map[string]rkentry.Entry)
			}
			boot.userEntries[v.GetType()][v.GetName()] = v
		}
	}

	for _, f := range rkentry.ListWebFrameEntryRegFunc() {
		for _, v := range f(raw) {
			if boot.webEntries[v.GetType()] == nil {
				boot.webEntries[v.GetType()] = make(map[string]rkentry.Entry)
			}
			boot.webEntries[v.GetType()][v.GetName()] = v
		}
	}

	return boot
}

// AddHookFuncBeforeBootstrap run functions before certain entry Bootstrap()
func (boot *Boot) AddHookFuncBeforeBootstrap(entryType, entryName string, f func(ctx context.Context)) {
	if f == nil {
		return
	}

	boot.beforeHookF.addFunc(entryType, entryName, f)
}

// AddHookFuncAfterBootstrap run functions before certain entry Bootstrap()
func (boot *Boot) AddHookFuncAfterBootstrap(entryType, entryName string, f func(ctx context.Context)) {
	if f == nil {
		return
	}

	boot.afterHookF.addFunc(entryType, entryName, f)
}

// Bootstrap entries as sequence of plugin, user defined and web framework
func (boot *Boot) Bootstrap(ctx context.Context) {
	defer syncLog(boot.EventId)

	ctx = context.WithValue(ctx, "eventId", boot.EventId)

	for entryType, byEntryName := range boot.pluginEntries {
		for entryName, e := range byEntryName {
			boot.beforeHookF.getFunc(entryType, entryName)(ctx)
			e.Bootstrap(ctx)
			boot.afterHookF.getFunc(entryType, entryName)(ctx)
		}
	}

	for entryType, byEntryName := range boot.userEntries {
		for entryName, e := range byEntryName {
			boot.beforeHookF.getFunc(entryType, entryName)(ctx)
			e.Bootstrap(ctx)
			boot.afterHookF.getFunc(entryType, entryName)(ctx)
		}
	}

	for entryType, byEntryName := range boot.webEntries {
		for entryName, e := range byEntryName {
			boot.beforeHookF.getFunc(entryType, entryName)(ctx)
			e.Bootstrap(ctx)
			boot.afterHookF.getFunc(entryType, entryName)(ctx)
		}
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

// Interrupt entries as sequence of plugin, user defined and web framework
func (boot *Boot) interrupt(ctx context.Context) {
	defer syncLog(boot.EventId)

	ctx = context.WithValue(ctx, "eventId", boot.EventId)

	// Interrupt external entries
	for _, m := range rkentry.GlobalAppCtx.ListEntries() {
		for _, e := range m {
			e.Interrupt(ctx)
		}
	}
}

// readYAML read YAML file
func (boot *Boot) readYAML() []byte {
	// case 1: if user provide raw then, continue
	if len(boot.bootConfigRaw) > 0 {
		return boot.bootConfigRaw
	}

	// case 2: if embed.FS is not nil, then try to read from it
	if boot.embedFS != nil {
		res, err := boot.embedFS.ReadFile(boot.bootConfigPath)
		if err != nil {
			rkentry.ShutdownWithError(err)
		}
		return res
	}

	// case 3: try to read from local, if bootConfigPath is empty, then try to read from default boot.yaml
	if len(boot.bootConfigPath) < 1 {
		boot.bootConfigPath = "boot.yaml"
	}
	if !filepath.IsAbs(boot.bootConfigPath) {
		wd, _ := os.Getwd()
		boot.bootConfigPath = filepath.Join(wd, boot.bootConfigPath)
	}

	res, err := os.ReadFile(boot.bootConfigPath)
	if err != nil {
		rkentry.ShutdownWithError(err)
	}
	return res
}

// sync logs
func syncLog(eventId string) {
	if r := recover(); r != nil {
		stackTrace := fmt.Sprintf("Panic occured, shutting down... \n%s", string(debug.Stack()))
		for _, v := range rkentry.GlobalAppCtx.ListEntriesByType(rkentry.LoggerEntryType) {
			logger, ok := v.(*rkentry.LoggerEntry)
			if !ok {
				continue
			}

			if logger != nil {
				logger.Error(stackTrace,
					zap.String("eventId", eventId),
					zap.Any("RootCause", r))
			}
			logger.Sync()
		}

		for _, v := range rkentry.GlobalAppCtx.ListEntriesByType(rkentry.EventEntryType) {
			event, ok := v.(*rkentry.EventEntry)
			if !ok {
				continue
			}

			event.Sync()
		}

		panic(r)
	}
}
