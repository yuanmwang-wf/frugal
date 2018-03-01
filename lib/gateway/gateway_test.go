// These tests require a running server and gateway before running.
//
// make gen-go
// go run exampleHttp/main.go
// go run exampleGateway/main.go
package gateway_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Workiva/frugal/lib/gateway/gen-go/gateway_test"
	"io/ioutil"
	"net/http"
	"testing"
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
