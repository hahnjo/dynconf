dynconf:
	go build
.PHONY: dynconf

fmt:
	go fmt ./...
.PHONY: fmt

TEST_DIRS = ./pkg
test:
	go test -cover $(TEST_DIRS)
.PHONY: test

bench:
	go test -bench . $(TEST_DIRS)
.PHONY: bench
