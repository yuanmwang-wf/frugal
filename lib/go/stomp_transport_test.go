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

const stompPort = ":61614"

// Ensures stomp transport is able to open and close.
func TestStompPublisherOpenPublish(t *testing.T) {
	// starts a tcp server.
	l, _ := net.Listen("tcp", stompPort)
	defer func() { l.Close() }()
	go server.Serve(l)

	// creates a tcp connection
	conn, err := net.Dial("tcp", "127.0.0.1"+stompPort)
	assert.Nil(t, err)
	defer conn.Close()

	// creates stomp client
	client, err := stomp.Connect(conn)
	assert.Nil(t, err)
	defer client.Disconnect()

	amazonMq := NewStompFPublisherTransport(client, 32 * 1024 * 1024)
	err = amazonMq.Open()
	assert.Nil(t, err)
	assert.True(t, amazonMq.IsOpen())
	assert.Equal(t, amazonMq.GetPublishSizeLimit(), uint(32 * 1024 * 1024))

	err = amazonMq.Close()
	assert.Nil(t, err)
}

// Ensures Amazon Mq transport is able to publish to the expected topic.
func TestAmazonMqPublisherPublish(t *testing.T) {
	workC := make(chan *stomp.Message)

	l, _ := net.Listen("tcp", stompPort)
	defer func() { l.Close() }()
	go server.Serve(l)

	// start subscriber subscribing to the expected topic.
	started := make(chan bool)
	go startSubscriber(t, "/topic/frugal.test123", started, workC)
	<-started

	conn, err := net.Dial("tcp", "127.0.0.1"+stompPort)
	assert.Nil(t, err)
	defer conn.Close()

	client, err := stomp.Connect(conn)
	assert.Nil(t, err)
	defer client.Disconnect()

	stompTransport := NewStompFPublisherTransport(client, 32 * 1024 * 1024)
	err = stompTransport.Open()
	assert.Nil(t, err)

	err = stompTransport.Publish("test123", []byte("foo"))
	assert.Nil(t, err)

	msg := <-workC
	assert.Equal(t, string(msg.Body[:]), "foo")
}

func startSubscriber(t *testing.T, topic string, started chan bool, workC chan *stomp.Message) {
	conn, err := net.Dial("tcp", "127.0.0.1"+stompPort)
	assert.Nil(t, err)

	client, err := stomp.Connect(conn)
	assert.Nil(t, err)

	sub, err := client.Subscribe(topic, stomp.AckClientIndividual)
	assert.Nil(t, err)

	started <- true
	msg := <-sub.C
	client.Ack(msg)
	workC <- msg
}

// Ensures Amazon Mq transport is able to subscribe to the expected queue and invoke callback on incoming messages.
func TestAmazonMqSubscriberSubscribe(t *testing.T) {
	started := make(chan bool, 1)

	l, _ := net.Listen("tcp", stompPort)
	defer func() { l.Close() }()
	go server.Serve(l)

	conn, err := net.Dial("tcp", "127.0.0.1"+stompPort)
	assert.Nil(t, err)
	defer conn.Close()

	client, err := stomp.Connect(conn)
	assert.Nil(t, err)

	cbCalled := make(chan bool, 1)
	cb := func(transport thrift.TTransport) error {
		cbCalled <- true
		return nil
	}
	amazonMq := NewStompFSubscriberTransport(client, "frugal.testConsumer")
	amazonMq.Subscribe("testQueue", cb)

	frame := make([]byte, 50)
	startPublisher(t, "/queue/frugal.testConsumer.testQueue", started, append(make([]byte, 4), frame...))
	<-started

	select {
	case <-cbCalled:
	case <-time.After(time.Second):
		assert.True(t, false, "Callback was not called")
	}
	assert.True(t, amazonMq.IsSubscribed())

	err = amazonMq.Unsubscribe()
	assert.Nil(t, err)
	assert.False(t, amazonMq.IsSubscribed())
}

// Ensures Amazon Mq transport is able to subscribe to the expected queue and discard messages with invalid frames (size<4).
func TestAmazonMqSubscriberSubscribeDiscardsInvalidFrames(t *testing.T) {
	started := make(chan bool, 1)

	l, _ := net.Listen("tcp", stompPort)
	defer func() { l.Close() }()
	go server.Serve(l)

	conn, err := net.Dial("tcp", "127.0.0.1"+stompPort)
	assert.Nil(t, err)
	defer conn.Close()

	client, err := stomp.Connect(conn)
	assert.Nil(t, err)

	cbCalled := false
	cb := func(transport thrift.TTransport) error {
		cbCalled = true
		return nil
	}
	amazonMq := NewStompFSubscriberTransport(client, "frugal.testConsumer")
	amazonMq.Subscribe("testQueue", cb)

	frame := make([]byte, 1)
	startPublisher(t, "/queue/frugal.testConsumer.testQueue", started, append(make([]byte, 1), frame...))
	<-started

	assert.True(t, amazonMq.IsSubscribed())
	time.Sleep(10 * time.Millisecond)
	assert.False(t, cbCalled)
}

func startPublisher(t *testing.T, queue string, started chan bool, frame []byte) {
	conn, err := net.Dial("tcp", "127.0.0.1"+stompPort)
	assert.Nil(t, err)

	client, err := stomp.Connect(conn)
	assert.Nil(t, err)

	started <- true

	err = client.Send(queue, "", frame)
	assert.Nil(t, err)
}
