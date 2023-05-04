.PHONY: all
all: lint fmt swag

.PHONY: lint
lint:
	@echo "[golangci-lint] Running golangci-lint..."
	@golangci-lint run 2>&1
	@echo "------------------------------------[Done]"

.PHONY: fmt
fmt:
	@echo "[fmt] Format go project..."
	@gofmt -s -w . 2>&1
	@echo "------------------------------------[Done]"

.PHONY: swag
swag:
	@echo "[swag] Running swag..."
	@swag init --generalInfo main.go --propertyStrategy camelcase
	@rm -rf docs/docs.go
	@echo "------------------------------------[Done]"

.PHONY: zip
zip:
	@echo "[zip] Compress to zip file..."
	@zip -r rk-demo.zip *
	@echo "------------------------------------[Done]"

