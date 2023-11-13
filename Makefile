# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL := /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# OS information
OS = $(shell go env GOOS)
ARCH = $(shell go env GOARCH)

# Current Operator version
export IMAGE_VERSION ?= 2.5.0-dev
BUNDLE_VERSION ?= $(IMAGE_VERSION)
DEFAULT_NAMESPACE ?= quay.io/cryostat
IMAGE_NAMESPACE ?= $(DEFAULT_NAMESPACE)
OPERATOR_NAME ?= cryostat-operator
OPERATOR_SDK_VERSION ?= v1.31.0
CLUSTER_CLIENT ?= kubectl
IMAGE_TAG_BASE ?= $(IMAGE_NAMESPACE)/$(OPERATOR_NAME)

# Default bundle image tag
BUNDLE_IMG ?= $(IMAGE_TAG_BASE)-bundle:$(BUNDLE_VERSION)
BUNDLE_IMGS ?= $(BUNDLE_IMG)

# Default catalog image tag
CATALOG_IMG ?= $(IMAGE_TAG_BASE)-catalog:$(BUNDLE_VERSION)
ifneq ($(origin CATALOG_BASE_IMG), undefined)
FROM_INDEX_OPT := --from-index $(CATALOG_BASE_IMG)
endif

# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)
# BUNDLE_GEN_FLAGS are the flags passed to the operator-sdk generate bundle command
BUNDLE_GEN_FLAGS ?= -q --overwrite --version $(BUNDLE_VERSION) $(BUNDLE_METADATA_OPTS)
# USE_IMAGE_DIGESTS defines if images are resolved via tags or digests
# You can enable this value if you would like to use SHA Based Digests
# To enable set flag to true
USE_IMAGE_DIGESTS ?= false
ifeq ($(USE_IMAGE_DIGESTS), true)
	BUNDLE_GEN_FLAGS += --use-image-digests
endif
BUNDLE_INSTALL_MODE ?= AllNamespaces

IMAGE_BUILDER ?= podman
# Image URL to use all building/pushing image targets
OPERATOR_IMG ?= $(IMAGE_TAG_BASE):$(IMAGE_VERSION)
	#
# PLATFORMS defines the target platforms for the manager image to provide support to multiple
# architectures. (i.e. make oci-buildx OPERATOR_IMG=quay.io/cryostat/cryostat-operator:latest).
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
MANIFEST_PUSH ?= true

# Name of the application deployed by the operator
export APP_NAME ?= Cryostat

# Images used by the operator
CORE_NAMESPACE ?= $(DEFAULT_NAMESPACE)
CORE_NAME ?= cryostat
CORE_VERSION ?= latest
export CORE_IMG ?= $(CORE_NAMESPACE)/$(CORE_NAME):$(CORE_VERSION)
DATASOURCE_NAMESPACE ?= $(DEFAULT_NAMESPACE)
DATASOURCE_NAME ?= jfr-datasource
DATASOURCE_VERSION ?= latest
export DATASOURCE_IMG ?= $(DATASOURCE_NAMESPACE)/$(DATASOURCE_NAME):$(DATASOURCE_VERSION)
GRAFANA_NAMESPACE ?= $(DEFAULT_NAMESPACE)
GRAFANA_NAME ?= cryostat-grafana-dashboard
GRAFANA_VERSION ?= latest
export GRAFANA_IMG ?= $(GRAFANA_NAMESPACE)/$(GRAFANA_NAME):$(GRAFANA_VERSION)
REPORTS_NAMESPACE ?= $(DEFAULT_NAMESPACE)
REPORTS_NAME ?= cryostat-reports
REPORTS_VERSION ?= latest
export REPORTS_IMG ?= $(REPORTS_NAMESPACE)/$(REPORTS_NAME):$(REPORTS_VERSION)

CERT_MANAGER_VERSION ?= 1.11.5
CERT_MANAGER_MANIFEST ?= \
	https://github.com/cert-manager/cert-manager/releases/download/v$(CERT_MANAGER_VERSION)/cert-manager.yaml

KUSTOMIZE_VERSION ?= 3.8.7
CONTROLLER_TOOLS_VERSION ?= 0.11.1
GOLICENSE_VERSION ?= 1.29.0
OPM_VERSION ?= 1.23.0
# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION ?= 1.26

# Scorecard ImagePullPolicy is hardcoded to IfNotPresent
# See: https://github.com/operator-framework/operator-sdk/pull/4762
#
# Suffix is the timestamp of the image build, compute with: date -u '+%Y%m%d%H%M%S'
CUSTOM_SCORECARD_VERSION ?= 2.5.0-$(shell date -u '+%Y%m%d%H%M%S')
export CUSTOM_SCORECARD_IMG ?= $(IMAGE_TAG_BASE)-scorecard:$(CUSTOM_SCORECARD_VERSION)

DEPLOY_NAMESPACE ?= cryostat-operator-system
TARGET_NAMESPACES ?= $(DEPLOY_NAMESPACE)
SCORECARD_NAMESPACE ?= cryostat-operator-scorecard

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Check whether this is a development or release version
ifneq (,$(shell echo $(IMAGE_VERSION) | grep -iE '(:latest|SNAPSHOT|dev|BETA[[:digit:]]+)$$'))
PULL_POLICY ?= Always
else
PULL_POLICY ?= IfNotPresent
endif
export PULL_POLICY

# Run tests with Ginkgo CLI if available
GINKGO ?= $(shell go env GOPATH)/bin/ginkgo
GO_TEST ?= go test
ifneq ("$(wildcard $(GINKGO))","")
GO_TEST="$(GINKGO)" -cover -output-dir=.
endif

# Optional Red Hat Insights integration
ENABLE_INSIGHTS ?= false
ifeq ($(ENABLE_INSIGHTS), true)
KUSTOMIZE_DIR ?= config/insights
INSIGHTS_PROXY_NAMESPACE ?= quay.io/3scale
INSIGHTS_PROXY_NAME ?= apicast
INSIGHTS_PROXY_VERSION ?= insights-01
export INSIGHTS_PROXY_IMG ?= $(INSIGHTS_PROXY_NAMESPACE)/$(INSIGHTS_PROXY_NAME):$(INSIGHTS_PROXY_VERSION)
export INSIGHTS_BACKEND ?= console.redhat.com
else
KUSTOMIZE_DIR ?= config/default
endif

##@ General

.PHONY: all
all: manager

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Tests

.PHONY: test ## Run tests.
test: test-envtest test-scorecard

.PHONY: test-envtest
test-envtest: generate manifests fmt vet setup-envtest ## Run tests using envtest.
ifneq ($(SKIP_TESTS), true)
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" $(GO_TEST) -v -coverprofile cover.out ./...
endif

.PHONY: test-scorecard
test-scorecard: check_cert_manager kustomize operator-sdk ## Run scorecard tests.
ifneq ($(SKIP_TESTS), true)
	$(call scorecard-setup)
	$(call scorecard-cleanup); \
	trap cleanup EXIT; \
	$(OPERATOR_SDK) scorecard -n $(SCORECARD_NAMESPACE) -s cryostat-scorecard -w 20m $(BUNDLE_IMG) --pod-security=restricted
endif

.PHONY: clean-scorecard
clean-scorecard: operator-sdk ## Clean up scorecard resources.
	- $(call scorecard-cleanup); cleanup

ifneq ($(and $(SCORECARD_REGISTRY_SERVER),$(SCORECARD_REGISTRY_USERNAME),$(SCORECARD_REGISTRY_PASSWORD)),)
SCORECARD_ARGS := --pull-secret-name registry-key --service-account cryostat-scorecard
endif

define scorecard-setup
@$(CLUSTER_CLIENT) get namespace $(SCORECARD_NAMESPACE) >/dev/null 2>&1 &&\
	echo "$(SCORECARD_NAMESPACE) namespace already exists, please remove it with \"make clean-scorecard\"" >&2 && exit 1 || true
$(CLUSTER_CLIENT) create namespace $(SCORECARD_NAMESPACE) && \
	$(CLUSTER_CLIENT) label --overwrite namespace $(SCORECARD_NAMESPACE) pod-security.kubernetes.io/warn=restricted pod-security.kubernetes.io/audit=restricted
cd internal/images/custom-scorecard-tests/rbac/ && $(KUSTOMIZE) edit set namespace $(SCORECARD_NAMESPACE)
$(KUSTOMIZE) build internal/images/custom-scorecard-tests/rbac/ | $(CLUSTER_CLIENT) apply -f -
@if [ -n "$(SCORECARD_ARGS)" ]; then \
	$(CLUSTER_CLIENT) create -n $(SCORECARD_NAMESPACE) secret docker-registry registry-key --docker-server="$(SCORECARD_REGISTRY_SERVER)" \
		--docker-username="$(SCORECARD_REGISTRY_USERNAME)" --docker-password="$(SCORECARD_REGISTRY_PASSWORD)"; \
	$(CLUSTER_CLIENT) patch sa cryostat-scorecard -n $(SCORECARD_NAMESPACE) -p '{"imagePullSecrets": [{"name": "registry-key"}]}'; \
fi
$(OPERATOR_SDK) run bundle -n $(SCORECARD_NAMESPACE) --timeout 20m $(BUNDLE_IMG) --security-context-config=restricted $(SCORECARD_ARGS)
endef

define scorecard-cleanup
function cleanup { \
	(\
	set +e; \
	$(OPERATOR_SDK) cleanup -n $(SCORECARD_NAMESPACE) $(OPERATOR_NAME); \
	$(KUSTOMIZE) build internal/images/custom-scorecard-tests/rbac/ | $(CLUSTER_CLIENT) delete --ignore-not-found=$(ignore-not-found) -f -; \
	$(CLUSTER_CLIENT) delete --ignore-not-found=$(ignore-not-found) -n $(SCORECARD_NAMESPACE) secret registry-key; \
	$(CLUSTER_CLIENT) delete --ignore-not-found=$(ignore-not-found) namespace $(SCORECARD_NAMESPACE); \
	)\
}
endef

##@ Build

.PHONY: manager
manager: manifests generate fmt vet ## Build the manager binary.
	go build -o bin/manager internal/main.go

.PHONY: run
run: manifests generate fmt vet ## Run against the configured Kubernetes cluster in ~/.kube/config
	go run ./internal/main.go

ifndef ignore-not-found
ignore-not-found = false
endif

.PHONY: custom-scorecard-tests
custom-scorecard-tests: fmt vet ## Build the custom scorecard binary.
	cd internal/images/custom-scorecard-tests/ && \
	go build -o bin/cryostat-scorecard-tests main.go

.PHONY: scorecard-build
scorecard-build: custom-scorecard-tests ## Build the custom scorecard OCI image.
	printf '# Code generated by hack/custom.config.yaml.in. DO NOT EDIT.\n' > config/scorecard/patches/custom.config.yaml
	envsubst < hack/custom.config.yaml.in >> config/scorecard/patches/custom.config.yaml
# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' internal/images/custom-scorecard-tests/Dockerfile > internal/images/custom-scorecard-tests/Dockerfile.cross
ifeq ($(IMAGE_BUILDER), docker)
	- $(IMAGE_BUILDER) buildx create --name project-v3-builder
	$(IMAGE_BUILDER) buildx use project-v3-builder
	- $(IMAGE_BUILDER) buildx build --push --platform=$(PLATFORMS) --tag $(CUSTOM_SCORECARD_IMG) -f internal/images/custom-scorecard-tests/Dockerfile.cross .
	- $(IMAGE_BUILDER) buildx rm project-v3-builder
else ifeq ($(IMAGE_BUILDER), podman)
	BUILDAH_FORMAT=docker $(IMAGE_BUILDER) build -f internal/images/custom-scorecard-tests/Dockerfile.cross --manifest $(CUSTOM_SCORECARD_IMG) --platform $(PLATFORMS) . ; \
	if [ "${MANIFEST_PUSH}" = "true" ] ; then \
		$(IMAGE_BUILDER) manifest push $(CUSTOM_SCORECARD_IMG) $(CUSTOM_SCORECARD_IMG) ; \
	fi
else
	$(error unsupported IMAGE_BUILDER: $(IMAGE_BUILDER))
endif
	rm internal/images/custom-scorecard-tests/Dockerfile.cross

.PHONY: oci-build
oci-build: manifests generate fmt vet test-envtest ## Build OCI image for the manager.
	BUILDAH_FORMAT=docker $(IMAGE_BUILDER) build --build-arg TARGETOS=$(OS) --build-arg TARGETARCH=$(ARCH) -t $(OPERATOR_IMG) .

# You may need to be able to push the image for your registry (i.e. if you do not inform a valid value via OPERATOR_IMG=<myregistry/image:<tag>> than the export will fail)
# If using podman, then you can set MANIFEST_PUSH to avoid this behaviour.
# If IMAGE_BUILDER is docker, you need to:
# - able to use docker buildx. More info: https://docs.docker.com/build/buildx/
# - have enable BuildKit, More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# To properly provided solutions that supports more than one platform you should use this option.
.PHONY: oci-buildx
oci-buildx: manifests generate fmt vet test-envtest ## Build OCI image for the manager for cross-platform support.
# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
ifeq ($(IMAGE_BUILDER), docker)
	- $(IMAGE_BUILDER) buildx create --name project-v3-builder
	$(IMAGE_BUILDER) buildx use project-v3-builder
	- $(IMAGE_BUILDER) buildx build --push --platform=$(PLATFORMS) --tag $(OPERATOR_IMG) -f Dockerfile.cross .
	- $(IMAGE_BUILDER) buildx rm project-v3-builder
else ifeq ($(IMAGE_BUILDER), podman)
	BUILDAH_FORMAT=docker $(IMAGE_BUILDER) build -f Dockerfile.cross --manifest $(OPERATOR_IMG) --platform $(PLATFORMS) . ; \
	if [ "${MANIFEST_PUSH}" = "true" ] ; then \
		$(IMAGE_BUILDER) manifest push $(OPERATOR_IMG) $(OPERATOR_IMG) ; \
	fi
else
	$(error unsupported IMAGE_BUILDER: $(IMAGE_BUILDER))
endif
	rm Dockerfile.cross

.PHONY: catalog-build
catalog-build: opm ## Build a catalog image.
	$(OPM) index add --container-tool $(IMAGE_BUILDER) --mode semver --tag $(CATALOG_IMG) --bundles $(BUNDLE_IMGS) $(FROM_INDEX_OPT)

.PHONY: bundle
bundle: manifests kustomize operator-sdk ## Generate bundle manifests and metadata, then validate generated files.
	$(OPERATOR_SDK) generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(OPERATOR_IMG)
	$(KUSTOMIZE) build config/manifests | $(OPERATOR_SDK) generate bundle $(BUNDLE_GEN_FLAGS)
	$(OPERATOR_SDK) bundle validate ./bundle

.PHONY: bundle-build
bundle-build: ## Build the bundle image.
	$(IMAGE_BUILDER) build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

##@ Development

.PHONY: manifests
manifests: controller-gen ## Generate manifests e.g. CRD, RBAC, etc.
	$(CONTROLLER_GEN) rbac:roleName=role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
	envsubst < hack/image_tag_patch.yaml.in > config/default/image_tag_patch.yaml
	envsubst < hack/image_pull_patch.yaml.in > config/default/image_pull_patch.yaml
ifeq ($(ENABLE_INSIGHTS), true)
	envsubst < hack/insights_patch.yaml.in > config/insights/insights_patch.yaml
endif

.PHONY: fmt
fmt: add-license ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	go generate ./...
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

LICENSE_FILE = $(shell pwd)/LICENSE
GO_PACKAGES := $(shell go list -test -f '{{.Dir}}' ./... | sed -e "s|^$$(pwd)||" | cut -d/ -f2 | sort -u)
.PHONY: check-license
check-license: golicense ## Check if license headers are missing from any code files.
	@echo "Checking license..."
	$(GOLICENSE) --config=go-license.yml --verify $(shell find ${GO_PACKAGES} -name "*.go")

.PHONY: add-license
add-license: golicense ## Add license headers to code files.
	@echo "Adding license..."
	$(GOLICENSE) --config=go-license.yml $(shell find ${GO_PACKAGES} -name "*.go")

.PHONY: remove-license
remove-license: golicense ## Remove license headers from code files.
	@echo "Removing license..."
	$(GOLICENSE) --config=go-license.yml --remove $(shell find ${GO_PACKAGES} -name "*.go")

# Local development/testing helpers
ifneq ($(origin SAMPLE_APP_NAMESPACE), undefined)
SAMPLE_APP_FLAGS += -n $(SAMPLE_APP_NAMESPACE)
endif

.PHONY: sample_app
sample_app: ## Deploy sample app.
	$(CLUSTER_CLIENT) apply $(SAMPLE_APP_FLAGS) -f config/samples/sample-app.yaml

.PHONY: undeploy_sample_app
undeploy_sample_app: ## Undeploy sample app.
	$(CLUSTER_CLIENT) delete $(SAMPLE_APP_FLAGS) --ignore-not-found=$(ignore-not-found) -f config/samples/sample-app.yaml

.PHONY: sample_app_agent
sample_app_agent: undeploy_sample_app_agent ## Deploy sample app with Cryostat Agent.
	@if [ -z "${AUTH_TOKEN}" ]; then \
		if [ "${CLUSTER_CLIENT}" = "oc" ]; then\
			AUTH_TOKEN=`oc whoami -t | base64`; \
		else \
			echo "'AUTH_TOKEN' must be specified."; \
			exit 1; \
		fi; \
	fi; \
	$(CLUSTER_CLIENT) apply $(SAMPLE_APP_FLAGS) -f config/samples/sample-app-agent.yaml; \
	$(CLUSTER_CLIENT) set env $(SAMPLE_APP_FLAGS) deployment/quarkus-test-agent CRYOSTAT_AGENT_AUTHORIZATION="Bearer $(AUTH_TOKEN)"

.PHONY: undeploy_sample_app_agent
undeploy_sample_app_agent: ## Undeploy sample app with Cryostat Agent.
	- $(CLUSTER_CLIENT) delete $(SAMPLE_APP_FLAGS) --ignore-not-found=$(ignore-not-found) -f config/samples/sample-app-agent.yaml

.PHONY: cert_manager
cert_manager: remove_cert_manager ## Install cert manager.
	$(CLUSTER_CLIENT) create --validate=false -f $(CERT_MANAGER_MANIFEST)

.PHONY: remove_cert_manager
remove_cert_manager: ## Remove cert manager.
	- $(CLUSTER_CLIENT) delete --ignore-not-found=$(ignore-not-found) -f $(CERT_MANAGER_MANIFEST)

.PHONY: check_cert_manager
check_cert_manager: ## Check cert manager.
	@api_versions=$$($(CLUSTER_CLIENT) api-versions) &&\
	if [ $$(echo "$${api_versions}" | grep -c '^cert-manager.io/v1$$') -eq 0 ]; then\
		if [ "$${DISABLE_SERVICE_TLS}" != "true" ]; then\
			echo 'cert-manager is not installed, install using "make cert_manager" or disable TLS for services by setting DISABLE_SERVICE_TLS to true' >&2 && exit 1;\
		fi;\
	fi

##@ Build Dependencies

LOCALBIN ?= $(shell pwd)/bin
PHONY: local-bin
local-bin: ## Location to install dependencies.
	mkdir -p $(LOCALBIN)

CONTROLLER_GEN = $(LOCALBIN)/controller-gen
.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary. If wrong version is installed, it will be overwritten.
$(CONTROLLER_GEN): local-bin
	test -s $(CONTROLLER_GEN) && $(CONTROLLER_GEN) --version | grep -q $(CONTROLLER_TOOLS_VERSION) || \
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@v$(CONTROLLER_TOOLS_VERSION)

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
KUSTOMIZE = $(LOCALBIN)/kustomize
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary. If wrong version is installed, it will be removed before downloading.
$(KUSTOMIZE): local-bin
	@if test -x $(LOCALBIN)/kustomize && ! $(LOCALBIN)/kustomize version | grep -q $(KUSTOMIZE_VERSION); then \
		echo "$(LOCALBIN)/kustomize version is not expected $(KUSTOMIZE_VERSION). Removing it before installing."; \
		rm -rf $(LOCALBIN)/kustomize; \
	fi
	test -s $(KUSTOMIZE) || { curl -Ss $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); }

GOLICENSE = $(LOCALBIN)/go-license
.PHONY: golicense
golicense: $(GOLICENSE) ## Download go-license locally if necessary.
$(GOLICENSE): local-bin
	test -s $(GOLICENSE) || GOBIN=$(LOCALBIN) go install github.com/palantir/go-license@v$(GOLICENSE_VERSION)

ENVTEST = $(LOCALBIN)/setup-envtest
.PHONY: setup-envtest
setup-envtest: $(ENVTEST) ## Download setup-envtest locally if necessary.
$(ENVTEST): local-bin
	test -s $(ENVTEST) || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

OPM = $(LOCALBIN)/opm
.PHONY: opm
opm: $(OPM) ## Download opm locally if necessary.
$(OPM): local-bin
	test -s $(OPM) || \
	{ \
	set -e ;\
	curl -sSLo $(OPM) https://github.com/operator-framework/operator-registry/releases/download/v$(OPM_VERSION)/$(OS)-$(ARCH)-opm ;\
	chmod +x $(OPM) ;\
	}

.PHONY: operator-sdk
OPERATOR_SDK ?= $(LOCALBIN)/operator-sdk
operator-sdk: ## Download operator-sdk locally if necessary.
ifeq (,$(wildcard $(OPERATOR_SDK)))
ifeq (, $(shell which operator-sdk 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p $(dir $(OPERATOR_SDK)) ;\
	OS=$(shell go env GOOS) && ARCH=$(shell go env GOARCH) && \
	curl -sSLo $(OPERATOR_SDK) https://github.com/operator-framework/operator-sdk/releases/download/$(OPERATOR_SDK_VERSION)/operator-sdk_$${OS}_$${ARCH} ;\
	chmod +x $(OPERATOR_SDK) ;\
	}
else
OPERATOR_SDK = $(shell which operator-sdk)
endif
endif

##@ Deployment

.PHONY: install
install: manifests kustomize ## Install CRDs into the cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | $(CLUSTER_CLIENT) apply -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the cluster specified in ~/.kube/config.
	- $(KUSTOMIZE) build config/crd | $(CLUSTER_CLIENT) delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: predeploy
predeploy:
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(OPERATOR_IMG)
	cd config/default && $(KUSTOMIZE) edit set namespace $(DEPLOY_NAMESPACE)

.PHONY: print_deploy_config
print_deploy_config: predeploy ## Print deployment configurations for the controller.
	$(KUSTOMIZE) build $(KUSTOMIZE_DIR)

.PHONY: deploy
deploy: check_cert_manager manifests kustomize predeploy ## Deploy controller in the configured cluster in ~/.kube/config
	$(KUSTOMIZE) build $(KUSTOMIZE_DIR) | $(CLUSTER_CLIENT) apply -f -
ifeq ($(DISABLE_SERVICE_TLS), true)
	@echo "Disabling TLS for in-cluster communication between Services"
	@$(CLUSTER_CLIENT) -n $(DEPLOY_NAMESPACE) set env deployment/cryostat-operator-controller-manager DISABLE_SERVICE_TLS=true
endif

.PHONY: undeploy
undeploy: ## Undeploy controller from the configured cluster in ~/.kube/config.
	- $(CLUSTER_CLIENT) delete --ignore-not-found=$(ignore-not-found) -f config/samples/operator_v1beta1_cryostat.yaml
	- $(CLUSTER_CLIENT) delete --ignore-not-found=$(ignore-not-found) -f config/samples/operator_v1beta1_clustercryostat.yaml
	- $(KUSTOMIZE) build $(KUSTOMIZE_DIR) | $(CLUSTER_CLIENT) delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy_bundle
deploy_bundle: check_cert_manager undeploy_bundle ## Deploy the controller in the bundle format with OLM.
	$(OPERATOR_SDK) run bundle --install-mode $(BUNDLE_INSTALL_MODE) $(BUNDLE_IMG)
ifeq ($(DISABLE_SERVICE_TLS), true)
	@echo "Disabling TLS for in-cluster communication between Services"
	@current_ns=`$(CLUSTER_CLIENT) config view --minify -o 'jsonpath={.contexts[0].context.namespace}'` && \
	if [ -z "$${current_ns}" ]; then \
		echo "Failed to determine Namespace in current context" >&2; \
		exit 1; \
	fi; \
	set -f -- `$(CLUSTER_CLIENT) get sub -l "operators.coreos.com/$(OPERATOR_NAME).$${current_ns}" -o name` && \
	if [ "$${#}" -ne 1 ]; then \
		echo -e "Expected 1 Subscription, found $${#}:\n$${@}" >&2; \
		exit 1; \
	fi; \
	$(CLUSTER_CLIENT) patch --type=merge -p '{"spec":{"config":{"env":[{"name":"DISABLE_SERVICE_TLS","value":"true"}]}}}' "$${1}"
endif

.PHONY: undeploy_bundle
undeploy_bundle: operator-sdk ## Undeploy the controller in the bundle format with OLM.
	- $(OPERATOR_SDK) cleanup $(OPERATOR_NAME)

.PHONY: create_cryostat_cr
create_cryostat_cr: destroy_cryostat_cr ## Create a namespaced Cryostat instance.
	$(CLUSTER_CLIENT) create -f config/samples/operator_v1beta1_cryostat.yaml

.PHONY: create_clustercryostat_cr
create_clustercryostat_cr: destroy_clustercryostat_cr ## Create a cluster-wide Cryostat instance.
	target_ns_json=$$(jq -nc '$$ARGS.positional' --args -- $(TARGET_NAMESPACES)) && \
	$(CLUSTER_CLIENT) patch -f config/samples/operator_v1beta1_clustercryostat.yaml --local=true --type=merge \
	-p "{\"spec\": {\"installNamespace\": \"$(DEPLOY_NAMESPACE)\", \"targetNamespaces\": $$target_ns_json}}" -o yaml | \
	oc apply -f -

.PHONY: destroy_cryostat_cr
destroy_cryostat_cr: ## Delete a namespaced Cryostat instance.
	- $(CLUSTER_CLIENT) delete --ignore-not-found=$(ignore-not-found) -f config/samples/operator_v1beta1_cryostat.yaml

.PHONY: destroy_clustercryostat_cr
destroy_clustercryostat_cr: ## Delete a cluster-wide Cryostat instance.
	- $(CLUSTER_CLIENT) delete --ignore-not-found=$(ignore-not-found) -f config/samples/operator_v1beta1_clustercryostat.yaml
