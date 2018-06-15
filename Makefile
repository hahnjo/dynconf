all: dep dynconf
.PHONY: all

dep:
	dep ensure
.PHONY: dep

dynconf:
	go build
.PHONY: dynconf

fmt:
	go fmt ./...
.PHONY: fmt

TEST_DIRS = ./pkg
COVERAGE ?= -cover
test:
	go test $(COVERAGE) $(TEST_DIRS)
.PHONY: test

bench:
	go test -bench . $(TEST_DIRS)
.PHONY: bench
