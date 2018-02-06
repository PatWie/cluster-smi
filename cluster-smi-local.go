package main

/*
This is simply reproducing "nvidia-smi" without networking.
*/

import (
	"flag"
	"github.com/patwie/cluster-smi/cluster"
	"github.com/patwie/cluster-smi/nvml"
	"log"
	"time"
)

var cls cluster.Cluster

func main() {

	// load ports and ip-address
	cfg := LoadConfig()

	showTimePtr := flag.Bool("t", false, "show time of events")
	showProcessesPtr := flag.Bool("p", false, "verbose process information")
	nodeRegex := flag.String("n", ".", "match node-names with regex for display information "+
		"(if not specified, all nodes will be shown)")
	flag.Parse()

	if err := nvml.InitNVML(); err != nil {
		log.Fatalf("Failed initializing NVML: %s\n", err.Error())
	}
	defer nvml.ShutdownNVML()

	node := cluster.Node{}
	InitNode(&node)

	cls.Nodes = append(cls.Nodes, node)

	log.Println("Cluster-SMI-Local is active. Press CTRL+C to shut down.")

	for {
		FetchNode(&cls.Nodes[0])
		cls.FilterNodes(*nodeRegex)
		cls.Print(*showProcessesPtr, *showTimePtr, cfg.Timeout)
		time.Sleep(time.Duration(cfg.Tick) * time.Second)
	}

}
