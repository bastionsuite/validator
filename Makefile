build: vet test
	go build ./...

test:
	go test -v -cover ./...

sast:
	gosec -quiet ./...

sca:
	snyk test

vet:
	go vet ./...
	staticcheck -checks all ./...

all: test

.PHONY: test all
