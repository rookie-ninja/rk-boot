module github.com/rookie-ninja/rk-demo

go 1.16

require (
	github.com/rookie-ninja/rk-boot v1.4.0
	github.com/rookie-ninja/rk-boot/zero v0.0.1
	github.com/tal-tech/go-zero v1.2.4
)

replace github.com/rookie-ninja/rk-boot => ../../

replace github.com/rookie-ninja/rk-boot/zero => ../../zero
