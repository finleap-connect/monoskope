FROM gitlab.figo.systems/platform/dependency_proxy/containers/golang:1.16-buster

# ensure versions are synched with the Makefile!
ENV GINKGO     ginkgo
ENV LINTER     golangci-lint
ENV MOCKGEN    mockgen

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.39.0
RUN go get -u github.com/onsi/ginkgo/ginkgo@v1.15.2
RUN go get -u github.com/golang/mock/mockgen@v1.5.0
