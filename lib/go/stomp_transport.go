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
	"fmt"
	"sync"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/go-stomp/stomp"
)

// FStompPublisherTransportFactory creates fStompPublisherTransports.
type FStompPublisherTransportFactory struct {
	conn           *stomp.Conn
	maxPublishSize int
	topicPrefix    string
}

// NewFStompPublisherTransportFactory creates an FStompPublisherTransportFactory using the
// provided stomp connection.
func NewFStompPublisherTransportFactory(conn *stomp.Conn, maxPublishSize int, topicPrefix string) *FStompPublisherTransportFactory {
	return &FStompPublisherTransportFactory{conn: conn, maxPublishSize: maxPublishSize, topicPrefix: topicPrefix}
}

// GetTransport creates a new stomp FPublisherTransport.
func (m *FStompPublisherTransportFactory) GetTransport() FPublisherTransport {
	return NewStompFPublisherTransport(m.conn, m.maxPublishSize, m.topicPrefix)
}

// fStompPublisherTransport implements FPublisherTransport.
type fStompPublisherTransport struct {
	conn           *stomp.Conn
	maxPublishSize int
	topicPrefix    string
}

// NewStompFPublisherTransport creates a new FPublisherTransport which is used for
// publishing using stomp protocol with scopes.
func NewStompFPublisherTransport(conn *stomp.Conn, maxPublishSize int, topicPrefix string) FPublisherTransport {
	return &fStompPublisherTransport{conn: conn, maxPublishSize: maxPublishSize, topicPrefix: topicPrefix}
}

// Open initializes the transport.
func (m *fStompPublisherTransport) Open() error {
	if m.conn == nil {
		return thrift.NewTTransportException(TRANSPORT_EXCEPTION_NOT_OPEN, "frugal: mq transport not open")
	}
	return nil
}

// IsOpen returns true if the transport is open, false otherwise.
func (m *fStompPublisherTransport) IsOpen() bool {
	return m.conn != nil
}

// Close closes the transport.
func (m *fStompPublisherTransport) Close() error {
	return nil
}

// GetPublishSizeLimit returns the maximum allowable size of a payload
// to be published. A non-positive number is returned to indicate an
// unbounded allowable size.
func (m *fStompPublisherTransport) GetPublishSizeLimit() uint {
	return uint(m.maxPublishSize)
}

// Publish sends the given payload with the transport.
func (m *fStompPublisherTransport) Publish(topic string, data []byte) error {
	if !m.IsOpen() {
		return thrift.NewTTransportException(TRANSPORT_EXCEPTION_NOT_OPEN, "frugal: stomp transport not open")
	}

	if len(data) > m.maxPublishSize {
		return thrift.NewTTransportException(
			TRANSPORT_EXCEPTION_REQUEST_TOO_LARGE,
			fmt.Sprintf("Message exceeds %d bytes, was %d bytes", m.maxPublishSize, len(data)))
	}

	destination := m.formatStompPublishTopic(topic)
	if err := m.conn.Send(destination, "application/octet-stream", data); err != nil {
		return thrift.NewTTransportExceptionFromError(err)
	}
	return nil
}

func (m *fStompPublisherTransport) formatStompPublishTopic(topic string) string {
	return fmt.Sprintf("/topic/%s%s%s", m.topicPrefix, frugalPrefix, topic)
}

// FStompSubscribeTransportFactory creates fStompSubscriberTransports.
type FStompSubscribeTransportFactory struct {
	conn           *stomp.Conn
	queue          string
	consumerPrefix string
}

// NewFStompSubscriberTransportFactory creates FStompSubscribeTransportFactory with the given stomp
// connection and consumer name.
func NewFStompSubscriberTransportFactory(conn *stomp.Conn, consumerPrefix string) *FStompSubscribeTransportFactory {
	return &FStompSubscribeTransportFactory{conn: conn, consumerPrefix: consumerPrefix}
}

// NewFStompSubscriberTransportFactory creates FStompSubscribeTransportFactory with the given stomp
// connection, consumer name and topic.
func NewFStompSubscriberTransportFactoryWithQueue(conn *stomp.Conn, consumerPrefix string, queue string) *FStompSubscribeTransportFactory {
	return &FStompSubscribeTransportFactory{conn: conn, consumerPrefix: consumerPrefix, queue: queue}
}

// GetTransport creates a new fStompSubscriberTransport.
func (m *FStompSubscribeTransportFactory) GetTransport() FSubscriberTransport {
	return NewStompFSubscriberTransportWithQueue(m.conn, m.consumerPrefix, m.queue)
}

// fStompSubscriberTransport implements FSubscriberTransport.
type fStompSubscriberTransport struct {
	conn           *stomp.Conn
	consumerPrefix string
	topic          string
	sub            *stomp.Subscription
	openMu         sync.RWMutex
	isSubscribed   bool
	callback       FAsyncCallback
	stopC          chan bool
}

// NewStompFSubscriberTransport creates a new FSubscriberTransport which is used for
// pub/sub.
func NewStompFSubscriberTransport(conn *stomp.Conn, consumerPrefix string) FSubscriberTransport {
	return &fStompSubscriberTransport{conn: conn, consumerPrefix: consumerPrefix}
}

// NewStompFSubscriberTransport creates a new FSubscriberTransport which is used for
// pub/sub with a topic.
func NewStompFSubscriberTransportWithQueue(conn *stomp.Conn, consumerPrefix string, topic string) FSubscriberTransport {
	return &fStompSubscriberTransport{conn: conn, consumerPrefix: consumerPrefix, topic: topic}
}

// Subscribe sets the subscribe topic and opens the transport.
func (m *fStompSubscriberTransport) Subscribe(topic string, callback FAsyncCallback) error {
	m.openMu.Lock()
	defer m.openMu.Unlock()

	if m.conn == nil {
		return thrift.NewTTransportException(TRANSPORT_EXCEPTION_NOT_OPEN, "frugal: stomp transport not open")
	}

	if m.isSubscribed {
		return thrift.NewTTransportException(TRANSPORT_EXCEPTION_ALREADY_OPEN, "frugal: stomp transport already has a subscription")
	}

	if topic == "" {
		return thrift.NewTTransportException(TRANSPORT_EXCEPTION_UNKNOWN, "frugal: stomp transport cannot subscribe to empty topic")
	}

	destination := fmt.Sprintf("/queue/%s%s", m.consumerPrefix, topic)
	sub, err := m.conn.Subscribe(destination, stomp.AckClientIndividual)
	if err != nil {
		return thrift.NewTTransportExceptionFromError(err)
	}
	m.stopC = make(chan bool, 1)
	m.sub = sub
	m.isSubscribed = true
	m.callback = callback
	go m.processMessages()
	return nil
}

// IsSubscribed returns true if the transport is subscribed to a topic, false
// otherwise.
func (m *fStompSubscriberTransport) IsSubscribed() bool {
	m.openMu.RLock()
	defer m.openMu.RUnlock()
	return m.conn != nil && m.isSubscribed
}

// Unsubscribe unsubscribes from the destination.
func (m *fStompSubscriberTransport) Unsubscribe() error {
	m.openMu.Lock()
	defer m.openMu.Unlock()
	if !m.isSubscribed {
		return nil
	}

	close(m.stopC)
	if err := m.sub.Unsubscribe(); err != nil {
		return thrift.NewTTransportExceptionFromError(err)
	}

	m.isSubscribed = false
	m.callback = nil
	return nil
}

// Processes messages from subscription channel with the given FAsyncCallback.
func (m *fStompSubscriberTransport) processMessages() {
	stopC := m.stopC
	for {
		select {
		case <-stopC:
			return
		case message, ok := <-m.sub.C:
			if !ok {
				logger().Errorf("frugal: error processing subscription messages, message channel closed unexpectedly")
				return
			}

			if len(message.Body) < 4 {
				continue
			}

			transport := &thrift.TMemoryBuffer{Buffer: bytes.NewBuffer(message.Body[4:])}
			if err := m.callback(transport); err != nil {
				logger().Warn("frugal: error executing callback: ", err)
				continue
			}

			go m.ackMessage(message)
		}
	}
}

// Acknowledges the stomp message.
func (m *fStompSubscriberTransport) ackMessage(message *stomp.Message) {
	if err := m.conn.Ack(message); err != nil {
		logger().Errorf("frugal: error acking mq message: ", err.Error())
	}
}
