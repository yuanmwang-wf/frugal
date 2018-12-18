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
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/Workiva/frugal/lib/gopherjs/thrift"
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// _opid0_cid123[1,"ping",1,0,{}]
var pingFrame = []byte{
	0, 0, 0, 0, 29, 0, 0, 0, 5, 95, 111, 112, 105, 100, 0, 0, 0, 1, 48, 0, 0,
	0, 4, 95, 99, 105, 100, 0, 0, 0, 3, 49, 50, 51, 91, 49, 44, 34, 112, 105,
	110, 103, 34, 44, 49, 44, 48, 44, 123, 125, 93,
}

type pingProcessor struct {
	t             *testing.T
	expectedProto *FProtocol
	called        bool
	err           error
}

func (p *pingProcessor) Process(ctx FContext, iprot, oprot *FProtocol) error {
	p.called = true
	assert.Equal(p.t, p.expectedProto, iprot)
	assert.Equal(p.t, p.expectedProto, oprot)
	assert.Equal(p.t, "123", ctx.CorrelationID())
	return p.err
}

func (p *pingProcessor) AddMiddleware(ServiceMiddleware) {}

// Ensures FBaseProcessor invokes the correct FProcessorFunction and returns
// nil on success.
func TestFBaseProcessorHappyPath(t *testing.T) {
	mockTransport := new(mockTTransport)
	reads := make(chan []byte, 4)
	reads <- pingFrame[0:1]  // version
	reads <- pingFrame[1:5]  // headers size
	reads <- pingFrame[5:34] // FContext headers
	reads <- pingFrame[34:]  // request body
	mockTransport.reads = reads
	proto := &FProtocol{thrift.NewTJSONProtocol(mockTransport)}
	processor := NewFBaseProcessor()
	processorFunction := &pingProcessor{t: t, expectedProto: proto}
	processor.AddToProcessorMap("ping", processorFunction)

	assert.Nil(t, processor.Process(proto, proto))
	assert.True(t, processorFunction.called)
}

// Ensures FBaseProcessor invokes the correct FProcessorFunction and logs
// errors while returning nil.
func TestFBaseProcessorError(t *testing.T) {
	tmpLogger := logrus.New()
	var logBuf bytes.Buffer
	tmpLogger.Out = &logBuf
	oldLogger := logger()
	SetLogger(tmpLogger)
	defer func() {
		SetLogger(oldLogger)
	}()

	mockTransport := new(mockTTransport)
	reads := make(chan []byte, 4)
	reads <- pingFrame[0:1]  // version
	reads <- pingFrame[1:5]  // headers size
	reads <- pingFrame[5:34] // FContext headers
	reads <- pingFrame[34:]  // request body
	mockTransport.reads = reads
	proto := &FProtocol{thrift.NewTJSONProtocol(mockTransport)}
	processor := NewFBaseProcessor()
	err := errors.New("error")
	processorFunction := &pingProcessor{t: t, expectedProto: proto, err: err}
	processor.AddToProcessorMap("ping", processorFunction)

	assert.NoError(t, processor.Process(proto, proto))
	assert.True(t, processorFunction.called)
	assert.True(t,
		strings.Contains(
			string(logBuf.Bytes()),
			"frugal: error occurred while processing request with correlation id 123: error"))
}

// Ensures FBaseProcessor returns a TTransportException if the transport read
// fails.
func TestFBaseProcessorReadError(t *testing.T) {
	mockTransport := new(mockTTransport)
	err := errors.New("error")
	mockTransport.readError = err
	proto := &FProtocol{thrift.NewTJSONProtocol(mockTransport)}
	processor := NewFBaseProcessor()

	err = processor.Process(proto, proto)
	assert.Error(t, err)
	trErr := err.(thrift.TTransportException)
	assert.Equal(t, int(TRANSPORT_EXCEPTION_UNKNOWN), int(trErr.TypeId()))
}

// Ensures FBaseProcessor writes an UNKNOWN_METHOD TApplicationException
// response and returns nil if there is no registered FProcessorFunction.
func TestFBaseProcessorNoProcessorFunction(t *testing.T) {
	tmpLogger := logrus.New()
	var logBuf bytes.Buffer
	tmpLogger.Out = &logBuf
	oldLogger := logger()
	SetLogger(tmpLogger)
	defer func() {
		SetLogger(oldLogger)
	}()

	mockTransport := new(mockTTransport)
	reads := make(chan []byte, 4)
	reads <- pingFrame[0:1]  // version
	reads <- pingFrame[1:5]  // headers size
	reads <- pingFrame[5:34] // FContext headers
	reads <- pingFrame[34:]  // request body
	mockTransport.reads = reads
	// _opid0, cid 123
	// The ordering of opid and cid in the header is non-deterministic,
	// so cant check for equality.
	responseCtx := []byte{0, 0, 0, 0, 29, 0, 0, 0, 5, 95, 111, 112, 105, 100, 0, 0, 0, 1, 48, 0, 0, 0, 4, 95, 99, 105, 100, 0, 0, 0, 3, 49, 50, 51}
	mockTransport.On("Write", mock.Anything).Return(len(responseCtx), nil).Once()
	// [1,"ping",3,0,{"1":{"str":"Unknown function ping"},"2":{"i32":1}}]
	responseBody := []byte{
		91, 49, 44, 34, 112, 105, 110, 103, 34, 44, 51, 44, 48, 44, 123, 34,
		49, 34, 58, 123, 34, 115, 116, 114, 34, 58, 34, 85, 110, 107, 110, 111,
		119, 110, 32, 102, 117, 110, 99, 116, 105, 111, 110, 32, 112, 105, 110,
		103, 34, 125, 44, 34, 50, 34, 58, 123, 34, 105, 51, 50, 34, 58, 49,
		125, 125, 93,
	}
	mockTransport.On("Write", responseBody).Return(len(responseBody), nil).Once()
	mockTransport.On("Flush").Return(nil)
	proto := &FProtocol{thrift.NewTJSONProtocol(mockTransport)}
	processor := NewFBaseProcessor()

	assert.NoError(t, processor.Process(proto, proto))
	assert.True(t,
		strings.Contains(
			string(logBuf.Bytes()),
			"frugal: client invoked unknown function ping on request with correlation id 123"))
	mockTransport.AssertExpectations(t)
}

// Ensures FBaseProcessor writes an UNKNOWN_METHOD TApplicationException if
// there is no registered FProcessorFunction and returns an error if the write
// fails.
func TestFBaseProcessorNoProcessorFunctionWriteError(t *testing.T) {
	tmpLogger := logrus.New()
	var logBuf bytes.Buffer
	tmpLogger.Out = &logBuf
	oldLogger := logger()
	SetLogger(tmpLogger)
	defer func() {
		SetLogger(oldLogger)
	}()

	mockTransport := new(mockTTransport)
	reads := make(chan []byte, 4)
	reads <- pingFrame[0:1]  // version
	reads <- pingFrame[1:5]  // headers size
	reads <- pingFrame[5:34] // FContext headers
	reads <- pingFrame[34:]  // request body
	mockTransport.reads = reads
	// _opid0, cid 123
	// The ordering of opid and cid in the header is non-deterministic,
	// so cant check for equality.
	//responseCtx := []byte{0, 0, 0, 0, 29, 0, 0, 0, 5, 95, 111, 112, 105, 100, 0, 0, 0, 1, 48, 0, 0, 0, 4, 95, 99, 105, 100, 0, 0, 0, 3, 49, 50, 51}
	mockTransport.On("Write", mock.Anything).Return(0, errors.New("error")).Once()
	proto := &FProtocol{thrift.NewTJSONProtocol(mockTransport)}
	processor := NewFBaseProcessor()

	assert.Error(t, processor.Process(proto, proto))
	assert.True(t,
		strings.Contains(
			string(logBuf.Bytes()),
			"frugal: client invoked unknown function ping on request with correlation id 123"))
	mockTransport.AssertExpectations(t)
}

// Ensures FBaseProcessor writes an UNKNOWN_METHOD TApplicationException if
// there is no registered FProcessorFunction and returns an error if the flush
// fails.
func TestFBaseProcessorNoProcessorFunctionFlushError(t *testing.T) {
	tmpLogger := logrus.New()
	var logBuf bytes.Buffer
	tmpLogger.Out = &logBuf
	oldLogger := logger()
	SetLogger(tmpLogger)
	defer func() {
		SetLogger(oldLogger)
	}()

	mockTransport := new(mockTTransport)
	reads := make(chan []byte, 4)
	reads <- pingFrame[0:1]  // version
	reads <- pingFrame[1:5]  // headers size
	reads <- pingFrame[5:34] // FContext headers
	reads <- pingFrame[34:]  // request body
	mockTransport.reads = reads
	// _opid0, cid 123
	// The ordering of opid and cid in the header is non-deterministic,
	// so cant check for equality.
	responseCtx := []byte{0, 0, 0, 0, 29, 0, 0, 0, 5, 95, 111, 112, 105, 100, 0, 0, 0, 1, 48, 0, 0, 0, 4, 95, 99, 105, 100, 0, 0, 0, 3, 49, 50, 51}
	mockTransport.On("Write", mock.Anything).Return(len(responseCtx), nil).Once()
	// [1,"ping",3,0,{"1":{"str":"Unknown function ping"},"2":{"i32":1}}]
	responseBody := []byte{
		91, 49, 44, 34, 112, 105, 110, 103, 34, 44, 51, 44, 48, 44, 123, 34,
		49, 34, 58, 123, 34, 115, 116, 114, 34, 58, 34, 85, 110, 107, 110, 111,
		119, 110, 32, 102, 117, 110, 99, 116, 105, 111, 110, 32, 112, 105, 110,
		103, 34, 125, 44, 34, 50, 34, 58, 123, 34, 105, 51, 50, 34, 58, 49,
		125, 125, 93,
	}
	mockTransport.On("Write", responseBody).Return(len(responseBody), nil).Once()
	mockTransport.On("Flush").Return(errors.New("error"))
	proto := &FProtocol{thrift.NewTJSONProtocol(mockTransport)}
	processor := NewFBaseProcessor()

	assert.Error(t, processor.Process(proto, proto))
	assert.True(t,
		strings.Contains(
			string(logBuf.Bytes()),
			"frugal: client invoked unknown function ping on request with correlation id 123"))
	mockTransport.AssertExpectations(t)
}

// Ensures FBaseProcessor correctly returns the annotations stored on the
// processor.
func TestFBaseProcessorAnnotations(t *testing.T) {
	assert := assert.New(t)
	processor := NewFBaseProcessor()
	processor.AddToAnnotationsMap("foo", map[string]string{
		"bar":   "baz",
		"boosh": "boom",
	})
	annoMap := processor.Annotations()
	assert.Equal("baz", annoMap["foo"]["bar"])
	assert.Equal("boom", annoMap["foo"]["boosh"])

	// Verify that we cannot modify the underlying map
	delete(annoMap, "foo")
	annoMap = processor.Annotations()
	assert.Equal("baz", annoMap["foo"]["bar"])
	assert.Equal("boom", annoMap["foo"]["boosh"])
}
