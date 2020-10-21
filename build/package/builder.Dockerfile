FROM gitlab.figo.systems/platform/dependency_proxy/containers/docker:19

# Install Alpine Dependencies
RUN apk update && apk upgrade && apk add --no-cache bash make

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# install golang
ENV GOLANG_VERSION 1.15.3
RUN wget -O go${GO_VERSION}.amd64.tar.gz https://dl.google.com/go/go${GOLANG_VERSION}.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go${GO_VERSION}.amd64.tar.gz
ENV PATH "/usr/local/go/bin:${PATH}"