package main

import (
	"strings"
)

type ByName []Node

func (a ByName) Len() int      { return len(a) }
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool {
	// return len(a[i].Name) < len(a[j].Name)
	return strings.Compare(a[i].Name, a[j].Name) < 0
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
	MemoryUtilization Memory
}

type Node struct {
	Name    string   `json:"name"`
	Devices []Device `json:"devices"`
}

type Cluster struct {
	Nodes []Node `json:"nodes"`
}
