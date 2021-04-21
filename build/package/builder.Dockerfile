FROM gitlab.figo.systems/platform/dependency_proxy/containers/golang:1.16.2-buster

WORKDIR /tmp/build

# Install Docker
RUN apt-get update && apt install docker.io -y \
    && rm -rf /var/lib/apt/lists/*

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.39.0

RUN go get -u github.com/onsi/ginkgo/ginkgo