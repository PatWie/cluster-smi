package main

import (
	"fmt"
	"github.com/apcera/termtables"
	"github.com/patwie/nvidia_exporter/nvml"
	"os"
	"strconv"
)

// Cluster
func (c *Cluster) update() {
	for i, _ := range c.Nodes {
		c.Nodes[i].update()
	}
}

func (c *Cluster) print() {

	table := termtables.CreateTable()

	table.AddHeaders("Node", "Gpu", "Memory-Usage", "Mem-Util", "GPU-Util")

	for n_id, n := range c.Nodes {
		for d_id, d := range n.Devices {
			memPercent := int(d.MemoryUtilization.Used * 100 / d.MemoryUtilization.Total)

			name := ""
			if d_id == 0 {
				name = n.Name
			}

			table.AddRow(
				name,
				strconv.Itoa(d.Id)+": "+d.Name,
				strconv.FormatInt(d.MemoryUtilization.Used/1024/1024, 10)+
					"MiB / "+
					strconv.FormatInt(d.MemoryUtilization.Total/1024/1024, 10)+"MiB",
				strconv.Itoa(memPercent)+"%",
				strconv.Itoa(d.Utilization)+"%",
			)
			table.SetAlign(termtables.AlignRight, 3)

		}
		if n_id < len(c.Nodes)-1 {
			table.AddSeparator()
		}
	}
	fmt.Printf("\033[2J")
	fmt.Println(table.Render())
}

// Node
func (n *Node) init() {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	n.Name = name
	devices, _ := nvml.GetDevices()

	for i := 0; i < len(devices); i++ {
		n.Devices = append(n.Devices, Device{0, "", 0, Memory{0, 0, 0, 0}})
	}
}

func (n *Node) print() {
	fmt.Println(n.Name)
	for _, device := range n.Devices {
		fmt.Println(device.Name)
	}
}

func (n *Node) update() {

	devices, _ := nvml.GetDevices()

	for idx, device := range devices {
		meminfo, _ := device.GetMemoryInfo()
		gpuPercent, _, _ := device.GetUtilization()
		memPercent := int(meminfo.Used / meminfo.Total)
		n.Devices[idx].Id = idx
		n.Devices[idx].Name = device.DeviceName
		n.Devices[idx].Utilization = gpuPercent
		n.Devices[idx].MemoryUtilization = Memory{meminfo.Used, meminfo.Free, meminfo.Total, memPercent}

	}
}
