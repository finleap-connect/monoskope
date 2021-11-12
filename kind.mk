.PHONY: template-clean dependency-update install uninstall template docs

KIND_VERSION ?= v1.15.12
CERT_MANAGER_VERSION ?= v1.1.0
K8S_CLUSTER_NAME := m8-dev-cluster

HELM ?= helm
HELM_RELEASE ?= m8
HELM_REGISTRY               ?= https://finleap-connect.github.io/charts
HELM_REGISTRY_ALIAS         ?= finleap-connect
HELM_RELEASE                ?= m8
HELM_VALUES_FILE            ?= examples/01-monoskope-kind-values.yaml

export KUBECONFIG=tmp/kind-kubeconfig

##@ Kind

kind-create: ## create kind cluster
	@kind get clusters | grep ${K8S_CLUSTER_NAME} || kind create cluster --name ${K8S_CLUSTER_NAME} --config build/deploy/kind/kind_config_${KIND_VERSION}.yaml --kubeconfig ${KUBECONFIG}

kind-delete: ## destroy kind cluster
	@kind delete cluster --name ${K8S_CLUSTER_NAME}

kind-trust-anchor: ## create trust-anchor in kind cluster
	@echo "Generating trust-anchor for m8 PKI..."
	@step certificate create root.monoskope.cluster.local tmp/ca.crt tmp/ca.key --profile root-ca --no-password --insecure --not-after=87600h
	@echo "Creating secret containing trust-anchor in kind cluster..."
	@kubectl create namespace monoskope --dry-run -o yaml | kubectl apply -f -
	@kubectl -n monoskope create secret tls m8-trust-anchor --cert=tmp/ca.crt --key=tmp/ca.key --dry-run -o yaml | kubectl apply -f -

kind-install-certmanager:
	@echo "Installing cert-manager into kind cluster..."
	@$(HELM) repo add jetstack https://charts.jetstack.io
	@$(HELM) repo update
	@$(HELM) upgrade --install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --version ${CERT_MANAGER_VERSION} --values examples/02-kind-cert-manager-values.yaml

kind-install-dex:
	@echo "Installing dex into kind cluster..."
	@$(HELM) repo add dex https://charts.dexidp.io
	@$(HELM) repo update
	@$(HELM) upgrade --install dex --wait dex/dex --values examples/03-kind-dex-values.yaml

kind-install-monoskope: kind-create kind-trust-anchor kind-install-certmanager kind-install-dex ## installs monoskope into kind cluster using the latest tag available
	@$(HELM) repo add $(HELM_REGISTRY_ALIAS) "$(HELM_REGISTRY)"
	@$(HELM) repo update
	@$(HELM) upgrade --install $(HELM_RELEASE) $(HELM_REGISTRY_ALIAS)/monoskope --namespace monoskope --create-namespace --version $(LATEST_TAG) --values $(HELM_VALUES_FILE)