EXTRA_DEPFLAGS ?=
EXTRA_GOFLAGS ?=
EXTRA_TESTFLAGS ?=
COVERAGE_FLAGS ?= -cover

all: dep dynconf
.PHONY: all

dep:
	dep ensure $(EXTRA_DEPFLAGS)
.PHONY: dep

dynconf:
	go build $(EXTRA_GOFLAGS)
.PHONY: dynconf

fmt:
	go fmt ./...
.PHONY: fmt

TEST_DIRS = ./pkg
test:
	go test $(EXTRA_TESTFLAGS) $(COVERAGE) $(TEST_DIRS)
.PHONY: test

bench:
	go test $(EXTRA_TESTFLAGS) -bench . $(TEST_DIRS)
.PHONY: bench
