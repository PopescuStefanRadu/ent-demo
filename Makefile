.PHONY: test install-dependencies-locally

install-dependencies-locally:
	go install entgo.io/ent/cmd/ent@v0.12.5

test:
	go test -count=1 ./...

generate:
	go generate ./...