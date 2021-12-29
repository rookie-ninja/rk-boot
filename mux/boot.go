package rkbootmux

import (
	_ "github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-mux/boot"
)

func GetMuxEntry(name string) *rkmux.MuxEntry {
	if raw := rkentry.GlobalAppCtx.GetEntry(name); raw != nil {
		if res, ok := raw.(*rkmux.MuxEntry); ok {
			return res
		}
	}

	return nil
}
