FROM registry.gitlab.figo.systems/finleap-cloud-tools/golang-builder:1.16-alpine3.13 AS builder

WORKDIR /tmp/build

# Install SSL ca certificates.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

FROM scratch

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ADD /clboreactor /

CMD ["/clboreactor", "serve"]
