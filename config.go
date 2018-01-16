package main

import (
	"github.com/patwie/cluster-smi/compiletimeconst"
	"strconv"
	"time"
)

type Config struct {
	ServerIp             string        // ip of cluster-smi-server
	ServerPortGather     string        // port of cluster-smi-server, which nodes send to
	ServerPortDistribute string        // port of cluster-smi-server, where clients subscribe to
	Tick                 time.Duration // tick between receiving data
	TimeoutThreshold     float64       // threshold after a node is considered/displayed as offline
}

func CreateConfig() Config {

	c := Config{}
	c.ServerIp = compiletimeconst.ServerIp
	c.ServerPortGather = compiletimeconst.PortGather
	c.ServerPortDistribute = compiletimeconst.PortDistribute

	if compiletimeconst.TimeoutThreshold == "" {
		c.TimeoutThreshold = 180
	} else {
		sec, _ := strconv.Atoi(compiletimeconst.TimeoutThreshold)
		c.TimeoutThreshold = float64(sec)
	}

	ms, _ := strconv.Atoi(compiletimeconst.Tick)
	c.Tick = time.Duration(ms) * time.Millisecond
	return c
}
