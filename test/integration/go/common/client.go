/*
 * Copyright 2017 Workiva
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"flag"
	"fmt"

	"net/http"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/Workiva/frugal/lib/go"
	"github.com/Workiva/frugal/test/integration/go/gen/frugaltest"
)

var debugClientProtocol bool

func init() {
	flag.BoolVar(&debugClientProtocol, "debug_client_protocol", false, "turn client protocol trace on")
}

func StartClient(
	host string,
	port int64,
	transport string,
	protocol string,
	pubSub chan bool,
	sent chan bool,
	clientMiddlewareCalled chan bool) (client *frugaltest.FFrugalTestClient, err error) {

	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	default:
		return nil, fmt.Errorf("Invalid protocol specified %s", protocol)
	}

	fProtocolFactory := frugal.NewFProtocolFactory(protocolFactory)
	natsConn := getNatsConn()


	// RPC client
	var trans frugal.FTransport
	switch transport {
	case NatsName:
		trans = frugal.NewFNatsTransport(natsConn, fmt.Sprintf("frugal.foo.bar.rpc.%d", port), "")
	case HttpName:
		// Set request and response capacity to 1mb
		maxSize := uint(1048576)
		trans = frugal.NewFHTTPTransportBuilder(&http.Client{}, fmt.Sprintf("http://localhost:%d", port)).WithRequestSizeLimit(maxSize).WithResponseSizeLimit(maxSize).Build()
	default:
		return nil, fmt.Errorf("Invalid transport specified %s", transport)
	}

	if err := trans.Open(); err != nil {
		return nil, fmt.Errorf("Error opening transport %s", err)
	}

	client = frugaltest.NewFFrugalTestClient(frugal.NewFServiceProvider(trans, fProtocolFactory), clientLoggingMiddleware(clientMiddlewareCalled))
	return
}
