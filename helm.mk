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
	@$(HELM) upgrade --install $* $(HELM_REGISTRY_ALIAS)/$* --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE) --version $(VERSION)

uninstall-%: 
	@$(HELM) uninstall $* --namespace $(KUBE_NAMESPACE)

template-%: clean 
	@mkdir -p $(HELM_OUTPUT_DIR)
	@$(HELM) template $* $(HELM_PATH)/$* --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE) --output-dir $(HELM_OUTPUT_DIR) --include-crds --debug

add-kubism:
	@$(HELM) repo add kubism.io https://kubism.github.io/charts/
	@$(HELM) repo update

package-%:
	@cp "$(HELM_PATH)/$*/values.yaml" "$(HELM_PATH)/$*/values.yaml.bkp"
	@yq write "$(HELM_PATH)/$*/values.yaml" image.tag "$(VERSION)" --inplace
	@$(HELM) package $(HELM_PATH)/$* --dependency-update --version $(VERSION)
	@mv "$(HELM_PATH)/$*/values.yaml.bkp" "$(HELM_PATH)/$*/values.yaml"

push-%:
	@curl --fail -u $(HELM_USER):$(HELM_PASSWORD) -T $*-$(VERSION).tgz "$(HELM_REGISTRY)/$*-$(VERSION).tgz"