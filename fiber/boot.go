package rkbootfiber

import (
	_ "github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-fiber/boot"
)

func GetFiberEntry(name string) *rkfiber.FiberEntry {
	if raw := rkentry.GlobalAppCtx.GetEntry(name); raw != nil {
		if res, ok := raw.(*rkfiber.FiberEntry); ok {
			return res
		}
	}

	return nil
}
