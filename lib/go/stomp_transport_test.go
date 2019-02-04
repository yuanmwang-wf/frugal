package frugal

import (
	"net"
	"testing"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/go-stomp/stomp"
	"github.com/go-stomp/stomp/server"
	"github.com/stretchr/testify/assert"
)

// Ensures stomp transport is able to open and close.
func TestStompPublisherOpenPublish(t *testing.T) {
	// starts a tcp server.
	l, _ := net.Listen("tcp", ":0")
	defer func() { l.Close() }()
	go server.Serve(l)

	// creates a tcp connection
	conn, err := net.Dial("tcp", l.Addr().String())
	assert.Nil(t, err)
	defer conn.Close()

	// creates stomp client
	client, err := stomp.Connect(conn)
	assert.Nil(t, err)

	publisherTransport := newStompFPublisherTransport(client, 32*1024*1024, "VirtualTopic.", nil, 3)
	err = publisherTransport.Open()
	assert.Nil(t, err)
	assert.True(t, publisherTransport.IsOpen())
	assert.Equal(t, publisherTransport.GetPublishSizeLimit(), uint(32*1024*1024))

	err = publisherTransport.Close()
	assert.Nil(t, err)
}

// Ensures stomp transport is able to publish to the expected topic.
func TestStompPublisherPublish(t *testing.T) {
	workC := make(chan *stomp.Message)

	l, _ := net.Listen("tcp", ":0")
	defer func() { l.Close() }()
	go server.Serve(l)

	// start subscriber subscribing to the expected topic.
	started := make(chan bool)
	go startSubscriber(t, "/topic/VirtualTopic.frugal.test123", l.Addr().String(), started, workC)
	<-started

	conn, err := net.Dial("tcp", l.Addr().String())
	assert.Nil(t, err)
	defer conn.Close()

	client, err := stomp.Connect(conn)
	assert.Nil(t, err)
	defer client.Disconnect()

	stompTransport := newStompFPublisherTransport(client, 32*1024*1024, "VirtualTopic.", nil, 3)
	err = stompTransport.Open()
	assert.Nil(t, err)

	err = stompTransport.Publish("test123", []byte("foo"))
	assert.Nil(t, err)

	msg := <-workC
	assert.Equal(t, string(msg.Body), "foo")
}

func TestStompPublisherTransport_Reconnect(t *testing.T) {
	workC := make(chan *stomp.Message)
	started := make(chan bool)
	reconCalled := false

	recon := func(maxAttempt int) (error, *stomp.Conn) {
		l, _ := net.Listen("tcp", ":0")
		defer func() { l.Close() }()
		go server.Serve(l)

		conn, err := net.Dial("tcp", l.Addr().String())
		client, err := stomp.Connect(conn)
		go startSubscriber(t, "/topic/VirtualTopic.frugal.test123", l.Addr().String(), started, workC)
		<-started
		reconCalled = true
		return err, client
	}

	l, _ := net.Listen("tcp", ":0")
	defer func() { l.Close() }()
	go server.Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	assert.Nil(t, err)

	client, err := stomp.Connect(conn)
	assert.Nil(t, err)

	publisherTransport := newStompFPublisherTransport(client, 32*1024*1024, "VirtualTopic.", recon, 3)
	err = publisherTransport.Open()
	assert.Nil(t, err)
	assert.True(t, publisherTransport.IsOpen())

	// close the stomp connection so frugal transport will attempt to reconnect when publishing
	err = client.Disconnect()
	assert.Nil(t, err)

	err = publisherTransport.Publish("test123", []byte("foo"))
	assert.Nil(t, err)
	assert.True(t, reconCalled)

	// Checks the message is published after reconnection
	msg := <-workC
	assert.Equal(t, string(msg.Body), "foo")

	err = publisherTransport.Close()
	assert.Nil(t, err)
}

func startSubscriber(t *testing.T, topic string, addr string, started chan bool, workC chan *stomp.Message) {
	conn, err := net.Dial("tcp", addr)
	assert.Nil(t, err)

	client, err := stomp.Connect(conn)
	assert.NoError(t, err)

	sub, err := client.Subscribe(topic, stomp.AckClientIndividual)
	assert.NoError(t, err)

	started <- true
	msg := <-sub.C
	// TODO ack is returning an error, not sure why
	//err = client.Ack(msg)
	//assert.NoError(t, err)
	workC <- msg
}

// Ensures stomp transport is able to subscribe to the expected topic and invoke callback on incoming messages.
func TestStompSubscriberSubscribe(t *testing.T) {
	started := make(chan bool, 1)

	l, _ := net.Listen("tcp", ":0")
	defer func() { l.Close() }()
	go server.Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	assert.Nil(t, err)

	client, err := stomp.Connect(conn)
	assert.Nil(t, err)

	cbCalled := make(chan bool, 1)
	cb := func(transport thrift.TTransport) error {
		cbCalled <- true
		return nil
	}
	stompTransport := newStompFSubscriberTransport(client, "Consumer.testConsumer.VirtualTopic.", true, nil, 3)
	err = stompTransport.Subscribe("testQueue", cb)
	assert.NoError(t, err)

	frame := make([]byte, 50)
	startPublisher(t, "/queue/Consumer.testConsumer.VirtualTopic.frugal.testQueue", l.Addr().String(), started, append(make([]byte, 4), frame...))
	<-started

	select {
	case <-cbCalled:
	case <-time.After(time.Second):
		assert.True(t, false, "Callback was not called")
	}
	assert.True(t, stompTransport.IsSubscribed())

	err = stompTransport.Unsubscribe()
	assert.Nil(t, err)
	assert.False(t, stompTransport.IsSubscribed())
}

func TestStompSubscriberReconnect_Subscribe(t *testing.T) {
	reconCalled := make(chan bool, 1)
	cbCalled := make(chan bool, 1)

	recon := func(maxAttempt int) (error, *stomp.Conn) {
		l, _ := net.Listen("tcp", ":0")
		defer func() { l.Close() }()
		go server.Serve(l)

		conn, err := net.Dial("tcp", l.Addr().String())
		client, err := stomp.Connect(conn)
		reconCalled <- true
		return err, client
	}

	l, _ := net.Listen("tcp", ":0")
	defer func() { l.Close() }()
	go server.Serve(l)

	cb := func(transport thrift.TTransport) error {
		cbCalled <- true
		return nil
	}

	conn, err := net.Dial("tcp", l.Addr().String())
	assert.NoError(t, err)

	client, err := stomp.Connect(conn)
	assert.NoError(t, err)

	stompTransport := newStompFSubscriberTransport(client, "Consumer.testConsumer.VirtualTopic.", true, recon, 3)
	err = stompTransport.Subscribe("testQueue", cb)
	assert.NoError(t, err)

	// Disconnect the stomp connection so frugal will attempt to reconnect
	err = client.Disconnect()
	assert.NoError(t, err)

	select {
	case <-reconCalled:
	case <-time.After(time.Second):
		assert.True(t, false, "Reconn handler was not called")
	}

	// Checks the transport is subscribed again after reconnect
	assert.True(t, stompTransport.IsSubscribed())
}

// Ensures stomp transport is able to subscribe to the expected topic and discard messages with invalid frames (size<4).
func TestStompSubscriberSubscribeDiscardsInvalidFrames(t *testing.T) {
	started := make(chan bool, 1)

	l, _ := net.Listen("tcp", ":0")
	defer func() { l.Close() }()
	go server.Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	assert.Nil(t, err)

	client, err := stomp.Connect(conn)
	assert.Nil(t, err)

	cbCalled := false
	cb := func(transport thrift.TTransport) error {
		cbCalled = true
		return nil
	}
	stompTransport := newStompFSubscriberTransport(client, "testConsumer.", false, nil, 3)
	err = stompTransport.Subscribe("testTopic", cb)
	assert.NoError(t, err)

	frame := make([]byte, 1)
	startPublisher(t, "/topic/testConsumer.frugal.testTopic", l.Addr().String(), started, append(make([]byte, 1), frame...))
	<-started

	assert.True(t, stompTransport.IsSubscribed())
	time.Sleep(10 * time.Millisecond)
	assert.False(t, cbCalled)
}

func startPublisher(t *testing.T, destination string, addr string, started chan bool, frame []byte) {
	conn, err := net.Dial("tcp", addr)
	assert.Nil(t, err)

	client, err := stomp.Connect(conn)
	assert.Nil(t, err)

	started <- true

	err = client.Send(destination, "", frame)
	assert.Nil(t, err)
}
