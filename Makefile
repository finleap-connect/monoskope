VERSION   ?= 0.0.1-local

# go

go-%:
	$(MAKE) -f go.mk $*
