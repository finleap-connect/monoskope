FROM gitlab.figo.systems/platform/dependency_proxy/containers/docker:18-git

ARG GOLANG_VERSION=1.15.3

RUN apk add --no-cache make
RUN wget https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.15.3.linux-amd64.tar.gz && \
    export PATH=$PATH:/usr/local/go/bin