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

package crossrunner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	// Default timeout in seconds for client/server configurations without a defined timeout
	DefaultTimeout     = 7
	TestFailure        = 101
	CrossrunnerFailure = 102
)

// client/server level options defined in tests.json
// this is useful if there is option supported by a client but not a server within a language.
type options struct {
	Command    []string
	Transports []string
	Protocols  []string
	Timeout    time.Duration
}

// language level options defined in tests.json.
type languages struct {
	Name          string  // Language name
	Client        options // client specific commands, protocols, transports, and timesouts
	Server        options // server specific commands, protocols, transports, and timesouts
	Publisher     options
	Subscriber    options
	RpcTransports []string // transports that apply to both clients and servers within a language
	Protocols     []string // protocols that apply to both clients and servers within a language
	Command       []string // command that applies to both clients and servers within a language
	Workdir       string   // working directory relative to /test/integration
}

//  Complete information required to shell out a client or server command.
type Config struct {
	Name      string
	Timeout   time.Duration
	Transport string
	Protocol  string
	Command   []string
	Workdir   string
	Logs      *os.File
}

// Matched client and server commands.
type Pair struct {
	Client     Config
	Server     Config
	ReturnCode int
	Err        error
}

func newPair(client, server Config) *Pair {
	return &Pair{
		Client: client,
		Server: server,
	}
}

// Load takes a json file of client/server definitions and returns a list of
// valid client/server pairs.
func Load(jsonFile string) ([]*Pair, []*Pair, error) {
	bytes, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, nil, err
	}

	var tests []languages

	// Unmarshal json into defined structs
	println(jsonFile)
	println(string(bytes))
	if err := json.Unmarshal(bytes, &tests); err != nil {
		return nil, nil, err
	}

	// Create empty lists of client and server configurations
	var clients []Config
	var servers []Config
	var publishers []Config
	var subscribers []Config

	// Iterate over each language to get all client/server configurations in that language
	for _, test := range tests {
		fmt.Printf("%+v\n", test.Publisher)
		fmt.Printf("%+v\n", test.Subscriber)

		// Append language level transports and protocols to client/server level
		test.Client.Transports = append(test.Client.Transports, test.RpcTransports...)
		test.Server.Transports = append(test.Server.Transports, test.RpcTransports...)
		test.Client.Protocols = append(test.Client.Protocols, test.Protocols...)
		test.Server.Protocols = append(test.Server.Protocols, test.Protocols...)

		//test.Publisher.Transports = append(test.Publisher.Transports, test.Transports...)
		//test.Subscriber.Transports = append(test.Subscriber.Transports, test.Transports...)
		test.Publisher.Protocols = append(test.Publisher.Protocols, test.Protocols...)
		test.Subscriber.Protocols = append(test.Subscriber.Protocols, test.Protocols...)

		// Get expanded list of clients/servers, using both language and Config level options
		clients = append(clients, getExpandedConfigs(test.Client, test)...)
		servers = append(servers, getExpandedConfigs(test.Server, test)...)
		publishers = append(publishers, getExpandedConfigs(test.Publisher, test)...)
		subscribers = append(subscribers, getExpandedConfigs(test.Subscriber, test)...)
	}

	// Find all valid client/server pairs
	// TODO: Accept some sort of flag(s) that would limit this list of pairs by
	// desired language(s) or other restrictions
	rpcPairs := make([]*Pair, 0)
	for _, client := range clients {
		for _, server := range servers {
			if server.Transport == client.Transport && server.Protocol == client.Protocol {
				rpcPairs = append(rpcPairs, newPair(client, server))
			}
		}
	}

	pubsubPairs := make([]*Pair, 0)
	println("len(publishers) == ", len(publishers))
	println("len(subscribers) == ", len(subscribers))
	for _, publisher := range publishers {
		for _, subscriber := range subscribers {
			println(publisher.Transport, subscriber.Transport, publisher.Protocol, subscriber.Protocol)
			if publisher.Transport == subscriber.Transport && publisher.Protocol == subscriber.Protocol {
				println("found a match!")
				pubsubPairs = append(pubsubPairs, newPair(publisher, subscriber))
			}
		}
	}

	return rpcPairs, pubsubPairs, nil
}

// getExpandedConfigs takes a client/server at the language level and the options
// associated with that client/server and returns a list of unique configs.
func getExpandedConfigs(options options, test languages) (apps []Config) {
	app := new(Config)

	// Loop through each transport and protocol to construct expanded list
	for _, transport := range options.Transports {
		for _, protocol := range options.Protocols {
			app.Name = test.Name
			app.Protocol = protocol
			app.Transport = transport
			app.Command = append(test.Command, options.Command...)
			app.Workdir = test.Workdir
			app.Timeout = DefaultTimeout * time.Second
			if options.Timeout != 0 {
				app.Timeout = options.Timeout
			}
			apps = append(apps, *app)
		}
	}
	return apps
}
