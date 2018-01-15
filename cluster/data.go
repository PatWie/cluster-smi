package cluster

import (
	"fmt"
	"github.com/apcera/termtables"
	"sort"
	"time"
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
	Name    string    `json:"name"`
	Devices []Device  `json:"devices"`
	Time    time.Time `json:"time"`
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

func (c *Cluster) Print(verbose bool) {

	table := termtables.CreateTable()

	tableHeader := []interface{}{"Node", "Gpu", "Memory-Usage", "Mem-Util", "GPU-Util"}

	if verbose {
		tableHeader = append(tableHeader, "Last Seen", "Timeout")
	}
	table.AddHeaders(tableHeader...)

	now := time.Now()

	for n_id, n := range c.Nodes {

		timeout := now.Sub(n.Time).Seconds() > 180

		if verbose == false {
			if timeout {
				continue
			}
		}

		for d_id, d := range n.Devices {
			usedMemoryPercentage := int(d.MemoryUtilization.Used * 100 / d.MemoryUtilization.Total)

			name := ""
			if d_id == 0 {
				name = n.Name
			}

			tableRow := []interface{}{
				name,
				fmt.Sprintf("%d:%s", d.Id, d.Name),
				fmt.Sprintf("%d MiB / %d MiB", d.MemoryUtilization.Used/1024/1024, d.MemoryUtilization.Total/1024/1024),
				fmt.Sprintf("%d %%", usedMemoryPercentage),
				fmt.Sprintf("%d %%", d.Utilization),
			}

			if verbose {
				lastseen := ""
				if d_id == 0 {
					lastseen = n.Time.Format("Mon Jan 2 15:04:05 2006")
				}
				tableRow = append(tableRow, lastseen)
				tableRow = append(tableRow, timeout)
			}

			table.AddRow(tableRow...)
			table.SetAlign(termtables.AlignRight, 3)

		}
		if n_id < len(c.Nodes)-1 {
			table.AddSeparator()
		}
	}
	fmt.Printf("\033[2J")
	fmt.Println(time.Now().Format("Mon Jan 2 15:04:05 2006"))
	fmt.Println(table.Render())
}
