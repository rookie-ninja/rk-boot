.PHONY: all
all: lint readme fmt

.PHONY: lint
lint:
	@echo "running golangci-lint..."
	@golangci-lint run 2>&1

.PHONY: fmt
fmt:
	@echo "format go project..."
	@gofmt -s -w . 2>&1

.PHONY: swag
swag:
	@echo "running swag..."
	@swag init --generalInfo gin.go --propertyStrategy camelcase

.PHONY: readme
readme:
	@echo "running doctoc..."
	@doctoc . 2>&1

