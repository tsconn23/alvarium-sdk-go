.PHONY: test

test:
	go test ./... -v -coverprofile=coverage.out ./...
	go vet ./...
	gofmt -l .
	[ "`gofmt -l .`" = "" ]