GO ?= go
BINDIR ?= bin
BIN ?= $(BINDIR)/smc-go

.PHONY: all build clean

all: build

build: $(BIN)

$(BIN): | $(BINDIR)
	$(GO) build -trimpath -ldflags="-s -w" -o $@ ./cmd/smc-go

$(BINDIR):
	mkdir -p $@

clean:
	rm -f $(BIN)
