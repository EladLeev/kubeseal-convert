SHELL := /bin/bash
export GOBIN := $(CWD)/.bin
NAME=kubeseal-convert


.PHONY: build
build:
	GOARCH=amd64 GOOS=darwin go build -o ${NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${NAME}-linux main.go

.PHONY: clean
clean:
	go clean
	rm ${NAME}-darwin
	rm ${NAME}-linux

.PHONY: test
test:
	export GOOGLE_APPLICATION_CREDENTIALS=$(PWD)/test/testdata/mock_gcp_creds.json && \
	go test -v ./...

.PHONY: test_coverage
test_coverage:
	go test -v ./... -coverprofile=coverage.out

.PHONY: clean_test
clean_test:
	go clean -testcache

.PHONY: dep
dep:
	go mod download

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vet
vet:
	go vet

# Too many dependencies on this one, won't work without installing each on of them manually first.
# localstack https://github.com/localstack/localstack
# awslocal https://github.com/localstack/awscli-local
# minikube https://minikube.sigs.k8s.io/docs/start/
# kubectx https://github.com/ahmetb/kubectx
# Helm https://helm.sh/docs/intro/install/
.PHONY: init-stack
init-stack:
	localstack start -d
	minikube start
	kubectx minikube

.PHONY: init-secretsmanager
init-secretsmanager:
	localstack status services --format json | jq -r .secretsmanager
	awslocal secretsmanager create-secret --name MyTestSecret --description "This is a test" --secret-string "{\"user\":\"Dwight_Schrute\",\"password\":\"beet4life\"}"
	awslocal secretsmanager list-secrets

.PHONY: init-sealedsecrets
init-sealedsecrets:
	kubectx minikube
	helm repo add sealed-secrets https://bitnami-labs.github.io/sealed-secrets
	helm install sealed-secrets -n kube-system --set-string fullnameOverride=sealed-secrets-controller sealed-secrets/sealed-secrets --wait

.PHONY: init-vault
init-vault:
	kubectl exec --stdin=true --tty=true vault-0 -- /bin/sh

.PHONY: init-dev
init-dev:  init-stack init-secretsmanager init-sealedsecrets

.PHONY: buildmocks
buildmocks:
	mockery --all --dir "./pkg/"
