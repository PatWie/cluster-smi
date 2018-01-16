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

type Process struct {
	Pid           int
	UsedGpuMemory int64
}

type Device struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Utilization       int    `json:"utilization"`
	MemoryUtilization Memory `json:"memory"`
	Processes         []Process
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

func (c *Cluster) Print(show_processes bool, show_time bool) {

	table := termtables.CreateTable()

	tableHeader := []interface{}{"Node", "Gpu", "Memory-Usage", "GPU-Util"}

	if show_processes {
		tableHeader = append(tableHeader, "Proccesses")
	}
	if show_time {
		tableHeader = append(tableHeader, "Last Seen")
	}
	table.AddHeaders(tableHeader...)

	now := time.Now()

	for n_id, n := range c.Nodes {

		timeout := now.Sub(n.Time).Seconds() > 180
		node_name := n.Name
		node_lastseen := n.Time.Format("Mon Jan 2 15:04:05 2006")

		if timeout {

			tableRow := []interface{}{
				node_name,
				"Offline",
				"",
				"",
			}

			if show_processes {
				tableRow = append(tableRow, "")
			}

			if show_time {
				tableRow = append(tableRow, node_lastseen)
			}

			table.AddRow(tableRow...)
			table.SetAlign(termtables.AlignRight, 3)

			if show_processes {
				table.SetAlign(termtables.AlignRight, 5)
			}

		} else {
			for d_id, d := range n.Devices {

				device_name := fmt.Sprintf("%d:%s", d.Id, d.Name)
				device_MemoryInfo := fmt.Sprintf("%d MiB / %d MiB (%d %%)",
					d.MemoryUtilization.Used/1024/1024,
					d.MemoryUtilization.Total/1024/1024,
					int(d.MemoryUtilization.Used*100/d.MemoryUtilization.Total))
				device_utilization := fmt.Sprintf("%d %%", d.Utilization)

				if timeout {
					device_MemoryInfo = "TimeOut"
					device_utilization = "TimeOut"
				}

				if d_id > 0 {
					node_name = ""
				}
				if d_id > 0 || !show_time {
					node_lastseen = ""
				}

				if len(d.Processes) > 0 && show_processes {
					for p_id, p := range d.Processes {

						if p_id > 0 {
							node_name = ""
							device_name = ""
							device_MemoryInfo = ""
							device_utilization = ""
						}

						tableRow := []interface{}{
							node_name,
							device_name,
							device_MemoryInfo,
							device_utilization,
							fmt.Sprintf("(%d) %3d MiB", p.Pid, p.UsedGpuMemory/1024/1024),
						}

						if show_time {
							tableRow = append(tableRow, node_lastseen)
						}

						table.AddRow(tableRow...)
						table.SetAlign(termtables.AlignRight, 3)
						if show_processes {
							table.SetAlign(termtables.AlignRight, 5)
						}
					}
				} else {

					tableRow := []interface{}{
						node_name,
						device_name,
						device_MemoryInfo,
						device_utilization,
					}

					if show_processes {
						tableRow = append(tableRow, "")
					}

					if show_time {
						tableRow = append(tableRow, node_lastseen)
					}

					table.AddRow(tableRow...)
					table.SetAlign(termtables.AlignRight, 3)
					if show_processes {
						table.SetAlign(termtables.AlignRight, 5)
					}

				}

			}
		}

		if n_id < len(c.Nodes)-1 {
			table.AddSeparator()
		}
	}
	fmt.Printf("\033[2J")
	fmt.Println(time.Now().Format("Mon Jan 2 15:04:05 2006"))
	fmt.Println(table.Render())
}
