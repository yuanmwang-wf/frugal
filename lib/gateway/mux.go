package gateway

import (
	"fmt"
	"net/http"
	"strings"
)

// A HandlerFunc handles a specific pair of path pattern and HTTP method.
type HandlerFunc func(w http.ResponseWriter, r *http.Request, pathParams map[string]string)

// ServeMux is a request multiplexer for the Frugal gateway.
// It matches http requests to patterns and invokes the corresponding
// Frugal processor for the request.
type ServeMux struct {
	// handlers maps HTTP method to a list of handlers.
	handlers   map[string][]handler
	marshalers marshalerRegistry
}

// ServeMuxOption is an option that can be given to a ServeMux on construction.
type ServeMuxOption func(*ServeMux)

// HeaderMatcherFunc checks whether a header key should be forwarded to/from gRPC context.
type HeaderMatcherFunc func(string) (string, bool)

// NewServeMux returns a new ServeMux whose internal mapping is empty.
func NewServeMux(opts ...ServeMuxOption) *ServeMux {
	serveMux := &ServeMux{
		handlers: make(map[string][]handler),
		// forwardResponseOptions: make([]func(context.Context, http.ResponseWriter, proto.Message) error, 0),
		marshalers: makeMarshalerMIMERegistry(),
	}

	for _, opt := range opts {
		opt(serveMux)
	}

	return serveMux
}

// Handle associates "h" to the pair of HTTP method and path pattern.
func (s *ServeMux) Handle(meth string, pat Pattern, h HandlerFunc) {
	s.handlers[meth] = append(s.handlers[meth], handler{pat: pat, h: h})
}

// ServeHTTP dispatches the request to the first handler whose pattern matches to r.Method and r.Path.
func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Println("ServeMux ServeHttp")
	path := r.URL.Path
	if !strings.HasPrefix(path, "/") {
		fmt.Println("Returning: no prefix?")
		return
	}

	components := strings.Split(path[1:], "/")
	fmt.Printf("comonents: %v\n", components)
	l := len(components)
	var verb string
	if idx := strings.LastIndex(components[l-1], ":"); idx == 0 {
		fmt.Println("Returning: malformed URL")
		return
	} else if idx > 0 {
		c := components[l-1]
		components[l-1], verb = c[:idx], c[idx+1:]
	}

	fmt.Printf("components: %v\n", components)
	fmt.Printf("verb: %v\n", verb)

	if override := r.Header.Get("X-HTTP-Method-Override"); override != "" && isPathLengthFallback(r) {
		r.Method = strings.ToUpper(override)
		if err := r.ParseForm(); err != nil {
			return
		}
	}

	fmt.Println("finding handler")
	for _, h := range s.handlers[r.Method] {
		pathParams, err := h.pat.Match(components, verb)
		if err != nil {
			continue
		}

		fmt.Println("success! let match handle it")
		h.h(w, r, pathParams)
		return
	}

	// lookup other methods to handle fallback from GET to POST and
	// to determine if it is MethodNotAllowed or NotFound.

	fmt.Println("no matching handler")
	for m, handlers := range s.handlers {
		if m == r.Method {
			continue
		}
		for _, h := range handlers {
			pathParams, err := h.pat.Match(components, verb)
			if err != nil {
				continue
			}
			// X-HTTP-Method-Override is optional. Always allow fallback to POST.
			if isPathLengthFallback(r) {
				if err := r.ParseForm(); err != nil {
					fmt.Println("parse error")
					return
				}
				h.h(w, r, pathParams)
				return
			}
			fmt.Println("method not allowed")
			return
		}
	}

	fmt.Println("No matching handler: return 404")
	// GatewayErrorHandler(w, r, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func isPathLengthFallback(r *http.Request) bool {
	return r.Method == "POST" && r.Header.Get("Content-Type") == "application/x-www-form-urlencoded"
}

type handler struct {
	pat Pattern
	h   HandlerFunc
}
