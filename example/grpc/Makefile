.PHONY: all
all: readme fmt buf

.PHONY: fmt
fmt:
	@echo "format go project..."
	@gofmt -s -w . 2>&1

.PHONY: readme
readme:
	@echo "running doctoc..."
	@doctoc . 2>&1

.PHONY: buf
buf:
	@echo "running buf..."
	@buf generate --path api/v1
