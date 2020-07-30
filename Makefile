GO     ?= go
GOFMT  ?= gofmt -s
GOPATH ?= $($(GO) env GOPATH)

RM		?= rm
CP		?= cp
CAT		?= cat
BIN		?= /usr/local/bin
MAKE    ?= make
SUDO    ?= sudo

ICON	:= icon/icon.png
EXEC    := listea
GOOUT   := $(shell pwd)/bin
GOICON	:= icon/icon.go
GOFILES := main.go icon/icon.go

.PHONY: build
build: godep
	$(GO) build -o $(GOOUT)/$(EXEC) -a

.PHONY: install
install: build
	$(SUDO) $(CP) $(GOOUT)/$(EXEC) $(BIN)/$(EXEC)
	$(MAKE) clean

.PHONY: remove
remove:
	$(SUDO) $(RM) -f $(BIN)/$(EXEC)

.PHONY: upgrade
upgrade: remove install

.PHONY: godep
godep:
	$(GO) get github.com/getlantern/systray
	$(GO) get github.com/skratchdot/open-golang/open

.PHONY: clean
clean:
	$(GO) clean
	$(RM) -f $(GOOUT)/*

.PHONY: icon
icon: 2goarray
	$(CAT) $(ICON) | $(GOPATH)/bin/2goarray Data icon > $(GOICON)
	$(MAKE) fmt

.PHONY: fmt
fmt:
	$(GOFMT) -e -w $(GOFILES)

.PHONY: 2goarray
2goarray: 2goarray
	$(GOFMT) get github.com/cratonica/2goarray

.PHONY: help
help:
	@echo "Make Routines:"
	@echo " - \"\"                get go deps and build listea"
	@echo " - install           install listea on system"
	@echo " - remove            remove binary from system"
	@echo " - upgrade           build and update installed binary"
	@echo " - fmt               automatically re-formats go code"