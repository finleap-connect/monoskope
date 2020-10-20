KIND         ?= $(TOOLS_DIR)/kind
KIND_CLUSTER ?= monoskope-test
KIND_CONFIG  ?= config/kind/kindconf.yaml
KIND_VERSION ?= v0.7.0

get:
	$(shell $(TOOLS_DIR)/goget-wrapper sigs.k8s.io/kind@$(KIND_VERSION))

clean:
	rm -Rf $(TOOLS_DIR)/kind

create: get
	$(KIND) create cluster --config $(KIND_CONFIG) --name $(KIND_CLUSTER) --kubeconfig /tmp/kind-$(KIND_CLUSTER)-config --wait 5m
	@echo "kubectl --kubeconfig \"/tmp/kind-$(KIND_CLUSTER)-config\" get no"

is-running: get
	@echo "Checking if kind cluster with name '$(KIND_CLUSTER)' is running..."
	@echo "(e.g. create cluster via 'make kind-create')"
	@{ \
	set -e; \
	$(KIND) get kubeconfig --name $(KIND_CLUSTER) > /dev/null; \
	}

get-kubeconfig: get
	$(KIND) get kubeconfig --name $(KIND_CLUSTER) > /tmp/kind-$(KIND_CLUSTER)-config
	@echo "Created untracked config file in '/tmp/kind-$(KIND_CLUSTER)-config. Use as follows:"
	@echo "kubectl --kubeconfig \"/tmp/kind-$(KIND_CLUSTER)-config\" get no"

delete: get
	$(KIND) delete cluster --name $(KIND_CLUSTER)