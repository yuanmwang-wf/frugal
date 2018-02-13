package gateway

import (
	"github.com/gorilla/mux"
)

// Router encapsulates an HTTP router for forwarding HTTP
// requests to Frugal endpoints
type Router struct {
	mux *mux.Router
}
