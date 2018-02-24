all:
		go build cluster-smi.go config.go
		go build cluster-smi-router.go config.go
		go build cluster-smi-node.go config.go cluster.go
		go build cluster-smi-local.go config.go cluster.go