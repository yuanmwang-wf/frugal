package gateway

import (
	"errors"
	"net/http"
)

// MIMEWildcard is the fallback MIME type used for requests which do not match
// a registered MIME type.
const MIMEWildcard = "*"

var (
	acceptHeader      = http.CanonicalHeaderKey("Accept")
	contentTypeHeader = http.CanonicalHeaderKey("Content-Type")

	defaultMarshaler = &JSONBuiltin{}
)

// MarshalerRegistry is a mapping from MIME types to Marshalers.
type MarshalerRegistry struct {
	mimeMap map[string]Marshaler
}

// add adds a marshaler for a case-sensitive MIME type string ("*" to match any
// MIME type).
func (m MarshalerRegistry) add(mime string, marshaler Marshaler) error {
	if len(mime) == 0 {
		return errors.New("empty MIME type")
	}

	m.mimeMap[mime] = marshaler

	return nil
}

// WithMarshaler returns a registry which associates inbound and outbound
// Marshalers to a MIME type.
func (m MarshalerRegistry) WithMarshaler(mime string, marshaler Marshaler) MarshalerRegistry {
	if err := m.add(mime, marshaler); err != nil {
		panic(err)
	}

	return m
}

// MarshalerForRequest returns the inbound/outbound marshalers for this request.
// It checks the registry on the Router for the MIME type set by the Content-Type header.
// If it isn't set (or the request Content-Type is empty), checks for "*".
// If there are multiple Content-Type headers set, choose the first one that it can
// exactly match in the registry.
// Otherwise, it follows the above logic for "*"/InboundMarshaler/OutboundMarshaler.
func (m MarshalerRegistry) MarshalerForRequest(r *http.Request) (inbound Marshaler, outbound Marshaler) {
	for _, acceptVal := range r.Header[acceptHeader] {
		if m, ok := m.mimeMap[acceptVal]; ok {
			outbound = m
			break
		}
	}

	for _, contentTypeVal := range r.Header[contentTypeHeader] {
		if m, ok := m.mimeMap[contentTypeVal]; ok {
			inbound = m
			break
		}
	}

	if inbound == nil {
		inbound = m.mimeMap[MIMEWildcard]
	}
	if outbound == nil {
		outbound = inbound
	}

	return inbound, outbound
	return nil, nil
}

// NewMarshalerRegistry returns a new registry of marshalers.
// It allows for a mapping of case-sensitive Content-Type strings to Marshaler interfaces.
//
// For example, you could allow the client to specify the use of the
// JSONBuiltin marshaler with an "application/json" Content-Type.
//
// "*" can be used to match any Content-Type.
func NewMarshalerRegistry() MarshalerRegistry {
	return MarshalerRegistry{
		mimeMap: map[string]Marshaler{
			MIMEWildcard: defaultMarshaler,
		},
	}
}
