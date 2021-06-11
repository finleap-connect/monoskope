FROM gitlab.figo.systems/platform/dependency_proxy/containers/golang:1.16-buster

# ensure versions are synched with the Makefile!
ARG PROTOC_IMPORTS_DIR=/include
ARG PROTOC_VERSION=3.17.0
ARG PROTOC_GEN_GO_VERSION=v1.26.0
ARG PROTOC_GEN_GO_GRPC_VERSION=v1.1.0
ARG ARCH=linux-x86_64

ENV TOOLS_DIR  /tools
ENV GINKGO     ginkgo
ENV LINTER     golangci-lint
ENV MOCKGEN    mockgen
ENV PROTOC     /tools/protoc

WORKDIR /tmp/build

# Install Docker
# RUN apt-get update && apt install unzip docker.io -y \
#    && rm -rf /var/lib/apt/lists/*

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.39.0

RUN go get -u github.com/onsi/ginkgo/ginkgo

RUN go get -u github.com/golang/mock/mockgen

RUN curl -LO "https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/protoc-$PROTOC_VERSION-$ARCH.zip" ; \
    echo curl -LO "https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/protoc-$PROTOC_VERSION-$ARCH.zip" ; \
    ls -la ; \
    unzip protoc-$PROTOC_VERSION-$ARCH.zip -d /.protoc-unpack ;\
    mkdir -p $PROTOC_IMPORTS_DIR/ $TOOLS_DIR ;\
    mv /.protoc-unpack/bin/protoc $PROTOC ;\
    cp -a /.protoc-unpack/include/* $PROTOC_IMPORTS_DIR/ ;\
    rm -rf /.protoc-unpack/ protoc-$PROTOC_VERSION-$ARCH.zip ;\
    go get -u google.golang.org/protobuf/cmd/protoc-gen-go@$PROTOC_GEN_GO_VERSION ;\
    go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc@$PROTOC_GEN_GO_GRPC_VERSION