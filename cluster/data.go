package cluster

import (
	"fmt"
	"github.com/apcera/termtables"
	"sort"
	"strconv"
)

type ByName []Node

func (a ByName) Len() int      { return len(a) }
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

type Memory struct {
	Used       int64 `json:"used"`
	Free       int64 `json:"free"`
	Total      int64 `json:"total"`
	Percentage int   `json:"percentage"`
}

type Device struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Utilization       int    `json:"utilization"`
	MemoryUtilization Memory `json:"memory"`
}

type Node struct {
	Name    string   `json:"name"`
	Devices []Device `json:"devices"`
}

func (n *Node) Print() {
	fmt.Println(n.Name)
	for _, device := range n.Devices {
		fmt.Println(device.Name)
	}
}

type Cluster struct {
	Nodes []Node `json:"nodes"`
}

func (c *Cluster) Sort() {
	sort.Sort(ByName(c.Nodes))
}

func (c *Cluster) Print() {

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
