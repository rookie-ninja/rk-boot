package rkbootzero

import (
	_ "github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-zero/boot"
)

func GetZeroEntry(name string) *rkzero.ZeroEntry {
	if raw := rkentry.GlobalAppCtx.GetEntry(name); raw != nil {
		if res, ok := raw.(*rkzero.ZeroEntry); ok {
			return res
		}
	}

	return nil
}
