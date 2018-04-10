// These tests require a running server and gateway before running.
//
// make gen-go
// go run exampleHttp/main.go
// go run exampleGateway/main.go
package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Workiva/frugal/lib/gateway/gen-go/gateway_test"
	"io/ioutil"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func unmarshalContainer(body []byte) (*gateway_test.ContainerType, error) {
	var s = new(gateway_test.ContainerType)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}

func TestStringPathParameter(t *testing.T) {
	payload := []byte(`{}`)
	rs, err := http.Post("http://localhost:5000/v1/container/", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		panic(err.Error())
	}

	c, err := unmarshalContainer([]byte(body))
	v := c.ListTest[0].StringTest
	if v != "container" {
		t.Error("Expected container, got ", v)
	}
}

func TestStringQueryOverridesPathParameter(t *testing.T) {
	payload := []byte(`{}`)
	rs, err := http.Post("http://localhost:5000/v1/container/?differentString=foo", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		panic(err.Error())
	}

	c, err := unmarshalContainer([]byte(body))
	v := c.ListTest[0].StringTest
	if v != "foo" {
		t.Error("Expected foo, got ", v)
	}
}

func TestJSONPayload(t *testing.T) {
	payload := []byte(`{"boolTest":true}`)
	rs, err := http.Post("http://localhost:5000/v1/container/", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		panic(err.Error())
	}

	c, err := unmarshalContainer([]byte(body))
	v := c.ListTest[0].BoolTest
	if !*v {
		t.Error("Expected true, got ", v)
	}
}

func TestBaseTypeSerialization(t *testing.T) {
	//1: optional bool boolTest;
	//2: optional byte byteTest;
	//3: optional i16 i16Test;
	//4: optional i32 i32Test;
	//5: optional i64 i64Test;
	//6: optional double doubleTest;
	//7: optional binary binaryTest;
	//8: string stringTest (http.jsonProperty="differentString")
	payload := []byte(`{
	"boolTest": true,
	"byteTest": 1,
	"i16Test": 2,
	"i32Test": 4,
	"i64Test": 8,
	"doubleTest": 1.24
}`)

	rs, err := http.Post("http://localhost:5000/base/container/", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		panic(err)
	}
	println(string(body))
	var baseType gateway_test.BaseType
	json.Unmarshal(body, &baseType)

	fmt.Printf("received %+v\n", baseType)
	assert.Equal(t, true, *baseType.BoolTest)
	assert.Equal(t, int8(1), *baseType.ByteTest)
	assert.Equal(t, int16(2), *baseType.I16Test)
	assert.Equal(t, int32(4), *baseType.I32Test)
	assert.Equal(t, int64(8), *baseType.I64Test)
	assert.Equal(t, float64(1.24), *baseType.DoubleTest)
}
