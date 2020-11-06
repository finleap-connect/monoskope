HELM                		?= helm3
HELM_OUTPUT_DIR             ?= tmp

.PHONY: helm-template-clean helm-dependency-update helm-install helm-uninstall helm-template

clean:
	@rm -Rf $(HELM_OUTPUT_DIR)

dep-%:
	@$(HELM) dep update $(HELM_PATH)/$*

lint-%:
	@$(HELM) lint $(HELM_PATH)/$*

install-%: 
	@$(HELM) upgrade --install $* $(HELM_PATH)/$* --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE)

uninstall-%: 
	@$(HELM) uninstall $* --namespace $(KUBE_NAMESPACE)

template-%: clean 
	@mkdir -p $(HELM_OUTPUT_DIR)
	@$(HELM) template $* $(HELM_PATH)/$* --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE) --output-dir $(HELM_OUTPUT_DIR) --include-crds --debug

add-kubism:
	@$(HELM) repo add kubism.io https://kubism.github.io/charts/
	@$(HELM) repo update