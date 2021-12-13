package rkbootgin

import (
	_ "github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-gf/boot"
)

func GetGfEntry(name string) *rkgf.GfEntry {
	if raw := rkentry.GlobalAppCtx.GetEntry(name); raw != nil {
		if res, ok := raw.(*rkgf.GfEntry); ok {
			return res
		}
	}

	return nil
}
