VERSION   ?= 0.0.1-local
KUBE_NAMESPACE ?= platform-monoskope-monoskope

export 

# go

go-%:
	@$(MAKE) -f go.mk $*

# helm

helm-%:
	@$(MAKE) -f helm.mk $*