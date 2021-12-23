package rkbootecho

import (
	_ "github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-echo/boot"
	"github.com/rookie-ninja/rk-entry/entry"
)

func GetEchoEntry(name string) *rkecho.EchoEntry {
	if raw := rkentry.GlobalAppCtx.GetEntry(name); raw != nil {
		if res, ok := raw.(*rkecho.EchoEntry); ok {
			return res
		}
	}

	return nil
}
