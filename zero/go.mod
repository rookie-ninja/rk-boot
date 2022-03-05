module github.com/rookie-ninja/rk-boot/zero

go 1.16

require (
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/openzipkin/zipkin-go v0.4.0 // indirect
	github.com/rookie-ninja/rk-boot v1.4.0
	github.com/rookie-ninja/rk-entry v1.0.11
	github.com/rookie-ninja/rk-zero v0.0.10
	go.opentelemetry.io/otel/exporters/zipkin v1.3.0 // indirect
)

replace github.com/rookie-ninja/rk-boot => ../
