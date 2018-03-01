package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/Workiva/frugal/lib/gateway/gen-go/gateway_test"
	"github.com/Workiva/frugal/lib/go"
)

// Run an HTTP server
func main() {
	// Set the protocol used for serialization.
	// The protocol stack must match between client and server
	fProtocolFactory := frugal.NewFProtocolFactory(thrift.NewTBinaryProtocolFactoryDefault())

	// Create a handler. Each incoming request at the processor is sent to
	// the handler. Responses from the handler are returned back to the
	// client
	handler := &GatewayHandler{}
	processor := gateway_test.NewFGatewayTestProcessor(handler)

	// Start the server using the configured processor, and protocol
	mux := http.NewServeMux()
	mux.HandleFunc("/frugal", frugal.NewFrugalHandlerFunc(processor, fProtocolFactory))
	corsOptions := cors.Options{
		AllowedHeaders: []string{"Content-Transfer-Encoding"},
	}
	httpHandler := cors.New(corsOptions).Handler(mux)

	fmt.Println("Starting the http server...")
	log.Fatal(http.ListenAndServe(":9090", httpHandler))
}

// GatewayHandler handles all incoming requests to the server.
// The handler must satisfy the interface the server exposes.
type GatewayHandler struct{}

// CreateTweet creates a new tweet
func (f *GatewayHandler) GetContainer(ctx frugal.FContext, baseType *gateway_test.BaseType) (r *gateway_test.ContainerType, err error) {
	str_ := baseType.StringTest
	fmt.Println("Received", str_)
	bool_ := false

	if baseType.BoolTest != nil {
		bool_ = *baseType.BoolTest
	}

	listTest := []*gateway_test.BaseType{&gateway_test.BaseType{StringTest: str_, BoolTest: &bool_}}
	mapTest := map[string]*gateway_test.BaseType{"foo": &gateway_test.BaseType{StringTest: str_, BoolTest: &bool_}}
	enumTest := gateway_test.EnumType_ANOPTION

	container := &gateway_test.ContainerType{
		ListTest: listTest,
		MapTest:  mapTest,
		EnumTest: &enumTest,
	}

	fmt.Println("Container", container)
	return container, nil // TODO: test how error response is handled
}
