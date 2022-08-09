.PHONY: all
all: lint fmt

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


