FROM gitlab.figo.systems/platform/dependency_proxy/containers/golang:1.15.3-buster

# Install Docker
RUN apt-get update && apt install docker.io -y \
    && rm -rf /var/lib/apt/lists/*
