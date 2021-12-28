module github.com/rookie-ninja/rk-demo

go 1.16

require (
	github.com/gofiber/fiber/v2 v2.23.0
	github.com/rookie-ninja/rk-boot v1.4.0
	github.com/rookie-ninja/rk-boot/fiber v0.0.0-00010101000000-000000000000
)

replace github.com/rookie-ninja/rk-boot => ../../

replace github.com/rookie-ninja/rk-boot/fiber => ../../fiber
