package gateway

import (
	"fmt"
	"net/http"

	"github.com/Workiva/frugal/lib/go"

	"git.apache.org/thrift.git/lib/go/thrift"
)

// HTTPStatusFromFrugalError converts a Frugal exception into a corresponding HTTP response status.
// https://github.com/Workiva/frugal/blob/master/lib/go/errors.go
func HTTPStatusFromFrugalError(err thrift.TException) int {
	switch err.(type) {
	case thrift.TTransportException:
		if e, ok := err.(thrift.TTransportException); ok {
			switch e.TypeId() {
			case frugal.TRANSPORT_EXCEPTION_UNKNOWN, frugal.TRANSPORT_EXCEPTION_NOT_OPEN, frugal.TRANSPORT_EXCEPTION_ALREADY_OPEN, frugal.TRANSPORT_EXCEPTION_END_OF_FILE:
				return http.StatusInternalServerError
			case frugal.TRANSPORT_EXCEPTION_TIMED_OUT:
				return http.StatusRequestTimeout
			case frugal.TRANSPORT_EXCEPTION_REQUEST_TOO_LARGE:
				return http.StatusBadRequest
			default:
				fmt.Println("unknown TTransportException")
				return http.StatusInternalServerError
			}
		}

	case thrift.TApplicationException:
		if e, ok := err.(thrift.TTransportException); ok {
			switch e.TypeId() {
			case frugal.APPLICATION_EXCEPTION_INVALID_MESSAGE_TYPE, frugal.APPLICATION_EXCEPTION_WRONG_METHOD_NAME, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE:
				return http.StatusBadRequest
			case frugal.APPLICATION_EXCEPTION_UNKNOWN, frugal.APPLICATION_EXCEPTION_UNKNOWN_METHOD, frugal.APPLICATION_EXCEPTION_UNSUPPORTED_CLIENT_TYPE, frugal.APPLICATION_EXCEPTION_BAD_SEQUENCE_ID, frugal.APPLICATION_EXCEPTION_MISSING_RESULT, frugal.APPLICATION_EXCEPTION_INTERNAL_ERROR, frugal.APPLICATION_EXCEPTION_PROTOCOL_ERROR, frugal.APPLICATION_EXCEPTION_INVALID_TRANSFORM, frugal.APPLICATION_EXCEPTION_INVALID_PROTOCOL:
				return http.StatusInternalServerError
			default:
				fmt.Println("unknown TTransportException")
				return http.StatusInternalServerError
			}
		}
	}

	fmt.Printf("Unknown Frugal error: %v", err)
	return http.StatusInternalServerError
}

// ErrorBody represents the default error response
type ErrorBody struct {
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}

// DefaultFrugalError is the default implementation of HTTPError.
//
// If the error is from Frugal, the function replies with the Frugal
// status code mapped to an HTTP status code using HTTPStatusFromFrugalError.
// If an unknown error occurs, it replies with http.StatusInternalServerError.
//
// The response body returned by this function is a JSON object, which
// contains a member whose key is "error" and whose value is err.Error().
// TODO: Make match API Guidelines
func DefaultFrugalError(marshaler Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "Failed to serialize error message"}`

	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType())

	// s, ok := status.FromError(err)
	// if !ok {
	// 	s = status.New(codes.Unknown, err.Error())
	// }

	// body := &errorBody{
	// 	Error: s.Message(),
	// 	Code:  int32(s.Code()),
	// }

	// buf, merr := marshaler.Marshal(body)
	// if merr != nil {
	// 	fmt.Printf("Failed to marshal error message %q: %v", body, merr)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	if _, err := io.WriteString(w, fallback); err != nil {
	// 		fmt.Printf("Failed to write response: %v", err)
	// 	}
	// 	return
	// }

	// md, ok := ServerMetadataFromContext(ctx)
	// if !ok {
	// 	fmt.Printf("Failed to extract ServerMetadata from context")
	// }

	// handleForwardResponseTrailerHeader(w, md)
	// st := HTTPStatusFromCode(s.Code())
	st := HTTPStatusFromFrugalError(err)
	fmt.Printf("Status codde: %v", st)
	// w.WriteHeader(st)
	// if _, err := w.Write(buf); err != nil {
	// 	fmt.Printf("Failed to write response: %v", err)
	// }

	// handleForwardResponseTrailer(w, md)
}

// DefaultRequestErrorHandler is the default implementation of an error handler for incoming HTTP requests.
func DefaultRequestErrorHandler(marshaler Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	f, _ := w.(http.Flusher)
	const fallback = "Problems parsing JSON"
	var status int
	var message string

	// Handle data type errors, etc.
	switch err {
	// case io.EOF:
	// 	fmt.Printf("EOF %v\n", err)

	default:
		status = http.StatusBadRequest
		message = fallback
	}

	// Serialize the error message to JSON
	body := ErrorBody{"<request-id>", message}

	buf, merr := marshaler.Marshal(body)
	if merr != nil {
		fmt.Printf("Failed to marshal error message %q: %v", body, merr)
		status = http.StatusInternalServerError
		message = fallback
	}

	// Write the response
	w.Header().Set("Content-Type", marshaler.ContentType())
	w.WriteHeader(status)
	w.Write(buf)
	f.Flush()
}

// DefaultFrugalErrorHandler is the default implementation of an error handler for incoming HTTP requests.
func DefaultFrugalErrorHandler(marshaler Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	f, _ := w.(http.Flusher)
	const fallback = "Internal Server Error"
	var status int
	var message string

	// Handle data type errors, etc.
	switch err.(type) {
	case thrift.TException:
		status = HTTPStatusFromFrugalError(err)
		// message = "" // TODO: return from above code (message as per API guidelines)
	default:
		status = http.StatusInternalServerError
		message = fallback
	}

	// Serialize the error message to JSON
	body := ErrorBody{"<request-id>", message}

	buf, merr := marshaler.Marshal(body)
	if merr != nil {
		fmt.Printf("Failed to marshal error message %q: %v", body, merr)
		status = http.StatusInternalServerError
		message = fallback
	}

	// Write the response
	w.Header().Set("Content-Type", marshaler.ContentType())
	w.WriteHeader(status)
	w.Write(buf)
	f.Flush()
}

// DefaultResponseErrorHandler is the default implementation of an error handler for outgoing HTTP responses.
func DefaultResponseErrorHandler(marshaler Marshaler, w http.ResponseWriter, r *http.Request, err error) {

}
