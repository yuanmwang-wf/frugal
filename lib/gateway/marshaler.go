package gateway

import (
	"io"
)

// Marshaler defines a conversion between a byte sequence and Frugal payloads / fields.
type Marshaler interface {
	// Marshal converts a struct "v" into a byte sequence.
	Marshal(v interface{}) ([]byte, error)

	// Unmarshal converts a byte sequence "data" into pointer value for struct "v".
	Unmarshal(data []byte, v interface{}) error

	// NewDecoder returns a Decoder which reads a byte sequence from "r".
	NewDecoder(r io.Reader) Decoder

	// NewEncoder returns an Encoder which writes bytes sequence into "w".
	NewEncoder(w io.Writer) Encoder

	// ContentType returns the Content-Type which this marshaler is responsible for.
	ContentType() string
}

// Decoder decodes a byte sequence
type Decoder interface {
	Decode(v interface{}) error
}

// Encoder encodes Frugal payloads / fields into byte sequence.
type Encoder interface {
	Encode(v interface{}) error
}

// DecoderFunc adapts an decoder function into Decoder.
type DecoderFunc func(v interface{}) error

// Decode delegates invocations to the underlying function itself.
func (f DecoderFunc) Decode(v interface{}) error { return f(v) }

// EncoderFunc adapts an encoder function into Encoder
type EncoderFunc func(v interface{}) error

// Encode delegates invocations to the underlying function itself.
func (f EncoderFunc) Encode(v interface{}) error { return f(v) }
