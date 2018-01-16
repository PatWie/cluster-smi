package main

import (
	"github.com/patwie/cluster-smi/cluster"
	"github.com/patwie/cluster-smi/nvml"
	"os"
	"time"
)

// Cluster
func FetchCluster(c *cluster.Cluster) {
	for i, _ := range c.Nodes {
		FetchNode(&c.Nodes[i])
	}
}

// Node
func InitNode(n *cluster.Node) {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	n.Name = name
	n.Time = time.Now()

	devices, _ := nvml.GetDevices()

	for i := 0; i < len(devices); i++ {
		n.Devices = append(n.Devices, cluster.Device{0, "", 0, cluster.Memory{0, 0, 0, 0}, nil})
	}
}

func FetchNode(n *cluster.Node) {

	devices, _ := nvml.GetDevices()
	n.Time = time.Now()

	for idx, device := range devices {
		meminfo, _ := device.GetMemoryInfo()
		gpuPercent, _, _ := device.GetUtilization()
		memPercent := int(meminfo.Used / meminfo.Total)

		// read processes
		deviceProcs, err := device.GetProcessInfo()
		if err != nil {
			panic(err)
		}

		var p []cluster.Process
		for i := 0; i < len(deviceProcs); i++ {
			if int(deviceProcs[i].Pid) == 0 {
				continue
			}

			p = append(p, cluster.Process{
				Pid:           deviceProcs[i].Pid,
				UsedGpuMemory: deviceProcs[i].UsedGpuMemory,
			})
		}

		n.Devices[idx].Id = idx
		n.Devices[idx].Name = device.DeviceName
		n.Devices[idx].Utilization = gpuPercent
		n.Devices[idx].MemoryUtilization = cluster.Memory{meminfo.Used, meminfo.Free, meminfo.Total, memPercent}
		n.Devices[idx].Processes = p

	}
}
