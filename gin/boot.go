package rkbootgin

import (
	_ "github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-gin/boot"
)

func GetGinEntry(name string) *rkgin.GinEntry {
	if raw := rkentry.GlobalAppCtx.GetEntry(name); raw != nil {
		if res, ok := raw.(*rkgin.GinEntry); ok {
			return res
		}
	}

	return nil
}
