package gateway

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/Workiva/frugal/lib/gateway/marshal_test"
	"reflect"
	"strings"
	"testing"
)

var str = "foo"
var strPtr = &str

func TestJSONBuiltinMarshalUnmarshal(t *testing.T) {
	// Given
	var m JSONBuiltin
	msg := marshal_test.BaseType{StringTest: strPtr}

	// When
	buf, err := m.Marshal(&msg)
	if err != nil {
		t.Errorf("m.Marshal(%v) failed with %v; want success", &msg, err)
	}

	var got marshal_test.BaseType
	if err := json.Unmarshal(buf, &got); err != nil {
		t.Errorf("json.Unmarshal(%q, &got) failed with %v; want success", buf, err)
	}

	// Then
	if want := msg; !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v; want %v", &got, &want)
	}
}

func TestJSONBuiltinMarshalField(t *testing.T) {
	var m JSONBuiltin
	for _, fixture := range thriftTypes {
		buf, err := m.Marshal(fixture.data)
		if err != nil {
			t.Errorf("m.Marshal(%v) failed with %v; want success", fixture.data, err)
		}
		if got, want := string(buf), fixture.json; got != want {
			t.Errorf("ttype = %v; got = %q; want %q; data = %#v", fixture.ttype, got, want, fixture.data)
		}
	}
}

func TestJSONBuiltinMarshalFieldKnownErrors(t *testing.T) {
	var m JSONBuiltin
	for _, fixture := range knownErrors {
		_, err := m.Marshal(fixture.data)
		if err == nil {
			t.Errorf("m.Marshal(%v) succeeded with %v; want failure", fixture.data, err)
		}
	}
}

func TestJSONBuiltinUnmarshalField(t *testing.T) {
	var m JSONBuiltin
	for _, fixture := range thriftTypes {
		dest := reflect.New(reflect.TypeOf(fixture.data))
		if err := m.Unmarshal([]byte(fixture.json), dest.Interface()); err != nil {
			t.Errorf("m.Unmarshal(%q, dest) failed with %v; want success", fixture.json, err)
		}

		if got, want := dest.Elem().Interface(), fixture.data; !reflect.DeepEqual(got, want) {
			t.Errorf("ttype = %v, got = %#v; want = %#v; input = %q", fixture.ttype, got, want, fixture.json)
		}
	}
}

// func TestJSONBuiltinUnmarshalFieldKnownErrors(t *testing.T) {
// 	var m JSONBuiltin
// 	for _, fixture := range knownErrors {
// 		dest := reflect.New(reflect.TypeOf(fixture.data))
// 		if err := m.Unmarshal([]byte(fixture.json), dest.Interface()); err == nil {
// 			t.Errorf("m.Unmarshal(%q, dest) succeeded; want an error", fixture.json)
// 		}
// 	}
// }

func TestJSONBuiltinEncoder(t *testing.T) {
	var m JSONBuiltin
	msg := marshal_test.BaseType{
		StringTest: strPtr,
	}

	var buf bytes.Buffer
	enc := m.NewEncoder(&buf)
	if err := enc.Encode(&msg); err != nil {
		t.Errorf("enc.Encode(%v) failed with %v; want success", &msg, err)
	}

	var got marshal_test.BaseType
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Errorf("json.Unmarshal(%q, &got) failed with %v; want success", buf.String(), err)
	}
	if want := msg; !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v; want %v", &got, &want)
	}
}

func TestJSONBuiltinEncoderFields(t *testing.T) {
	var m JSONBuiltin
	for _, fixture := range thriftTypes {
		var buf bytes.Buffer
		enc := m.NewEncoder(&buf)
		if err := enc.Encode(fixture.data); err != nil {
			t.Errorf("enc.Encode(%#v) failed with %v; want success", fixture.data, err)
		}

		if got, want := buf.String(), fixture.json+"\n"; got != want {
			t.Errorf("got = %q; want %q; data = %#v", got, want, fixture.data)
		}
	}
}

func TestJSONBuiltinDecoder(t *testing.T) {
	var (
		m   JSONBuiltin
		got marshal_test.BaseType

		data = `{"StringTest": "foo"}`
	)
	r := strings.NewReader(data)
	dec := m.NewDecoder(r)
	if err := dec.Decode(&got); err != nil {
		t.Errorf("m.Unmarshal(&got) failed with %v; want success", err)
	}

	want := marshal_test.BaseType{
		StringTest: strPtr,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v; want = %v", &got, &want)
	}
}

func TestJSONBuiltinDecoderFields(t *testing.T) {
	var m JSONBuiltin
	for _, fixture := range thriftTypes {
		r := strings.NewReader(fixture.json)
		dec := m.NewDecoder(r)
		dest := reflect.New(reflect.TypeOf(fixture.data))
		if err := dec.Decode(dest.Interface()); err != nil {
			t.Errorf("dec.Decode(dest) failed with %v; want success; data = %q", err, fixture.json)
		}

		if got, want := dest.Elem().Interface(), fixture.data; !reflect.DeepEqual(got, want) {
			t.Errorf("ttype = %v; got = %v; want = %v; input = %q", fixture.ttype, got, want, fixture.json)
		}
	}
}

var (
	thriftTypes = []struct {
		ttype string
		data  interface{}
		json  string
	}{
		{ttype: "boolean", data: true, json: "true"},
		{ttype: "byte", data: int8(-1), json: "-1"},
		{ttype: "int16", data: int16(-1), json: "-1"},
		{ttype: "int32", data: int32(-1), json: "-1"},
		{ttype: "int64", data: int64(-1), json: "-1"},
		{ttype: "double", data: float64(-1.5), json: "-1.5"},
		{ttype: "binary",
			data: []byte("AAAAAQID"),
			json: `"` + base64.StdEncoding.EncodeToString([]byte("AAAAAQID")) + `"`},
		{ttype: "string", data: "AAAAAQID", json: `"AAAAAQID"`},
		{ttype: "enum", data: marshal_test.EnumType_ANOPTION, json: `"ANOPTION"`},
		{ttype: "struct", data: marshal_test.BaseType{StringTest: strPtr}, json: `{"stringTest":"foo"}`},
		{ttype: "map",
			data: marshal_test.ContainerType{
				MapTest: map[string]*marshal_test.BaseType{"foo": &marshal_test.BaseType{StringTest: strPtr}}},
			json: `{"mapTest":{"foo":{"stringTest":"foo"}}}`},
		{ttype: "list",
			data: marshal_test.ContainerType{
				ListTest: []*marshal_test.BaseType{&marshal_test.BaseType{StringTest: strPtr}}},
			json: `{"listTest":[{"stringTest":"foo"}]}`},
	}
	knownErrors = []struct {
		ttype string
		data  interface{}
		json  string
	}{
		{ttype: "set",
			data: marshal_test.ContainerType{
				SetTest: map[*marshal_test.BaseType]bool{&marshal_test.BaseType{StringTest: strPtr}: true}},
			json: `{"setTest":[{"stringTest":"foo"}]}`},
	}
)
