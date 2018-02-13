package main

import (
	"fmt"
	"log"
	"net/http"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/rs/cors"

	"github.com/Workiva/frugal/lib/go"
	"github.com/Workiva/frugal/test/gateway/gen-go/type_test"
)

// Run an HTTP server
func main() {
	// Set the protocol used for serialization.
	// The protocol stack must match between client and server
	fProtocolFactory := frugal.NewFProtocolFactory(thrift.NewTBinaryProtocolFactoryDefault())

	// Create a handler. Each incoming request at the processor is sent to
	// the handler. Responses from the handler are returned back to the
	// client
	handler := &Handler{}
	processor := type_test.NewFTypeTestProcessor(handler)

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

// The handler must satisfy the interface the server exposes.
type Handler struct{}

// TODO: Test that each base type can be mapped correctly
func (f *Handler) GetBoolArgument(ctx frugal.FContext, request *type_test.BaseType) (r *type_test.BaseType, err error) {
	return nil, nil
}

func (f *Handler) GetByteArgument(ctx frugal.FContext, request *type_test.BaseType) (r *type_test.BaseType, err error) {
	return nil, nil
}

func (f *Handler) GetI16Argument(ctx frugal.FContext, request *type_test.BaseType) (r *type_test.BaseType, err error) {
	return nil, nil
}

func (f *Handler) GetI32Argument(ctx frugal.FContext, request *type_test.BaseType) (r *type_test.BaseType, err error) {
	return nil, nil
}

func (f *Handler) GetI64Argument(ctx frugal.FContext, request *type_test.BaseType) (r *type_test.BaseType, err error) {
	return nil, nil
}

func (f *Handler) GetDoubleArgument(ctx frugal.FContext, request *type_test.BaseType) (r *type_test.BaseType, err error) {
	return nil, nil
}

func (f *Handler) GetBinaryArgument(ctx frugal.FContext, request *type_test.BaseType) (r *type_test.BaseType, err error) {
	return nil, nil
}

func (f *Handler) GetStringArgument(ctx frugal.FContext, request *type_test.BaseType) (r *type_test.BaseType, err error) {
	return nil, nil
}
