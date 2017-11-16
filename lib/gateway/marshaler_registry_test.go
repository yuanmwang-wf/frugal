package gateway_test

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/Workiva/frugal/lib/gateway"
)

func TestMarshalerForRequest(t *testing.T) {
	r, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf(`http.NewRequest("GET", "http://example.com", nil) failed with %v; want success`, err)
	}
	r.Header.Set("Accept", "application/x-out")
	r.Header.Set("Content-Type", "application/x-in")

	mux := gateway.NewServeMux()

	in, out := gateway.MarshalerForRequest(mux, r)
	if _, ok := in.(*gateway.JSON); !ok {
		t.Errorf("in = %#v; want a JSON", in)
	}
	if _, ok := out.(*gateway.JSON); !ok {
		t.Errorf("out = %#v; want a JSON", in)
	}

	var marshalers [3]dummyMarshaler
	specs := []struct {
		opt gateway.ServeMuxOption

		wantIn  gateway.Marshaler
		wantOut gateway.Marshaler
	}{
		{
			opt:     gateway.WithMarshalerOption(gateway.MIMEWildcard, &marshalers[0]),
			wantIn:  &marshalers[0],
			wantOut: &marshalers[0],
		},
		{
			opt:     gateway.WithMarshalerOption("application/x-in", &marshalers[1]),
			wantIn:  &marshalers[1],
			wantOut: &marshalers[0],
		},
		{
			opt:     gateway.WithMarshalerOption("application/x-out", &marshalers[2]),
			wantIn:  &marshalers[1],
			wantOut: &marshalers[2],
		},
	}
	for i, spec := range specs {
		var opts []gateway.ServeMuxOption
		for _, s := range specs[:i+1] {
			opts = append(opts, s.opt)
		}
		mux = gateway.NewServeMux(opts...)

		in, out = gateway.MarshalerForRequest(mux, r)
		if got, want := in, spec.wantIn; got != want {
			t.Errorf("in = %#v; want %#v", got, want)
		}
		if got, want := out, spec.wantOut; got != want {
			t.Errorf("out = %#v; want %#v", got, want)
		}
	}

	r.Header.Set("Content-Type", "application/x-another")
	in, out = gateway.MarshalerForRequest(mux, r)
	if got, want := in, &marshalers[1]; got != want {
		t.Errorf("in = %#v; want %#v", got, want)
	}
	if got, want := out, &marshalers[0]; got != want {
		t.Errorf("out = %#v; want %#v", got, want)
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

func (dummyMarshaler) NewDecoder(r io.Reader) gateway.Decoder {
	return dummyDecoder{}
}
func (dummyMarshaler) NewEncoder(w io.Writer) gateway.Encoder {
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
