FROM nvidia/cuda:9.0-devel-ubuntu16.04


RUN apt-get update && apt-get install -yy \
    wget \
    libtool-bin \
    autoconf \
    g++ \
    git \
    make \
    golang-go


RUN mkdir /zmq
RUN wget http://files.patwie.com/mirror/zeromq-4.1.0-rc1.tar.gz
RUN tar -xf zeromq-4.1.0-rc1.tar.gz -C /zmq
WORKDIR /zmq/zeromq-4.1.0
RUN ./autogen.sh
RUN ./configure
RUN ./configure --prefix=/zmq/zeromq-4.1.0/dist
RUN make
RUN make install


ENV PKG_CONFIG_PATH="/zmq/zeromq-4.1.0/dist/lib/pkgconfig:${PKG_CONFIG_PATH}"
ENV LD_LIBRARY_PATH="/zmq/zeromq-4.1.0/dist/lib:${LD_LIBRARY_PATH}"
ENV GOPATH="/gocode"
RUN mkdir -p /gocode/src/github.com/patwie/cluster-smi
ENV CLUSTER_SMI_CONFIG_PATH="/cluster-smi.yml"


