/*
 * Copyright 2017 Workiva
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package frugal

import (
	"fmt"
	"testing"
	"time"

	"github.com/Workiva/frugal/lib/gopherjs/thrift"
	"github.com/nats-io/go-nats"
	"github.com/stretchr/testify/assert"
)

// Ensures FStatelessNatsServer receives requests and sends back responses on
// the correct subject.
func TestFStatelessNatsServer(t *testing.T) {
	s := runServer(nil)
	defer s.Shutdown()
	conn, err := nats.Connect(fmt.Sprintf("nats://localhost:%d", defaultOptions.Port))
	if err != nil {
		t.Fatal(err)
	}
	processor := &processor{t}
	protoFactory := NewFProtocolFactory(thrift.NewTBinaryProtocolFactoryDefault())
	server := NewFNatsServerBuilder(conn, processor, protoFactory, []string{"foo"}).WithQueueGroup("queue").Build()
	go func() {
		assert.Nil(t, server.Serve())
	}()
	time.Sleep(10 * time.Millisecond)
	defer server.Stop()

	tr := NewFNatsTransport(conn, "foo", "bar").(*fNatsTransport)
	ctx := NewFContext("")
	assert.Nil(t, tr.Open())

	// Send a request.
	buffer := NewTMemoryOutputBuffer(0)
	proto := protoFactory.GetProtocol(buffer)
	proto.WriteRequestHeader(ctx)
	proto.WriteBinary([]byte{1, 2, 3, 4, 5})
	resultTrans, err := tr.Request(ctx, buffer.Bytes())
	assert.Nil(t, err)

	resultProto := protoFactory.GetProtocol(resultTrans)
	ctx = NewFContext("")
	err = resultProto.ReadResponseHeader(ctx)
	assert.Nil(t, err)
	resultBytes, err := resultProto.ReadBinary()
	assert.Nil(t, err)
	assert.Equal(t, "foo", string(resultBytes))
}

type processor struct {
	t *testing.T
}

func (p *processor) Process(in, out *FProtocol) error {
	ctx, err := in.ReadRequestHeader()
	if err != nil {
		return err
	}
	bytes, err := in.ReadBinary()
	if err != nil {
		return err
	}
	assert.Equal(p.t, []byte{1, 2, 3, 4, 5}, bytes)
	out.WriteResponseHeader(ctx)
	out.WriteString("foo")
	return nil
}

func (p *processor) AddMiddleware(middleware ServiceMiddleware) {}

func (p *processor) Annotations() map[string]map[string]string {
	return nil
}
