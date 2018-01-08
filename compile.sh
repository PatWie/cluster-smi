go build cluster-smi.go cluster.go config.go data.go
go build cluster-smi-server.go config.go data.go
go build cluster-smi-node.go cluster.go config.go data.go

# PKG_CONFIG_PATH=/graphics/opt/opt_Ubuntu16.04/libzmq/dist/lib/pkgconfig \
# go build -v --ldflags '-extldflags "-static"' -a cluster-smi-node.go