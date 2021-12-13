package rkbootgrpc

import (
	_ "github.com/rookie-ninja/rk-boot"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-grpc/boot"
)

func GetGrpcEntry(name string) *rkgrpc.GrpcEntry {
	if raw := rkentry.GlobalAppCtx.GetEntry(name); raw != nil {
		if res, ok := raw.(*rkgrpc.GrpcEntry); ok {
			return res
		}
	}

	return nil
}
