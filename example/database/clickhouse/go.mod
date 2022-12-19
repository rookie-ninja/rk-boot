module github.com/rookie-ninja/rk-demo

go 1.18

require (
	github.com/gin-gonic/gin v1.8.1
	github.com/rookie-ninja/rk-boot/v2 v2.2.6
	github.com/rookie-ninja/rk-db/clickhouse v1.2.14
	github.com/rookie-ninja/rk-gin/v2 v2.2.16
	github.com/rs/xid v1.3.0
	gorm.io/gorm v1.24.0
)

replace github.com/rookie-ninja/rk-boot/v2 => ../../../

require (
	github.com/ClickHouse/ch-go v0.48.0 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.3.0 // indirect
	github.com/andres-erbsen/clock v0.0.0-20160526145045-9e14626cd129 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/gin-contrib/pprof v1.4.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.6.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.1 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/openzipkin/zipkin-go v0.4.0 // indirect
	github.com/paulmach/orb v0.7.1 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.13.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/rookie-ninja/rk-entry/v2 v2.2.14 // indirect
	github.com/rookie-ninja/rk-logger v1.2.13 // indirect
	github.com/rookie-ninja/rk-query v1.2.14 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.12.0 // indirect
	github.com/subosito/gotenv v1.3.0 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	go.opentelemetry.io/contrib v1.8.0 // indirect
	go.opentelemetry.io/otel v1.10.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.8.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.10.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.10.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.10.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.8.0 // indirect
	go.opentelemetry.io/otel/exporters/zipkin v1.10.0 // indirect
	go.opentelemetry.io/otel/sdk v1.10.0 // indirect
	go.opentelemetry.io/otel/trace v1.10.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/ratelimit v0.2.0 // indirect
	go.uber.org/zap v1.23.0 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20220920203100-d0c6ba3f52d9 // indirect
	golang.org/x/sys v0.0.0-20221006211917-84dc82d7e875 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220920201722-2b89144ce006 // indirect
	google.golang.org/grpc v1.49.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/ini.v1 v1.66.6 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/clickhouse v0.5.0 // indirect
)
