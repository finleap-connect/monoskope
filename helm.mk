
SHELL := bash
TOOLS_DIR := $(shell cd tools && pwd)

HELM                		?= helm3
HELM_OUTPUT_DIR             ?= tmp
HELM_REGISTRY 				?= https://artifactory.figo.systems/artifactory/virtual_helm
HELM_REGISTRY_ALIAS			?= finleap

.PHONY: template-clean dependency-update install uninstall template docs

clean:
	@rm -Rf $(HELM_OUTPUT_DIR)

dep-%:
	@$(HELM) dep update $(HELM_PATH)/$*

lint-%:
	@$(HELM) lint $(HELM_PATH)/$*

install-%:
	@cat $(HELM_VALUES_FILE) | sed "s/0.0.1-local/$(VERSION)/g" > $(HELM_VALUES_FILE).tag
	@$(HELM) upgrade --install m8dev $(HELM_PATH)/$* --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE).tag --skip-crds
	@rm $(HELM_VALUES_FILE).tag

install-from-repo-%:
	@$(HELM) repo update
	@$(HELM) upgrade --install m8dev $(HELM_REGISTRY_ALIAS)/$* --namespace $(KUBE_NAMESPACE) --version $(VERSION) --values $(HELM_VALUES_FILE) --skip-crds --atomic --timeout 10m

uninstall-%: 
	@$(HELM) uninstall m8dev --namespace $(KUBE_NAMESPACE)

template-%: clean
	@mkdir -p $(HELM_OUTPUT_DIR)
	@$(HELM) template m8dev $(HELM_PATH)/$* --namespace $(KUBE_NAMESPACE) --values $(HELM_VALUES_FILE) --output-dir $(HELM_OUTPUT_DIR) --include-crds --debug
	@echo "ATTENTION:"
	@echo "If you want to have the latest dependencies (e.g. gateway chart changes)"
	@echo "execute the following command prior to the current command:"
	@echo "$$ $(MAKE) helm-dep-$*"
	@echo

add-finleap:
	@$(HELM) repo add --username $(HELM_USER) --password $(HELM_PASSWORD) $(HELM_REGISTRY_ALIAS) "$(HELM_REGISTRY)"

set-chart-version-%:
	@yq write $(HELM_PATH)/$*/Chart.yaml version "$(VERSION)" --inplace

set-app-version-%:
	@yq write $(HELM_PATH)/$*/Chart.yaml appVersion "$(VERSION)" --inplace
	@yq write $(HELM_PATH)/$*/values.yaml image.tag "$(VERSION)" --inplace

set-version-%:
	@$(MAKE) helm-set-chart-version-$*
	@$(MAKE) helm-set-app-version-$*

set-app-version-latest-%:
	@yq write $(HELM_PATH)/$*/Chart.yaml appVersion "$(LATEST_TAG)" --inplace
	@yq write $(HELM_PATH)/$*/values.yaml image.tag "$(LATEST_TAG)" --inplace

set-version-all:
	@$(MAKE) helm-set-version-gateway
	@$(MAKE) helm-set-version-eventstore
	@$(MAKE) helm-set-version-commandhandler
	@$(MAKE) helm-set-version-queryhandler
	@$(MAKE) helm-set-version-cluster-bootstrap-reactor

docs:
	@docker run --rm --volume "$(PWD):/helm-docs" -u $(shell id -u) jnorwood/helm-docs:v1.4.0 --template-files=./README.md.gotmpl

.PHONY: test
test:
	bash $(TOOLS_DIR)/render_all_fixtures.sh ./test/pkg/templates_unit_test/testdata
	go test ./test/pkg/templates_unit_test/
