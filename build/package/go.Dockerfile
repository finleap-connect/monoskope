# Copyright 2021 Monoskope Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM registry.gitlab.figo.systems/finleap-cloud-tools/golang-builder:1.16-alpine3.13 AS builder

ARG VERSION
ARG GO_MODULE
ARG SRC
ARG COMMIT


# Install SSL ca certificates.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /workdir

ENV GOPATH /workdir/.go
ENV GRPC_HEALTH_PROBE_VERSION=v0.3.5
RUN wget -qOgrpc-health-probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x grpc-health-probe

COPY go.mod .
COPY go.sum .
COPY .go/ .go/

COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -o app -a -installsuffix cgo --tags release -ldflags "-X=${GO_MODULE}/internal/version.Version=${VERSION} -X=${GO_MODULE}/internal/version.Commit=${COMMIT}" ${SRC}

FROM scratch

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /workdir/grpc-health-probe /bin/grpc-health-probe
COPY --from=builder /workdir/app .

# Run as non root user
USER 1001:1001

CMD ["/app", "server"]
