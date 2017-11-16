package gateway

import (
	"fmt"
	"net/http"
)

// ForwardResponseMessage forwards the response from a Frugal server to a REST client.
func ForwardResponseMessage(mux *ServeMux, marshaler Marshaler, w http.ResponseWriter, req *http.Request, resp interface{}) {
	w.Header().Set("Content-Type", marshaler.ContentType())

	buf, err := marshaler.Marshal(resp)
	if err != nil {
		fmt.Printf("Marshal error: %v", err)
		return
	}

	if _, err = w.Write(buf); err != nil {
		fmt.Printf("Failed to write response: %v", err)
	}

}
