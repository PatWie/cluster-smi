go build cluster-smi.go cluster.go const.go data.go
go build cluster-smi-server.go const.go data.go
go build cluster-smi-node.go cluster.go const.go data.go

# PKG_CONFIG_PATH=/graphics/opt/opt_Ubuntu16.04/libzmq/dist/lib/pkgconfig \
# go build -v --ldflags '-extldflags "-static"' -a cluster-smi-node.go