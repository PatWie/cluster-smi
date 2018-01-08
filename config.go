package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	_ "time"
)

type Config struct {
	ServerIp             string `yaml:"server_ip"`
	ServerPortGather     string `yaml:"server_port_gather"`
	ServerPortDistribute string `yaml:"server_port_distribute"`
	Tick                 int    `yaml:"tick_ms"`
}

func (C *Config) ReadConfig(fn string) *Config {
	yamlFile, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatalf("Error: %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, C)
	if err != nil {
		log.Fatalf("Error in %s %v", fn, err)
	}

	return C
}
