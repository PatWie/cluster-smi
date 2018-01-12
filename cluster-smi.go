package main

import (
	"fmt"
	"github.com/patwie/cluster-smi/cluster"
	"github.com/pebbe/zmq4"
	"github.com/vmihailenco/msgpack"
	"log"
	"os"
	"time"
)

// dummy request for REQ-ROUTER pattern
type Request struct {
	Identity string
}

func RequestUpdateMessage() (buf []byte, err error) {
	id := fmt.Sprintf("REQ %v", os.Getpid())
	req := Request{id}
	return msgpack.Marshal(&req)
}

func main() {

	request_attempts := 0

	cfg := CreateConfig()

	// ask for updates messages (REQ-ROUTER)
	request_socket, err := zmq4.NewSocket(zmq4.REQ)
	if err != nil {
		log.Fatalf("Failed open Socket ZMQ: %s\n", err.Error())
		panic(err)
	}
	defer request_socket.Close()

	SocketAddr := "tcp://" + cfg.ServerIp + ":" + cfg.ServerPortDistribute
	request_socket.Connect(SocketAddr)
	for {

		// request new update
		msg, err := RequestUpdateMessage()
		if err != nil {
			log.Fatal("request messsage error:", err)
			panic(err)
		}
		_, err = request_socket.SendBytes(msg, 0)
		if err != nil {
			log.Fatal("sending request messsage error:", err)
			panic(err)
		}

		// response from cluster-smi-server
		s, err := request_socket.RecvBytes(0)
		if err != nil {
			log.Println(err)

			time.Sleep(10 * time.Second)
			request_attempts += 1

			if request_attempts == 0 {
				panic("too many request attempts yielding an error")
			}
			continue
		}

		var clus cluster.Cluster
		err = msgpack.Unmarshal(s, &clus)
		clus.Sort()
		clus.Print()
		time.Sleep(cfg.Tick)
	}

}
