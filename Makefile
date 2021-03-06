.PHONY: all clean help

SRC := $(filter-out $(wildcard examples/*/*_test.go), $(wildcard examples/*/*.go))
BIN := $(subst .go,,$(SRC))

all: $(BIN) linter	# build all binary

clean:		# clean-up the environment
	rm -f $(BIN)

help:		# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

GO      := go
GOFMT   := $(GO)fmt -w -s
GOFLAG  := -ldflags="-s -w"
GOTEST  := $(GO) test -cover -failfast -timeout 2s
GOBENCH := $(GO) test -bench=. -cover -failfast -benchmem

linter:
	$(GOFMT) $(shell find . -name '*.go')
	$(GOTEST) ./...
	$(GOBENCH)

$(BIN): linter

%: %.go
	$(GO) build $(GOFLAG) -o $@ $<
