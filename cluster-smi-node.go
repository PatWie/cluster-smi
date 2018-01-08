package main

import (
	"github.com/patwie/cluster-smi/nvml"
	"github.com/pebbe/zmq4"
	"github.com/vmihailenco/msgpack"
	"log"
	"time"
)

var node Node

func main() {

	cfg := CreateConfig()

	if err := nvml.InitNVML(); err != nil {
		log.Fatalf("Failed initializing NVML: %s\n", err.Error())
	}
	defer nvml.ShutdownNVML()

	SocketAddr := "tcp://" + cfg.ServerIp + ":" + cfg.ServerPortGather
	log.Println("Now pushing to", SocketAddr)
	socket, err := zmq4.NewSocket(zmq4.PUSH)
	if err != nil {
		panic(err)
	}
	defer socket.Close()
	socket.Connect(SocketAddr)

	node := &Node{}
	node.Init()

	log.Println("Cluster-SMI-Node is active. Press CTRL+C to shut down.")
	for _ = range time.Tick(cfg.Tick) {
		node.Fetch()

		// encode data
		msg, err := msgpack.Marshal(&node)
		if err != nil {
			log.Fatal("encode error:", err)
			panic(err)
		}

		// send data
		socket.SendBytes(msg, 0)
	}

}
