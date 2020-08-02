LISTEAV ?= 0.0.0

GO     ?= go
GOFMT  ?= gofmt -s
GOPATH ?= $($(GO) env GOPATH)

RM		?= rm
CP		?= cp
SHA		?= sha256sum
CAT		?= cat
BIN		?= /usr/local/bin
MAKE    ?= make
SUDO    ?= sudo

ICON		:= icon/icon.png
ICONNEW		:= icon/icon_new.png
ICONERR		:= icon/icon_err.png
EXEC	    := listea-$(LISTEAV)
GOOUT		:= $(shell pwd)/bin
GOICON		:= icon/icon.go
GOICONNEW	:= icon/icon_new.go
GOICONERR	:= icon/icon_err.go
GOFILES		:= main.go $(GOICON) $(GOICONNEW) $(GOICONERR)

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

.PHONY: mod
mod:
	#$(GO) mod init git.iglou.eu/Laboratory/listea
	$(GO) mod tidy
	$(GO) mod verify

.PHONY: clean
clean:
	$(GO) clean
	$(RM) -f $(GOOUT)/*

.PHONY: icon
icon: 2goarray
	$(CAT) $(ICON) | $(GOPATH)/bin/2goarray Data icon > $(GOICON)
	$(CAT) $(ICONNEW) | $(GOPATH)/bin/2goarray New icon > $(GOICONNEW)
	$(CAT) $(ICONERR) | $(GOPATH)/bin/2goarray Err icon > $(GOICONERR)
	$(MAKE) fmt

.PHONY: fmt
fmt:
	$(GOFMT) -e -w $(GOFILES)

.PHONY: 2goarray
2goarray: 2goarray
	$(GO) get github.com/cratonica/2goarray

.PHONY: releases
releases: mob
	@env GOOS=linux GOARCH=amd64 $(GO) build -o $(GOOUT)/$(EXEC)-linux-amd64 -a
	@cd $(GOOUT) && $(SHA) $(EXEC)-linux-amd64 > $(EXEC)-linux-amd64.sha256

	#@env GOOS=darwin GOARCH=amd64 $(GO) build -o $(GOOUT)/$(EXEC)-darwin-amd64 -a

	@env GO111MODULE=on GOOS=windows GOARCH=amd64 $(GO) build -ldflags "-H=windowsgui" -o $(GOOUT)/$(EXEC)-windows-amd64.exe -a
	@cd $(GOOUT) && $(SHA) $(EXEC)-windows-amd64.exe > $(EXEC)-windows-amd64.exe.sha256

.PHONY: help
help:
	@echo "Make Routines:"
	@echo " - \"\"                get go deps and build listea"
	@echo " - install           install listea on system"
	@echo " - remove            remove binary from system"
	@echo " - upgrade           build and update installed binary"
	@echo " - fmt               automatically re-formats go code"