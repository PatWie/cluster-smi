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
	selectedNodes :=flag.String("n", "", "select node(s) to display information, " +
                                         "muliple nodes are separated by comma without space " +
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
	for _ = range time.Tick(time.Duration(cfg.Tick) * time.Second) {
		FetchNode(&cls.Nodes[0])
		cls.Print(*selectedNodes, *showProcessesPtr, *showTimePtr, cfg.Timeout)
	}

}
