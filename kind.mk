KIND_VERSION ?= v1.15.12
KIND_KUBECONFIG=tmp/kind-kubeconfig
CERT_MANAGER_VERSION ?= v1.1.0
K8S_CLUSTER_NAME := m8-dev-cluster
HELM_KIND_VALUES_FILE ?= examples/01-monoskope-kind-values.yaml

##@ Kind

kind-create-cluster: ## create kind cluster
	@kind get clusters | grep ${K8S_CLUSTER_NAME} || kind create cluster --name ${K8S_CLUSTER_NAME} --config build/deploy/kind/kind_config_${KIND_VERSION}.yaml --kubeconfig ${KIND_KUBECONFIG}
	@kind get kubeconfig --name ${K8S_CLUSTER_NAME} > ${KIND_KUBECONFIG}

kind-delete-cluster: ## destroy kind cluster
	@kind delete cluster --name ${K8S_CLUSTER_NAME}

kind-helm-repos: ## add & update helm repos necessary
	@$(HELM) repo add $(HELM_REGISTRY_ALIAS) "$(HELM_REGISTRY)"
	@$(HELM) repo add jetstack https://charts.jetstack.io
	@$(HELM) repo add dex https://charts.dexidp.io
	@$(HELM) repo update

kind-trust-anchor: ## create trust-anchor in kind cluster
	@echo "Generating trust-anchor for m8 PKI..."
	@step certificate create root.monoskope.cluster.local tmp/ca.crt tmp/ca.key --profile root-ca --no-password --insecure --not-after=87600h
	@echo "Creating secret containing trust-anchor in kind cluster..."
	@kubectl --kubeconfig ${KIND_KUBECONFIG} create namespace monoskope --dry-run -o yaml | kubectl --kubeconfig ${KIND_KUBECONFIG} apply -f -
	@kubectl --kubeconfig ${KIND_KUBECONFIG} -n monoskope create secret tls m8-trust-anchor --cert=tmp/ca.crt --key=tmp/ca.key --dry-run -o yaml | kubectl --kubeconfig ${KIND_KUBECONFIG} apply -f -

kind-install-certmanager:
	@echo "Installing cert-manager into kind cluster..."
	@$(HELM) upgrade --install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --version ${CERT_MANAGER_VERSION} --values examples/02-kind-cert-manager-values.yaml
	@kubectl --kubeconfig ${KIND_KUBECONFIG} apply -f examples/04-kind-cert-manager-issuer.yaml
	
kind-install-dex:
	@echo "Installing dex into kind cluster..."
	@$(HELM) --kubeconfig ${KIND_KUBECONFIG} upgrade --install dex --wait dex/dex --values examples/03-kind-dex-values.yaml

kind-helm-clean: ## clean up templated helm charts
	@rm -Rf $(HELM_OUTPUT_DIR)

kind-helm-template-monoskope: kind-helm-clean
	@$(HELM) --kubeconfig ${KIND_KUBECONFIG} template $(HELM_RELEASE) $(HELM_REGISTRY_ALIAS)/monoskope --version $(LATEST_TAG) --namespace monoskope --values $(HELM_KIND_VALUES_FILE) --output-dir $(HELM_OUTPUT_DIR)

kind-install-monoskope: ## installs monoskope into kind cluster using the latest tag available
	@$(HELM) --kubeconfig ${KIND_KUBECONFIG} upgrade --install $(HELM_RELEASE) $(HELM_REGISTRY_ALIAS)/monoskope --namespace monoskope --create-namespace --version $(LATEST_TAG) --values $(HELM_KIND_VALUES_FILE)

kind-setup-monoskope: kind-create-cluster kind-trust-anchor kind-helm-repos kind-install-certmanager kind-install-dex kind-install-monoskope ## install monoskope with kind