package gateway

import (
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestDefaultMarshalerForRequest(t *testing.T) {
	// Given
	r, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf(`http.NewRequest("GET", "http://example.com", nil) failed with %v; want success`, err)
	}
	r.Header.Set("Accept", "application/x-out")
	r.Header.Set("Content-Type", "application/x-in")

	registry := NewMarshalerRegistry()

	// When
	in, out := registry.MarshalerForRequest(r)

	// Then
	if _, ok := in.(*JSONBuiltin); !ok {
		t.Errorf("in = %#v; want a runtime.JSONBuiltin", in)
	}
	if _, ok := out.(*JSONBuiltin); !ok {
		t.Errorf("out = %#v; want a runtime.JSONBuiltin", in)
	}
}

func TestMarshalerForRequest(t *testing.T) {
	r, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf(`http.NewRequest("GET", "http://example.com", nil) failed with %v; want success`, err)
	}
	r.Header.Set("Accept", "application/x-out")
	r.Header.Set("Content-Type", "application/x-in")

	var marshalers [3]dummyMarshaler
	specs := []struct {
		mimeType string
		wantIn   Marshaler
		wantOut  Marshaler
	}{
		{
			mimeType: MIMEWildcard,
			wantIn:   &marshalers[0],
			wantOut:  &marshalers[0],
		},
		{
			mimeType: "application/x-in",
			wantIn:   &marshalers[1],
			wantOut:  &marshalers[0],
		},
		{
			mimeType: "application/x-out",
			wantIn:   &marshalers[1],
			wantOut:  &marshalers[2],
		},
	}
	for i, spec := range specs {
		registry := NewMarshalerRegistry()
		for _, s := range specs[:i+1] {
			registry.WithMarshaler(s.mimeType, s.wantOut)
		}

		in, out := registry.MarshalerForRequest(r)
		if got, want := in, spec.wantIn; got != want {
			t.Errorf("in = %#v; want %#v", got, want)
		}
		if got, want := out, spec.wantOut; got != want {
			t.Errorf("out = %#v; want %#v", got, want)
		}
	}
}

type dummyMarshaler struct{}

func (dummyMarshaler) ContentType() string { return "" }
func (dummyMarshaler) Marshal(interface{}) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (dummyMarshaler) Unmarshal([]byte, interface{}) error {
	return errors.New("not implemented")
}

func (dummyMarshaler) NewDecoder(r io.Reader) Decoder {
	return dummyDecoder{}
}
func (dummyMarshaler) NewEncoder(w io.Writer) Encoder {
	return dummyEncoder{}
}

type dummyDecoder struct{}

func (dummyDecoder) Decode(interface{}) error {
	return errors.New("not implemented")
}

type dummyEncoder struct{}

func (dummyEncoder) Encode(interface{}) error {
	return errors.New("not implemented")
}
