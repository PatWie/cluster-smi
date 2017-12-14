package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
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
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Node", "Gpu", "Memory-Usage", "GPU-Util"})
	for _, n := range c.Nodes {
		for _, d := range n.Devices {
			table.Append([]string{
				n.Name,
				strconv.Itoa(d.Id) + ": " + d.Name,
				strconv.FormatInt(d.MemoryUtilization.Used/1024/1024, 10) +
					"MiB / " +
					strconv.FormatInt(d.MemoryUtilization.Total/1024/1024, 10) + "MiB" +
					" (" + strconv.Itoa(d.MemoryUtilization.Percentage) + "%)",
				strconv.Itoa(d.Utilization) + "%",
			})
		}
	}
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	fmt.Printf("\033[2J")
	table.Render()
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
		gpuPercent, memPercent, _ := device.GetUtilization()

		n.Devices[idx].Id = idx
		n.Devices[idx].Name = device.DeviceName
		n.Devices[idx].Utilization = gpuPercent
		n.Devices[idx].MemoryUtilization = Memory{meminfo.Used, meminfo.Free, meminfo.Total, memPercent}

	}
}
