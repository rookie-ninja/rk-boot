module github.com/rookie-ninja/rk-demo

go 1.16

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.1
	github.com/rookie-ninja/rk-boot v1.3.6
	github.com/rookie-ninja/rk-grpc v1.2.14
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/rookie-ninja/rk-boot => ../../
