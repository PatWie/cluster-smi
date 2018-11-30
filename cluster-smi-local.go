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

var clus cluster.Cluster

func main() {

	// load ports and ip-address
	cfg := LoadConfig()

	showTimePtr := flag.Bool("t", false, "show time of events")
	showExtendedPtr := flag.Bool("e", false, "extended view")
	showProcessesPtr := flag.Bool("p", false, "verbose process information")
	showDetailPtr := flag.Bool("d", false, "detail view with fan, temp, and power info")
	nodeRegex := flag.String("n", ".", "match node-names with regex for display information "+
		"(if not specified, all nodes will be shown)")
	usernameFilter := flag.String("u", "", "show all information only for specific user")
	useColor := flag.Bool("color", true, "use colored output")
	flag.Parse()

	if err := nvml.InitNVML(); err != nil {
		log.Fatalf("Failed initializing NVML: %s\n", err.Error())
	}
	defer nvml.ShutdownNVML()

	node := cluster.Node{}
	InitNode(&node)

	clus.Nodes = append(clus.Nodes, node)

	log.Println("Cluster-SMI-Local is active. Press CTRL+C to shut down.")

	for {
		FetchNode(&clus.Nodes[0])

		if *usernameFilter != "" {
			clus = cluster.FilterByUser(clus, *usernameFilter)
		}

		clus.FilterNodes(*nodeRegex)
		clus.Print(*showProcessesPtr, *showTimePtr, cfg.Timeout, *useColor, *showExtendedPtr, *showDetailPtr)
		time.Sleep(time.Duration(cfg.Tick) * time.Second)
	}

}
