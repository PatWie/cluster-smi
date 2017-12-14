package main

import (
	"github.com/pebbe/zmq4"
	"github.com/vmihailenco/msgpack"
	"log"
	"sort"
)

var cluster Cluster
var allNodes map[string]Node

func main() {
	allNodes = make(map[string]Node)

	// incoming messages
	SocketAddr := "tcp://" + "*" + ":" + NodePort
	log.Printf("Open socket on %v\n", SocketAddr)
	node_socket, err := zmq4.NewSocket(zmq4.PULL)
	if err != nil {
		panic(err)
	}
	defer node_socket.Close()
	node_socket.Bind(SocketAddr)

	// outgoing messages
	SocketAddr = "tcp://" + "*" + ":" + ClientPort
	log.Printf("Open socket on %v\n", SocketAddr)
	publisher, err := zmq4.NewSocket(zmq4.PUB)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()
	publisher.Bind(SocketAddr)

	// message loop
	log.Println("Start loop")
	for {
		// read node information
		s, err := node_socket.RecvBytes(0)
		if err != nil {
			log.Println(err)
			continue
		}

		var node Node
		err = msgpack.Unmarshal(s, &node)
		if err != nil {
			panic(err)
		}

		if err != nil {
			log.Println(err)
			continue
		}

		allNodes[node.Name] = node

		cluster := Cluster{}
		for _, n := range allNodes {
			cluster.Nodes = append(cluster.Nodes, n)
		}
		sort.Sort(ByName(cluster.Nodes))

		msg, err := msgpack.Marshal(&cluster)
		publisher.SendBytes(msg, 0)

	}

}
