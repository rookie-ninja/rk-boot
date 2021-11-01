.PHONY: all
all: lint doctoc fmt buf

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

.PHONY: doctoc
doctoc:
	@echo "[doctoc] Running doctoc..."
	@doctoc . 2>&1
	@echo "------------------------------------[Done]"

.PHONY: buf
buf:
	@echo "[buf] Running buf..."
	@buf generate --path api/v1
