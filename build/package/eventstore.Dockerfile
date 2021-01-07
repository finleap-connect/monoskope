FROM registry.gitlab.figo.systems/platform/monoskope/monoskope/builder:latest AS builder

WORKDIR /tmp/build

# Install SSL ca certificates.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

FROM scratch

# Import from builder.
COPY --from=builder /tmp/build/grpc-health-probe /bin/grpc-health-probe
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ADD /eventstore /

CMD ["/eventstore", "server"]
