PACKAGER="github.com/patwie/cluster-smi/compiletimeconst"

# add compile-time constants to avoid path issues to configuration files
include cluster-smi.env
LDFLAGS="-X ${PACKAGER}.TimeoutThreshold=${cluster_smi_timeout_threshold_sec} -X ${PACKAGER}.ServerIp=${cluster_smi_server_ip} -X ${PACKAGER}.PortGather=${cluster_smi_server_port_gather} -X ${PACKAGER}.PortDistribute=${cluster_smi_server_port_distribute} -X ${PACKAGER}.Tick=${cluster_smi_tick_ms}"

all:
		go build -ldflags ${LDFLAGS} cluster-smi.go config.go cluster.go
		go build -ldflags ${LDFLAGS} cluster-smi-server.go config.go
		go build -ldflags ${LDFLAGS} cluster-smi-node.go config.go cluster.go
		go build -ldflags ${LDFLAGS} cluster-smi-local.go config.go cluster.go

# PKG_CONFIG_PATH=/graphics/opt/opt_Ubuntu16.04/libzmq/dist/lib/pkgconfig \
# go build -v --ldflags '-extldflags "-static"' -a cluster-smi-node.go