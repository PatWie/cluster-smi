module github.com/patwie/cluster-smi

go 1.13

require (
	github.com/brettski/go-termtables v0.0.0-00010101000000-000000000000
	github.com/c9s/goprocinfo v0.0.0-20210130143923-c95fcf8c64a8
	github.com/pebbe/zmq4 v1.2.5
	github.com/vmihailenco/msgpack/v5 v5.2.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.0.0
)

replace github.com/brettski/go-termtables => github.com/yan12125/go-termtables v0.0.0-20210222125219-99eaf4bd18fd

replace github.com/pebbe/zmq4 => github.com/yan12125/zmq4 v1.2.6-0.20210222135945-91b7a5fdef51
