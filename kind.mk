.PHONY: template-clean dependency-update install uninstall template docs

kind-create: ## create kind cluster
	@kind create cluster --name m8kind --config build/deploy/kind/kubeadm_conf.yaml --kubeconfig tmp/kind-kubeconfig