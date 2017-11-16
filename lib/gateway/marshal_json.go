package gateway

import (
	"encoding/json"
	"io"
)

// JSON is a Marshaler which marshals/unmarshals into/from json
// with the standard "encoding/json" package of Golang.
type JSON struct{}

// ContentType always returns "application/json".
func (*JSON) ContentType() string {
	return "application/json"
}

// Marshal marshals "v" into JSON
func (j *JSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal unmarshals JSON data into "v".
func (j *JSON) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// NewDecoder returns a Decoder which reads JSON stream from "r".
func (j *JSON) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (j *JSON) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}
