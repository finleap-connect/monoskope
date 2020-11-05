FROM registry.gitlab.figo.systems/platform/golang-builder:go1.15-alpine3.12 AS builder

# Install SSL ca certificates.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

FROM scratch

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ADD /gateway /

CMD ["/gateway", "server"]
