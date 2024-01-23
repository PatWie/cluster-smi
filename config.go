package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Tick     int    `yaml:"tick"`      // tick (in seconds) between receiving data
	Timeout  int    `yaml:"timeout"`   // threshold (in seconds) after a node is considered/displayed as offline
	RouterIp string `yaml:"router_ip"` // ip of cluster-smi-router
	Ports    struct {
		Nodes   string `yaml:"nodes"`   // port of cluster-smi-router, which nodes send to
		Clients string `yaml:"clients"` // port of cluster-smi-router, where clients subscribe to
	} `yaml:"ports"`
}

func ConfigFilePath() string {
	configFilePath := "/usr/local/etc/cluster-smi.yml"
	customPath := false
	if os.Getenv("CLUSTER_SMI_CONFIG_PATH") != "" {
		configFilePath = os.Getenv("CLUSTER_SMI_CONFIG_PATH")
	}
	_, err := os.Stat(configFilePath)
	if err != nil {
		if customPath {
			log.Println("Config file %s not accessible", configFilePath)
		}
		configFilePath = ""
	}
	return configFilePath
}

func LoadConfig() Config {

	c := Config{}

	c.RouterIp = "127.0.0.1"
	c.Tick = 3
	c.Timeout = 180
	c.Ports.Nodes = "9080"
	c.Ports.Clients = "9081"

	fn := ConfigFilePath()
	if fn != "" {
		yamlFile, err := ioutil.ReadFile(fn)
		if err != nil {
			log.Fatalf("Error: %v ", err)
		}

		err = yaml.Unmarshal(yamlFile, &c)
		if err != nil {
			log.Fatalf("Error in %s %v", fn, err)
		}
	}

	return c
}

func (c Config) Print() {

	fn := ConfigFilePath()
	if fn != "" {
		log.Println("Read configuration from", fn)
	} else {
		log.Println("use default configuration")
	}

	log.Println("  Tick:", c.Tick)
	log.Println("  Timeout:", c.Timeout)
	log.Println("  RouterIp:", c.RouterIp)
	log.Println("  Ports:")
	log.Println("    Nodes:", c.Ports.Nodes)
	log.Println("    Clients:", c.Ports.Clients)

}
