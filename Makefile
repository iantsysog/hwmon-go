GO ?= go
BINDIR ?= bin

TAGS ?= smc hid eventsystem iohid
CGO_ENABLED ?= 1

.PHONY: build test

HWMON_BIN := $(BINDIR)/hwmon-go

build: $(HWMON_BIN)

$(HWMON_BIN): | $(BINDIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -trimpath -tags="$(TAGS)" -ldflags="-s -w" -o $@ ./cmd/hwmon-go

$(BINDIR):
	mkdir -p $@

test:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test -tags="$(TAGS)" ./...