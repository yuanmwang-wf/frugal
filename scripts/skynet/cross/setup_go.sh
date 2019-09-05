#!/usr/bin/env bash

set -ex

export FRUGAL_HOME=$GOPATH/src/github.com/Workiva/frugal

cd $FRUGAL_HOME/lib/go && GO111MODULE=on go mod vendor
cd $FRUGAL_HOME

# setup a vendor folder with frugals dependencies and frugal
cd test/integration/go
mkdir vendor
cp -r $FRUGAL_HOME/lib/go/vendor/* vendor/
rm -rf $FRUGAL_HOME/lib/go/vendor
mkdir -p vendor/github.com/Workiva/frugal/lib/go
cp -r $FRUGAL_HOME/lib/go/* vendor/github.com/Workiva/frugal/lib/go

# Create Go binaries
rm -rf test/integration/go/bin/*
go build -o bin/testclient src/bin/testclient/main.go
go build -o bin/testserver src/bin/testserver/main.go
go build -o bin/testpublisher src/bin/testpublisher/main.go
go build -o bin/testsubscriber src/bin/testsubscriber/main.go
