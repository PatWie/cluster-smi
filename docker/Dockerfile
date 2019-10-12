FROM centos:7

ARG ZEROMQ_VER=4.3.2
ARG GO_VER=1.13
ARG CUDA_VER=8.0.61

ADD https://github.com/zeromq/libzmq/releases/download/v${ZEROMQ_VER}/zeromq-${ZEROMQ_VER}.tar.gz \
    https://dl.google.com/go/go${GO_VER}.linux-amd64.tar.gz \
    https://developer.download.nvidia.com/compute/cuda/repos/rhel7/x86_64/cuda-repo-rhel7-${CUDA_VER}-1.x86_64.rpm /

RUN yum -y update && yum -y install gcc-c++ make && \
    tar xf go${GO_VER}.linux-amd64.tar.gz && \
    rpm -i cuda-repo-rhel7-${CUDA_VER}-1.x86_64.rpm && \
    bash -c 'TMP=${CUDA_VER%.*} ; yum -y install cuda-nvml-dev-${TMP/./-}'

# Need to patch libzmq.pc.in due to https://github.com/zeromq/libzmq/issues/3710
RUN tar xf zeromq-${ZEROMQ_VER}.tar.gz && \
    pushd /zeromq-${ZEROMQ_VER} && \
    sed -i 's#@pkg_config_libs_private@#-lm @pkg_config_libs_private@#' src/libzmq.pc.in && \
    ./configure \
        --disable-shared \
        --enable-static \
        --without-libsodium \
        --disable-libunwind && \
    make && make install && \
    popd && rm -rvf /zeromq-${ZEROMQ_VER}

ENV PATH="/go/bin:${PATH}" GOPATH=/gopath PKG_CONFIG_PATH=/usr/local/lib/pkgconfig

WORKDIR /gopath/src/github.com/patwie/cluster-smi

CMD ["bash"]
