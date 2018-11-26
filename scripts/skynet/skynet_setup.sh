#!/usr/bin/env bash

set -exo pipefail

which glide > /dev/null || {
    curl https://glide.sh/get | sh
}

mkdir -p /go/src/github.com/Workiva/

# Symlink frugal to gopath - this allows skynet-cli editing for interactive/directmount
ln -s /testing/ /go/src/github.com/Workiva/frugal

# Install frugal
cd $GOPATH/src/github.com/Workiva/frugal && go install

# Start gnatsd
gnatsd &

# TODO this install should be in messaging-docker-images
wget http://archive.apache.org/dist/activemq/5.15.6/apache-activemq-5.15.6-bin.tar.gz
tar -xzf apache-activemq-5.15.2-bin.tar.gz
cd apache-activemq-5.15.2/bin
./activemq start
cd $GOPATH/src/github.com/Workiva/frugal
