FROM registry.gitlab.figo.systems/platform/golang-builder:go1.15-alpine3.12 AS builder

WORKDIR /tmp/build

ENV GRPC_HEALTH_PROBE_VERSION=v0.3.5
RUN wget -qOgrpc-health-probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x grpc-health-probe

# Install SSL ca certificates.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

FROM gitlab.figo.systems/platform/dependency_proxy/containers/scratch

# Import from builder.
COPY --from=builder /tmp/build/grpc-health-probe /bin/grpc-health-probe
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ADD /eventstore /

CMD ["/eventstore", "server"]
