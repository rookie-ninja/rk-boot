.PHONY: all
all: gomod doctoc fmt

.PHONY: fmt
fmt:
	@echo "[fmt] Formatting go project..."
	@gofmt -s -w . 2>&1
	@echo "------------------------------------[Done]"

.PHONY: doctoc
doctoc:
	@echo "[doctoc] Running doctoc..."
	@doctoc . 2>&1
	@echo "------------------------------------[Done]"

.PHONY: gomod
gomod:
	@echo "[gomod] Running go mod tidy..."
	@go mod tidy
	@echo "------------------------------------[Done]"
