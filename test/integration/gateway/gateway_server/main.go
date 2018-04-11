package main

import (
	"fmt"
	"net/http"
	"reflect"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/Workiva/frugal/lib/gateway"
	gateway_gen "github.com/Workiva/frugal/test/integration/gateway/gen-go/gateway_test"
	"github.com/Workiva/frugal/lib/go"
)

var mockEndpoint = "http://localhost:9090/frugal"

func newLoggingMiddleware() frugal.ServiceMiddleware {
	return func(next frugal.InvocationHandler) frugal.InvocationHandler {
		return func(service reflect.Value, method reflect.Method, args frugal.Arguments) frugal.Results {
			fmt.Printf("==== CALLING %s.%s ====\n", service.Type(), method.Name)
			ret := next(service, method, args)
			fmt.Printf("==== CALLED  %s.%s ====\n", service.Type(), method.Name)
			return ret
		}
	}
}

func newHeaderLoggingMiddleware() gateway.GatewayMiddleware {
	return func(next gateway.InvocationHandler) gateway.InvocationHandler {
		return func(service reflect.Value, method reflect.Method, headers http.Header, args gateway.Arguments) gateway.Results {
			fmt.Printf("headers: %+v\n", headers)
			return next(service, method, headers, args)
		}
	}
}

// Create a new Frugal client connected to the backing service
func newClient() *gateway_gen.FGatewayTestClient {
	// Set the protocol used for serialization.
	// The protocol stack must match between client and server
	fProtocolFactory := frugal.NewFProtocolFactory(thrift.NewTBinaryProtocolFactoryDefault())

	// Create an HTTP transport listening
	httpTransport := frugal.NewFHTTPTransportBuilder(&http.Client{}, mockEndpoint).Build()
	defer httpTransport.Close()
	if err := httpTransport.Open(); err != nil {
		panic(err)
	}

	// Create a provider with the transport and protocol factory. The provider
	// can be used to create multiple Clients.
	provider := frugal.NewFServiceProvider(httpTransport, fProtocolFactory)

	// Create a client used to send messages with our desired protocol.  You
	// can also pass middleware in here if you only want it to intercept calls
	// for this specific client.
	storeClient := gateway_gen.NewFGatewayTestClient(provider, newLoggingMiddleware())

	return storeClient
}

func main() {
	context := gateway_gen.NewGatewayTestContext(newClient(), gateway.NewMarshalerRegistry(), newHeaderLoggingMiddleware())
	//context := gateway_gen.GatewayTestContext{
	//	Marshalers: gateway.NewMarshalerRegistry(),
	//	Client:     newClient(),
	//}

	router, _ := gateway_gen.MakeRouter(context)

	// 	// TODO: Compile function MakeRouter(context) to return an HTTP mux router
	// 	handler := &gateway_test.GatewayTestHandler{&context, gateway_test.GatewayTestGetContainerHandler}

	// 	router := mux.NewRouter()

	// 	router.Methods("POST").Path("/v1/{differentString}/").Name("GatewayTestGetContainerHandler").Handler(handler)

	http.ListenAndServe(":5000", router)
}
