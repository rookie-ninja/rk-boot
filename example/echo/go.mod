module github.com/rookie-ninja/rk-demo

go 1.16

require (
	github.com/labstack/echo/v4 v4.6.1
	github.com/rookie-ninja/rk-boot v1.4.0
	github.com/rookie-ninja/rk-boot/echo v0.0.5
)

replace github.com/rookie-ninja/rk-boot => ../../

replace github.com/rookie-ninja/rk-boot/echo => ../../echo
