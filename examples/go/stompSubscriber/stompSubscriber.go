package main

import (
	"fmt"
	"sync"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/go-stomp/stomp"

	"github.com/Workiva/frugal/examples/go/gen-go/v1/music"
	"github.com/Workiva/frugal/lib/go"
)

// Run a Stomp Subscriber
func main() {
	fProtocalFactory := frugal.NewFProtocolFactory(thrift.NewTBinaryProtocolFactoryDefault())

	conn, err := stomp.Dial("tcp", "localhost:61613")
	if err != nil {
		panic(err)
	}
	defer conn.Disconnect()

	subFactory := frugal.NewFStompSubscriberTransportFactoryBuilder(conn).Build()
	provider := frugal.NewFScopeProvider(nil, subFactory, fProtocalFactory)
	subscriber := music.NewAlbumWinnersSubscriber(provider)

	// Subscribe to messages
	var wg sync.WaitGroup
	wg.Add(2)

	subscriber.SubscribeWinner(func(ctx frugal.FContext, m *music.Album) {
		fmt.Printf("received %+v : %+v\n", ctx, m)
		wg.Done()
	})
	subscriber.SubscribeContestStart(func(ctx frugal.FContext, albums []*music.Album) {
		fmt.Printf("received %+v : %+v\n", ctx, albums)
		wg.Done()
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Subscriber started...")
	wg.Wait()
}
