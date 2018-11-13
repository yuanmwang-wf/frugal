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

const frugalVirtualTopicPrefix = "frugal.VirtualTopic."

// FAmazonMqPublisherTransportFactory creates FMqPublisherTransports.
type FAmazonMqPublisherTransportFactory struct {
	conn *stomp.Conn
}

// NewFAmazonMqPublisherTransportFactory creates an FAmazonMqPublisherTransportFactory using the
// provided stomp connection.
func NewFAmazonMqPublisherTransportFactory(conn *stomp.Conn) *FAmazonMqPublisherTransportFactory {
	return &FAmazonMqPublisherTransportFactory{conn: conn}
}

// GetTransport creates a new Amazon MQ FPublisherTransport.
func (m *FAmazonMqPublisherTransportFactory) GetTransport() FPublisherTransport {
	return NewAmazonMqFPublisherTransport(m.conn)
}

// fAmazonMqPublisherTransport implements FPublisherTransport.
type fAmazonMqPublisherTransport struct {
	conn *stomp.Conn
}

// NewAmazonMqFPublisherTransport creates a new FPublisherTransport which is used for
// publishing with scopes.
func NewAmazonMqFPublisherTransport(conn *stomp.Conn) FPublisherTransport {
	return &fAmazonMqPublisherTransport{conn: conn}
}

// Open initializes the transport.
func (m *fAmazonMqPublisherTransport) Open() error {
	if m.conn == nil {
		return thrift.NewTTransportException(TRANSPORT_EXCEPTION_NOT_OPEN, "frugal: mq transport not open")
	}
	return nil
}

// IsOpen returns true if the transport is open, false otherwise.
func (m *fAmazonMqPublisherTransport) IsOpen() bool {
	return m.conn != nil
}

// Close closes the transport.
func (m *fAmazonMqPublisherTransport) Close() error {
	if err := m.conn.Disconnect(); err != nil {
		return thrift.NewTTransportException(TRANSPORT_EXCEPTION_UNKNOWN, fmt.Sprintf("frugal: error closing transport: %s", err))
	}
	m.conn = nil
	return nil
}

// GetPublishSizeLimit returns the maximum allowable size of a payload
// to be published. A non-positive number is returned to indicate an
// unbounded allowable size.
func (m *fAmazonMqPublisherTransport) GetPublishSizeLimit() uint {
	return 0
}

// Publish sends the given payload with the transport. Data are sent to virtual topics created in AmazonMq.
func (m *fAmazonMqPublisherTransport) Publish(topic string, data []byte) error {
	if !m.IsOpen() {
		return thrift.NewTTransportException(TRANSPORT_EXCEPTION_NOT_OPEN, "frugal: mq transport not open")
	}

	destination := fmt.Sprintf("/topic/%s%s", frugalVirtualTopicPrefix, topic)

	if err := m.conn.Send(destination, "text/plain", data); err != nil {
		return thrift.NewTTransportExceptionFromError(err)
	}
	return nil
}

// FAmazonMqSubscribeTransportFactory creates FAmazonMqSubscriberTransports.
type FAmazonMqSubscribeTransportFactory struct {
	conn         *stomp.Conn
	queue        string
	consumerName string
}

// NewFAmazonMqSubscriberTransportFactory creates FAmazonMqSubscribeTransportFactory with the given stomp
// connection and consumer name.
func NewFAmazonMqSubscriberTransportFactory(conn *stomp.Conn, consumerName string) *FAmazonMqSubscribeTransportFactory {
	return &FAmazonMqSubscribeTransportFactory{conn: conn, consumerName: consumerName}
}

// NewFAmazonMqSubscriberTransportFactory creates FAmazonMqSubscribeTransportFactory with the given stomp
// connection, consumer name and queue.
func NewFAmazonMqSubscriberTransportFactoryWithQueue(conn *stomp.Conn, consumerName string, queue string) *FAmazonMqSubscribeTransportFactory {
	return &FAmazonMqSubscribeTransportFactory{conn: conn, consumerName: consumerName, queue: queue}
}

// GetTransport creates a new fAmazonMqSubscriberTransport.
func (m *FAmazonMqSubscribeTransportFactory) GetTransport() FSubscriberTransport {
	return NewAmazonMqFSubscriberTransportWithQueue(m.conn, m.consumerName, m.queue)
}

// fAmazonMqSubscriberTransport implements FSubscriberTransport.
type fAmazonMqSubscriberTransport struct {
	conn         *stomp.Conn
	consumerName string
	queue        string
	sub          *stomp.Subscription
	openMu       sync.RWMutex
	isSubscribed bool
	callback     FAsyncCallback
	stopC        chan bool
}

// NewAmazonMqFSubscriberTransport creates a new FSubscriberTransport which is used for
// pub/sub.
func NewAmazonMqFSubscriberTransport(conn *stomp.Conn, consumerName string) FSubscriberTransport {
	return &fAmazonMqSubscriberTransport{conn: conn, consumerName: consumerName}
}

// NewAmazonMqFSubscriberTransport creates a new FSubscriberTransport which is used for
// pub/sub with a queue.
func NewAmazonMqFSubscriberTransportWithQueue(conn *stomp.Conn, consumerName string, queue string) FSubscriberTransport {
	return &fAmazonMqSubscriberTransport{conn: conn, consumerName: consumerName, queue: queue}
}

// Subscribe sets the subscribe queue and opens the transport.
func (m *fAmazonMqSubscriberTransport) Subscribe(queue string, callback FAsyncCallback) error {
	m.openMu.Lock()
	defer m.openMu.Unlock()

	if m.conn == nil {
		return thrift.NewTTransportException(TRANSPORT_EXCEPTION_NOT_OPEN, "frugal: mq transport not open")
	}

	if m.isSubscribed {
		return thrift.NewTTransportException(TRANSPORT_EXCEPTION_ALREADY_OPEN, "frugal: mq transport already has a subscription")
	}

	if queue == "" {
		return thrift.NewTTransportException(TRANSPORT_EXCEPTION_UNKNOWN, "frugal: mq transport cannot subscribe to empty queue")
	}

	destination := fmt.Sprintf("/queue/frugalConsumer.%s.%s", m.consumerName, queue)
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

// IsSubscribed returns true if the transport is subscribed to a queue, false
// otherwise.
func (m *fAmazonMqSubscriberTransport) IsSubscribed() bool {
	m.openMu.RLock()
	defer m.openMu.RUnlock()
	return m.conn != nil && m.isSubscribed
}

// Unsubscribe unsubscribes and disconnects the stomp connection.
func (m *fAmazonMqSubscriberTransport) Unsubscribe() error {
	m.openMu.Lock()
	defer m.openMu.Unlock()
	if !m.isSubscribed {
		return nil
	}

	m.stopC <- true
	if err := m.sub.Unsubscribe(); err != nil {
		return thrift.NewTTransportExceptionFromError(err)
	}

	if err := m.conn.Disconnect(); err != nil {
		return thrift.NewTTransportExceptionFromError(err)
	}

	m.isSubscribed = false
	m.callback = nil
	return nil
}

// Processes messages from subscription channel with the given FAsyncCallback.
func (m *fAmazonMqSubscriberTransport) processMessages() {
	for {
		select {
		case <-m.stopC:
			return
		case message, ok := <-m.sub.C:
			if !ok {
				logger().Errorf("messaging_sdk: error processing subscription messages, message channel closed unexpectedly")
				return
			}

			if len(message.Body) < 4 {
				continue
			}

			transport := &thrift.TMemoryBuffer{Buffer: bytes.NewBuffer(message.Body[4:])}
			if m.callback(transport) != nil {
				continue
			}

			go m.ackMessage(message)
		}
	}
}

// Acknowledges the stomp message.
func (m *fAmazonMqSubscriberTransport) ackMessage(message *stomp.Message) {
	if err := m.conn.Ack(message); err != nil {
		logger().Errorf("messaging_sdk: error acking mq message: ", err.Error())
	}
}
