module github.com/rookie-ninja/rk-boot/database/sqlserver

go 1.16

require (
	github.com/rookie-ninja/rk-boot v1.4.1
	github.com/rookie-ninja/rk-db/sqlserver v0.0.3
	gorm.io/gorm v1.22.4
)

replace github.com/rookie-ninja/rk-boot => ../../