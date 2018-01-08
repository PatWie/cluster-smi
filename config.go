package main

import (
	"github.com/patwie/cluster-smi/compiletimeconst"
	"log"
	"strconv"
	"time"
)

type Config struct {
	ServerIp             string        // ip of cluster-smi-server
	ServerPortGather     string        // port of cluster-smi-server, which nodes send to
	ServerPortDistribute string        // port of cluster-smi-server, where clients subscribe to
	Tick                 time.Duration // tick between receiving data
}

func CreateConfig() Config {

	log.Println(compiletimeconst.ServerIp)
	log.Println(compiletimeconst.PortGather)
	log.Println(compiletimeconst.PortDistribute)
	log.Println(compiletimeconst.Tick)

	c := Config{}
	c.ServerIp = compiletimeconst.ServerIp
	c.ServerPortGather = compiletimeconst.PortGather
	c.ServerPortDistribute = compiletimeconst.PortDistribute

	ms, _ := strconv.Atoi(compiletimeconst.Tick)
	c.Tick = time.Duration(ms) * time.Millisecond
	return c
}
