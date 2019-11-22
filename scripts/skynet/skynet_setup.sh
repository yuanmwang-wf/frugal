#!/usr/bin/env bash

set -exo pipefail

mkdir -p /go/src/github.com/Workiva/

# Symlink frugal to gopath - this allows skynet-cli editing for interactive/directmount
ln -s /testing/ /go/src/github.com/Workiva/frugal

# Install frugal
cd $GOPATH/src/github.com/Workiva/frugal && go install

# Start nats-server
nats-server &

# Start activemq broker
cd /opt/activemq/bin
./activemq start

cd $GOPATH/src/github.com/Workiva/frugal
