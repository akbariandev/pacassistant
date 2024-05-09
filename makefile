generate-mock:
	mockery --dir=internal/domain --output=internal/domain/mocks --outpkg=mocks --all

build-bot:
	go build -mod vendor -ldflags "-X 'github.com/akbariandev/pacassistant/cmd/bot/commands.CommitID=$(shell git rev-parse --short HEAD 2>/dev/null)' -X 'github.com/akbariandev/pacassistant/cmd/bot/commands.BuildDate=$(shell date)'" -o main cmd/bot/main.go

check:
	gofumpt -l -w cmd
	gofumpt -l -w config
	gofumpt -l -w internal
	gofumpt -l -w pkg
	gofumpt -l -w migration
	gofumpt -l -w transport
	govulncheck ./...
	golangci-lint run --timeout=20m0s

unit_test:
	go test ./...

race_test:
	go test ./... --race

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/vektra/mockery/v2@v2.36.0

submodule:
	git submodule update --init --recursive

