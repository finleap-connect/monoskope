FROM gitlab.figo.systems/platform/dependency_proxy/containers/golang:1.15-buster

WORKDIR /tmp/build

# Install Docker
RUN apt-get update && apt install docker.io -y \
    && rm -rf /var/lib/apt/lists/*

ENV GRPC_HEALTH_PROBE_VERSION=v0.3.5
RUN wget -qOgrpc-health-probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x grpc-health-probe
