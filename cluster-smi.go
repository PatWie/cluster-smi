package main

import (
	"github.com/pebbe/zmq4"
	"github.com/vmihailenco/msgpack"
	"log"
	"sort"
	"time"
)

func main() {

	subscriber, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		panic(err)
	}
	defer subscriber.Close()

	SocketAddr := "tcp://" + Addr + ":" + ClientPort
	subscriber.Connect(SocketAddr)
	subscriber.SetLinger(0)
	subscriber.SetSubscribe("")
	// subscriber.SetRcvhwm(1)

	for {
		s, err := subscriber.RecvBytes(0)
		if err != nil {
			log.Println(err)
			continue
		}

		var cluster Cluster
		err = msgpack.Unmarshal(s, &cluster)
		sort.Sort(ByName(cluster.Nodes))
		cluster.print()
		time.Sleep(Tick)
	}

}

