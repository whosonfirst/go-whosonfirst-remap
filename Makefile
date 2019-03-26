CWD=$(shell pwd)
GOPATH := $(CWD)

build:	rmdeps deps fmt bin

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-remap
	cp -r *.go src/github.com/whosonfirst/go-whosonfirst-remap/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-index"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-uri"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-export"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/sjson"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	cp -r src/ vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt *.go

bin:	self
	@GOPATH=$(GOPATH) go build -o bin/wof-remap cmd/wof-remap.go
	@GOPATH=$(GOPATH) go build -o bin/wof-contract cmd/wof-contract.go
	@GOPATH=$(GOPATH) go build -o bin/wof-expand cmd/wof-expand.go
