SHELL=PATH='$(PATH)' /bin/sh

GOBUILD=CGO_ENABLED=0 go build -ldflags '-w -s'

PLATFORM := $(shell uname -o)

NAME := BAS.exe
OS := windows

ifeq ($(PLATFORM), Msys)
    INCLUDE := ${shell echo "$(GOPATH)"|sed -e 's/\\/\//g'}
else ifeq ($(PLATFORM), Cygwin)
    INCLUDE := ${shell echo "$(GOPATH)"|sed -e 's/\\/\//g'}
else
	INCLUDE := $(GOPATH)
	NAME=BAS
	OS=linux
endif

# enable second expansion
.SECONDEXPANSION:

.PHONY: all

BINDIR=$(INCLUDE)/bin

all: build

build:
	GOOS=$(OS) GOARCH=amd64 $(GOBUILD) -o $(BINDIR)/$(NAME)

mac:
	GOOS=darwin go build -ldflags '-w -s' -o $(BINDIR)/$(NAME).mac
arm:
	GOOS=linux GOARM=7 GOARCH=arm go build -ldflags '-w -s' -o $(BINDIR)/$(NAME).arm
lnx:
	GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o $(BINDIR)/$(NAME).lnx

ios:
	gomobile bind -v -o $(BINDIR)/bas.framework -target=ios github.com/hyperorchidlab/BAS/ios

clean:
	rm $(BINDIR)/$(NAME)
