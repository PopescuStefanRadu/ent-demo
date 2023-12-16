.PHONY: test install-dependencies-locally test generate

install-dependencies-locally:
	go install entgo.io/ent/cmd/ent@v0.12.5
	go install go.uber.org/mock/mockgen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest

test:
	go test -count=1 ./...

generate:
	go generate ./...

lint:
	golangci-lint run

format:
	gofumpt -l -w .
	goimports -l -w .