HELM                		?= helm3
HELM_OUTPUT_DIR             ?= tmp
HELM_REGISTRY 				?= https://artifactory.figo.systems/artifactory/virtual_helm
HELM_REGISTRY_ALIAS			?= finleap


.PHONY: helm-template-clean helm-dependency-update helm-install helm-uninstall helm-template

clean:
	@rm -Rf $(HELM_OUTPUT_DIR)

dep-%:
	@$(HELM) dep update $(HELM_PATH)/$*

lint-%:
	@$(HELM) lint $(HELM_PATH)/$*

install-%: 
	@$(HELM) upgrade --install $* $(HELM_PATH)/$* --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE)

install-from-repo-%:
	@$(HELM) repo update
	@$(HELM) upgrade --install $* $(HELM_REGISTRY_ALIAS)/$* --namespace $(KUBE_NAMESPACE) --version $(VERSION) --values $(HELM_VALUES_FILE)

uninstall-%: 
	@$(HELM) uninstall $* --namespace $(KUBE_NAMESPACE)

template-%: clean 
	@mkdir -p $(HELM_OUTPUT_DIR)
	@$(HELM) template $* $(HELM_PATH)/$* --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE) --output-dir $(HELM_OUTPUT_DIR) --include-crds --debug

add-kubism:
	@$(HELM) repo add kubism.io https://kubism.github.io/charts/

add-finleap:
	@$(HELM) repo add --username $(HELM_USER) --password $(HELM_PASSWORD) $(HELM_REGISTRY_ALIAS) "$(HELM_REGISTRY)"

update-chart-deps:
	@sed -i 's/latest/$(VERSION)/g' "$(HELM_PATH)/monoskope/Chart.yaml"
