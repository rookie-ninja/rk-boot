module github.com/rookie-ninja/rk-example

go 1.15

require (
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/rookie-ninja/rk-boot v1.1.2
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)

replace github.com/rookie-ninja/rk-boot => ../../
