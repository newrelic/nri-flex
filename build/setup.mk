KUBECTL_BIN ?= bin/kubectl
KUBE_VERSION ?= 1.15.10
KIND_BIN ?= bin/kind
KIND_VERSION ?= 0.7.0

OS := $(shell go env GOOS)

$(KUBECTL_BIN): bin
	@echo "=== $(PROJECT_NAME) === [ setup/kubectl ]: Getting kubectl for $(OS)"
	@(wget -qO $(KUBECTL_BIN) https://storage.googleapis.com/kubernetes-release/release/v$(KUBE_VERSION)/bin/$(OS)/amd64/kubectl)
	@(chmod +x $(KUBECTL_BIN))


$(KIND_BIN): bin
	@echo "=== $(PROJECT_NAME) === [ setup/kind ]: Getting kind for $(OS)"
	@(wget -qO $(KIND_BIN) https://kind.sigs.k8s.io/dl/v$(KIND_VERSION)/kind-$(OS)-amd64)
	@(chmod +x $(KIND_BIN))

.PHONY : setup
setup: $(KUBECTL_BIN) $(KIND_BIN)
