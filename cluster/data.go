package cluster

import (
	"fmt"
	"github.com/brettski/go-termtables"
	"os/user"
	"regexp"
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
	Pid             int
	UsedGpuMemory   int64
	Name            string
	Username        string
	RunTime         int64
	ExtendedCommand string
}

type Device struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Utilization       int    `json:"utilization"`
	MemoryUtilization Memory `json:"memory"`
	FanSpeed          int    `json:"fan_speed"`
 	Temperature       int    `json:"temperature"`
 	PowerUsage        int    `json:"power_usage"`
	Processes         []Process
}

type Node struct {
	Name       string    `json:"name"`       // hostname
	Devices    []Device  `json:"devices"`    // devices
	Time       time.Time `json:"time"`       // current timestamp from message
	BootTime   int64     `json:"boot_time"`  // uptime of system
	ClockTicks int64     `json:"clock_ticks` // cpu ticks per second
	CPUUtil    float32   `json:"cpu_util"`   // cpu utilization
	LastRun    uint64    `json:"last_run"`
	LastIdle   uint64    `json:"last_idle"`
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

func FilterByUser(c Cluster, username string) Cluster {

	var clus Cluster

	for _, n := range c.Nodes {
		var Devices []Device

		for _, d := range n.Devices {
			var Processes []Process
			for _, p := range d.Processes {
				if p.Username == username {
					Processes = append(Processes, p)
				}
			}

			if len(Processes) > 0 {
				current_device := Device{
					d.Id, d.Name, d.Utilization, d.MemoryUtilization,
 					d.FanSpeed, d.Temperature, d.PowerUsage, Processes,
				}
				Devices = append(Devices, current_device)
			}
		}

		if len(Devices) > 0 {
			current_node := Node{
				n.Name, Devices, n.Time, n.BootTime, n.ClockTicks, n.CPUUtil, n.LastRun, n.LastIdle,
			}
			clus.Nodes = append(clus.Nodes, current_node)
		}

	}

	return clus
}

func (c *Cluster) Sort() {
	sort.Sort(ByName(c.Nodes))
}

func HumanizeSeconds(secs int64) string {

	days := secs / 60 / 60 / 24
	hours := (secs / 60 / 60) % 24
	minutes := (secs / 60) % 60
	seconds := secs % 60

	// a bug in term-tables ? cannot "right-align" last cell to
	has_prefix := false

	answer := ""
	if days > 0 {
		answer = fmt.Sprintf("%s %2d d", answer, days)
		has_prefix = true
	} else {
		answer = fmt.Sprintf("%s     ", answer)
	}
	if hours > 0 {
		answer = fmt.Sprintf("%s %2d h", answer, hours)
		has_prefix = true
	} else {
		if has_prefix {
			answer = fmt.Sprintf("%s  0 h", answer)
		} else {
			answer = fmt.Sprintf("%s     ", answer)
		}

	}
	if minutes > 0 {
		answer = fmt.Sprintf("%s %2d min", answer, minutes)
		has_prefix = true
	} else {
		if has_prefix {
			answer = fmt.Sprintf("%s  0 min", answer)
		} else {
			answer = fmt.Sprintf("%s     ", answer)
		}

	}
	if seconds > 0 {
		if has_prefix {
			answer = fmt.Sprintf("%s %2d sec", answer, seconds)
		} else {
			answer = fmt.Sprintf("%s   %2d sec", answer, seconds)
		}
		has_prefix = true
	} else {
		if has_prefix {
			answer = fmt.Sprintf("%s  0 sec", answer)
		} else {
			answer = fmt.Sprintf("%s     ", answer)
		}
	}

	return answer

}

func (c *Cluster) FilterNodes(node_regex string) {
	r, _ := regexp.Compile(node_regex)
	var match_nodes []Node

	for _, node := range c.Nodes {
		if r.MatchString(node.Name) {
			match_nodes = append(match_nodes, node)
		}
	}

	c.Nodes = match_nodes
}

func highlight(s string, useHTML bool) string {
	if useHTML {
		return fmt.Sprintf("<span style=\"color: yellow\">%s</span>", s)
	} else {
		return fmt.Sprintf("\033[0;33m%s\033[0m", s)
	}
}

func (c *Cluster) Print(show_processes bool, show_time bool, timeout_threshold int, useColor bool, useHTML bool, extended bool, show_detail bool) {

	table := termtables.CreateTable()

	if useHTML {
		table.SetModeHTML()
	}

	tableHeader := []interface{}{"Node", "Gpu", "Memory-Usage", "GPU-Util", "CPU-Util"}

	if show_detail {
		tableHeader = append(tableHeader, "Fan")
		tableHeader = append(tableHeader, "Temp")
		tableHeader = append(tableHeader, "Power")
	}
	if show_processes {
		tableHeader = append(tableHeader, "PID")
		tableHeader = append(tableHeader, "User")
		tableHeader = append(tableHeader, "Command")
		tableHeader = append(tableHeader, "GPU Mem")
		tableHeader = append(tableHeader, "Runtime")
	}
	if show_time {
		tableHeader = append(tableHeader, "Last Seen")
	}
	table.AddHeaders(tableHeader...)

	now := time.Now()

	currentUser, _ := user.Current()

	for n_id, n := range c.Nodes {

		timeout := now.Sub(n.Time).Seconds() > float64(timeout_threshold)
		node_name := n.Name
		node_lastseen := n.Time.Format("Mon Jan 2 15:04:05 2006")

		if timeout {

			tableRow := []interface{}{
				node_name,
				"Offline",
				"",
				"",
				"",
			}

			if show_detail {
				tableRow = append(tableRow, []interface{}{"", "", ""}...)
			}

			if show_processes {
				tableRow = append(tableRow, []interface{}{"", "", "", "", ""}...)
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
				device_MemoryInfo := fmt.Sprintf("%d MiB / %d MiB (%3d %%)",
					d.MemoryUtilization.Used/1024/1024,
					d.MemoryUtilization.Total/1024/1024,
					int(d.MemoryUtilization.Used*100/d.MemoryUtilization.Total))
				device_utilization := fmt.Sprintf("%3d %%", d.Utilization)
				cpu_utilization := fmt.Sprintf("%.2f %%", n.CPUUtil * 100)

				device_FanSpeed := ""
				device_Temp := ""
				device_PowerUtil := ""
				if show_detail {
					device_FanSpeed = fmt.Sprintf("%3d %%", d.FanSpeed)
					device_Temp = fmt.Sprintf("%3d °C", d.Temperature)
					device_PowerUtil = fmt.Sprintf("%3dW", d.PowerUsage)
				}

				if timeout {
					device_MemoryInfo = "TimeOut"
					device_utilization = "TimeOut"
				}

				if d_id > 0 {
					node_name = ""
					cpu_utilization = ""
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
							if show_detail {
								device_FanSpeed = ""
								device_Temp = ""
								device_PowerUtil = ""
							}
							cpu_utilization = ""
						}

						processName := p.Name
						if extended {
							processName = fmt.Sprintf("%.55s", p.ExtendedCommand)
							// cmdName =
                                                }
						processUseGPUMemory := fmt.Sprintf("%3d MiB", p.UsedGpuMemory/1024/1024)
						processRuntime := HumanizeSeconds(p.RunTime)
						processPID := fmt.Sprintf("%v", p.Pid)
						processUsername := p.Username
						if p.Username == currentUser.Username {
							node_name = highlight(node_name, useHTML)
							device_name = highlight(device_name, useHTML)
							device_MemoryInfo = highlight(device_MemoryInfo, useHTML)
							device_utilization = highlight(device_utilization, useHTML)
							if show_detail {
							   device_FanSpeed = highlight(device_FanSpeed, useHTML)
							   device_Temp = highlight(device_Temp, useHTML)
							   device_PowerUtil = highlight(device_PowerUtil, useHTML)
							}
							processPID = highlight(fmt.Sprintf("%v", p.Pid), useHTML)
							processUsername = highlight(processUsername, useHTML)
							processName = highlight(processName, useHTML)
							processUseGPUMemory = highlight(processUseGPUMemory, useHTML)
							processRuntime = highlight(processRuntime, useHTML)
						}

						if d.Utilization < 10 {
							if useHTML {
								device_utilization = fmt.Sprintf("<span style=\"color: red\">%s</span>", device_utilization)
							} else {
								device_utilization = fmt.Sprintf("\033[0;31m%s\033[0m", device_utilization)
							}
						}

						tableRow := []interface{}{
							node_name,
							device_name,
							device_MemoryInfo,
							device_utilization,
							cpu_utilization,
						}
						if show_detail {
							tableRow = append(tableRow, []interface{}{
								device_FanSpeed,
								device_Temp,
								device_PowerUtil,
								}...)
						}
						tableRow = append(tableRow, []interface{}{
								processPID,
								processUsername,
								processName,
								processUseGPUMemory,
								processRuntime,
							}...)

						// fmt.Sprintf("%s (%d, %s) %3d MiB %v", p.Name, p.Pid, p.Username, p.UsedGpuMemory/1024/1024, p.RunTime),

						if show_time {
							if p_id > 0 {
								tableRow = append(tableRow, "")

							} else {
								tableRow = append(tableRow, node_lastseen)
								//FIXME
							}
						}

						table.AddRow(tableRow...)
						table.SetAlign(termtables.AlignRight, 3)
						table.SetAlign(termtables.AlignCenter, 4)
						if show_processes && show_detail {
							table.SetAlign(termtables.AlignRight, 5)
							table.SetAlign(termtables.AlignRight, 7)
							table.SetAlign(termtables.AlignRight, 8)
							table.SetAlign(termtables.AlignRight, 9)
							table.SetAlign(termtables.AlignRight, 10)
							table.SetAlign(termtables.AlignRight, 11)
							table.SetAlign(termtables.AlignRight, 12)
							table.SetAlign(termtables.AlignRight, 13)
						}
					}

				} else {
					if len(d.Processes) == 0 && useColor {
						if useHTML {
							device_name = fmt.Sprintf("<span style=\"color: green\">%s</span>", device_name)
						} else {
							device_name = fmt.Sprintf("\033[0;32m%s\033[0m", device_name)
						}
					}

					tableRow := []interface{}{
						node_name,
						device_name,
						device_MemoryInfo,
						device_utilization,
						cpu_utilization,
					}

					if show_detail {
						tableRow = append(tableRow, []interface{}{
							device_FanSpeed,
							device_Temp,
							device_PowerUtil,
							}...)
					}

					if show_processes {
						tableRow = append(tableRow, []interface{}{"", "", "", "", ""}...)
					}

					if show_time {
						tableRow = append(tableRow, node_lastseen)
					}

					table.AddRow(tableRow...)
					table.SetAlign(termtables.AlignRight, 3)
					if show_processes {
						table.SetAlign(termtables.AlignRight, 5)
						table.SetAlign(termtables.AlignRight, 8)
						table.SetAlign(termtables.AlignRight, 9)
					}

				}

			}
		}

		if n_id < len(c.Nodes)-1 {
			table.AddSeparator()
		}
	}
	if !useHTML {
		fmt.Printf("\033[2J")
	}
	// fmt.Printf("\033[0;30m color here \033[0m") // Black - Regular
	// fmt.Printf("\033[0;31m color here \033[0m") // Red
	// fmt.Printf("\033[0;32m color here \033[0m") // Green
	// fmt.Printf("\033[0;33m color here \033[0m") // Yellow
	// fmt.Printf("\033[0;34m color here \033[0m") // Blue
	// fmt.Printf("\033[0;35m color here \033[0m") // Purple
	// fmt.Printf("\033[0;36m color here \033[0m") // Cyan
	// fmt.Printf("\033[0;37m color here \033[0m") // White
	// fmt.Printf("\033[1;30m color here \033[0m") // Black - Bold
	// fmt.Printf("\033[1;31m color here \033[0m") // Red
	// fmt.Printf("\033[1;32m color here \033[0m") // Green
	// fmt.Printf("\033[1;33m color here \033[0m") // Yellow
	// fmt.Printf("\033[1;34m color here \033[0m") // Blue
	// fmt.Printf("\033[1;35m color here \033[0m") // Purple
	// fmt.Printf("\033[1;36m color here \033[0m") // Cyan
	// fmt.Printf("\033[1;37m color here \033[0m") // White
	// fmt.Printf("\033[4;30m color here \033[0m") // Black - Underline
	// fmt.Printf("\033[4;31m color here \033[0m") // Red
	// fmt.Printf("\033[4;32m color here \033[0m") // Green
	// fmt.Printf("\033[4;33m color here \033[0m") // Yellow
	// fmt.Printf("\033[4;34m color here \033[0m") // Blue
	// fmt.Printf("\033[4;35m color here \033[0m") // Purple
	// fmt.Printf("\033[4;36m color here \033[0m") // Cyan
	// fmt.Printf("\033[4;37m color here \033[0m") // White
	// fmt.Printf("\033[0m color here \033[0m")
	fmt.Println(time.Now().Format("Mon Jan 2 15:04:05 2006") + " (http://github.com/patwie/cluster-smi)")
	fmt.Println(table.Render())
}
