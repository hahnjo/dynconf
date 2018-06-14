dynconf:
	go build
.PHONY: dynconf

fmt:
	go fmt ./...
.PHONY: fmt

test:
	go test ./pkg
.PHONY: test
