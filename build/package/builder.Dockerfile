FROM gitlab.figo.systems/platform/dependency_proxy/containers/docker:18-git

RUN apk add --no-cache go make
RUN go version

