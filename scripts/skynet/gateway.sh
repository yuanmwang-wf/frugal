#!/usr/bin/env bash


set -ex

./scripts/skynet/skynet_setup.sh

export FRUGAL_HOME=$GOPATH/src/github.com/Workiva/frugal

if [ ! -e "$FRUGAL_HOME/lib/go/glide.lock" ]; then
    cd $FRUGAL_HOME/lib/go && glide install
fi

cd $FRUGAL_HOME

# setup a vendor folder with frugals dependencies and frugal
cd test/integration/gateway
mkdir vendor
ls $FRUGAL_HOME/lib/go/vendor
cp -r $FRUGAL_HOME/lib/go/vendor/* vendor/
rm -rf $FRUGAL_HOME/lib/go/vendor
mkdir -p vendor/github.com/Workiva/frugal/lib/go
cp -r $FRUGAL_HOME/lib/go/* vendor/github.com/Workiva/frugal/lib/go


frugal --gen go gateway_test.frugal
frugal --gen gateway gateway_test.frugal

rm -rf $FRUGAL_HOME/vendor

# TODO janky stuff because vendor isn't being found
# And not all are in the copied vendor folder
mkdir $GOPATH/src/git.apache.org/
cp -r vendor/git.apache.org/* $GOPATH/src/git.apache.org/thrift.git
go get github.com/Sirupsen/logrus
go get github.com/gorilla/mux
go get github.com/mattrobenolt/gocql/uuid
go get github.com/nats-io/go-nats
go get github.com/rs/cors
rm -rf vendor

go run http_server/main.go &
go run gateway_server/main.go &
sleep 2s # make sure everything is running

go test
