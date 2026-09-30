CONTROLLER_GEN_VERSION ?= v0.19.0
CONTROLLER_GEN = go run sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)

.PHONY: all
all: manifests generate fmt vet build

.PHONY: manifests
manifests:
	$(CONTROLLER_GEN) crd paths=./api/... output:crd:artifacts:config=config/crd/bases
	$(CONTROLLER_GEN) rbac:roleName=threat-policy-manager-role webhook paths=./... output:rbac:artifacts:config=config/rbac

.PHONY: generate
generate:
	$(CONTROLLER_GEN) object:headerFile=hack/boilerplate.go.txt paths=./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build: fmt vet
	go build -o bin/threat-policy .

.PHONY: docker-build
docker-build:
	docker build -t threat-policy:dev .
