module github.com/rookie-ninja/rk-boot/database/sqlite

go 1.16

require (
	github.com/rookie-ninja/rk-boot v1.4.1
	github.com/rookie-ninja/rk-db/sqlite v0.0.4
	gorm.io/gorm v1.22.4
)

replace github.com/rookie-ninja/rk-boot => ../../
