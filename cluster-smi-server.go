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

	var cfg Config
	cfg.ReadConfig("cluster-smi.yml")

	// incoming messages (Push-Pull)
	SocketAddr := "tcp://" + "*" + ":" + cfg.ServerPortGather
	log.Println("Now listening on", SocketAddr)
	node_socket, err := zmq4.NewSocket(zmq4.PULL)
	if err != nil {
		panic(err)
	}
	defer node_socket.Close()
	node_socket.Bind(SocketAddr)

	// outgoing messages (Pub-Sub)
	SocketAddr = "tcp://" + "*" + ":" + cfg.ServerPortDistribute
	log.Println("Now publishing to", SocketAddr)
	publisher, err := zmq4.NewSocket(zmq4.PUB)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()
	publisher.Bind(SocketAddr)

	// message loop
	log.Println("Cluster-SMI-Server is active. Press CTRL+C to shut down.")
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
			log.Println(err)
			// panic(err)
			continue
		}

		// update information
		allNodes[node.Name] = node

		// rebuild cluster struct
		cluster := Cluster{}
		for _, n := range allNodes {
			cluster.Nodes = append(cluster.Nodes, n)
		}
		sort.Sort(ByName(cluster.Nodes))

		// send cluster information
		msg, err := msgpack.Marshal(&cluster)
		publisher.SendBytes(msg, 0)

	}

}
