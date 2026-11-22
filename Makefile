SERVICE ?= gateway
export GOWORK := off

.PHONY: test lint marketplace

test:
	@cd "$(SERVICE)" && go test ./...
	@printf '✅ Tests passed\n'

lint: bin/golangci-lint
	@cd "$(SERVICE)" && "$(CURDIR)/bin/golangci-lint" run --config "$(CURDIR)/.golangci.yml"
	@printf '✅ Lint passed\n'

bin/golangci-lint:
	@GOBIN="$(CURDIR)/bin" go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0

marketplace:
	@goreleaser build --snapshot --clean --config collector/.goreleaser.yaml
