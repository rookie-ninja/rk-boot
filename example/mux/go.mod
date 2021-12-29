module github.com/rookie-ninja/rk-demo

go 1.16

require (
	github.com/rookie-ninja/rk-boot v1.4.0
	github.com/rookie-ninja/rk-boot/mux v0.0.0-20211229200216-30d9fd733c94
)

replace github.com/rookie-ninja/rk-boot => ../../

replace github.com/rookie-ninja/rk-boot/mux => ../../mux
