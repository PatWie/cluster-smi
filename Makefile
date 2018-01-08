PACKAGER="github.com/patwie/cluster-smi/compiletimeconst"

# add compile-time constants to avoid path issues to configuration files
include cluster-smi.env
LDFLAGS="-X ${PACKAGER}.ServerIp=${cluster_smi_server_ip} -X ${PACKAGER}.PortGather=${cluster_smi_server_port_gather} -X ${PACKAGER}.PortDistribute=${cluster_smi_server_port_distribute} -X ${PACKAGER}.Tick=${cluster_smi_tick_ms}"

all:
		go build -ldflags ${LDFLAGS} cluster-smi.go cluster.go config.go data.go
		go build -ldflags ${LDFLAGS} cluster-smi-server.go config.go data.go
		go build -ldflags ${LDFLAGS} cluster-smi-node.go cluster.go config.go data.go

# PKG_CONFIG_PATH=/graphics/opt/opt_Ubuntu16.04/libzmq/dist/lib/pkgconfig \
# go build -v --ldflags '-extldflags "-static"' -a cluster-smi-node.go