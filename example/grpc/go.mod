module github.com/rookie-ninja/rk-demo

go 1.16

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.1
	github.com/rookie-ninja/rk-boot v1.4.0
	github.com/rookie-ninja/rk-boot/grpc v1.2.15
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/rookie-ninja/rk-boot => ../../

replace github.com/rookie-ninja/rk-boot/grpc => ../../grpc
