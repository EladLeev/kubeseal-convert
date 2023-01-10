SHELL := /bin/bash
export GOBIN := $(CWD)/.bin

NAME := kubeseal-convert
VERSION := "0.0.2"
GO_MODULE := github.com/eladleev/kubeseal-convert

run:
	go run main.go 
clean:
	go clean
	rm ${NAME}-darwin
	rm ${NAME}-linux

test:
	go test -v ./...

test_coverage:
	go test -v ./... -coverprofile=coverage.out

dep:
	go mod download

tidy:
	go mod tidy

vet:
	go vet

# Too many dependencies on this one, won't work without installing each on of them manually first.
# localstack https://github.com/localstack/localstack
# awslocal https://github.com/localstack/awscli-local
# minikube https://minikube.sigs.k8s.io/docs/start/
# kubectx https://github.com/ahmetb/kubectx
# Helm https://helm.sh/docs/intro/install/
init-stack:
	localstack start -d
	minikube start
	kubectx minikube

init-secretsmanager:
	localstack status services --format json | jq -r .secretsmanager
	awslocal secretsmanager create-secret --name MyTestSecret --description "This is a test" --secret-string "{\"user\":\"Dwight_Schrute\",\"password\":\"beet4life\"}"
	awslocal secretsmanager list-secrets

init-sealedsecrets:
	kubectx minikube
	helm repo add sealed-secrets https://bitnami-labs.github.io/sealed-secrets
	helm install sealed-secrets -n kube-system --set-string fullnameOverride=sealed-secrets-controller sealed-secrets/sealed-secrets --wait

init-vault:
	kubectl exec --stdin=true --tty=true vault-0 -- /bin/sh

init-dev:  init-stack init-secretsmanager init-sealedsecrets

buildmocks:
	mockery --all --dir "./pkg/"

# Build
build: build-linux build-darwin build-windows

build-os:
	CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build  -ldflags="-X '${GO_MODULE}/cmd/kubeseal-convert.version=${VERSION}' -w -extldflags '-static'" -o out/${OS}/${ARCH}/${NAME}  ./main.go

build-darwin:
	OS=darwin ARCH=amd64 make build-os
	OS=darwin ARCH=arm64 make build-os

build-linux:
	OS=linux ARCH=amd64 make build-os
	OS=linux ARCH=arm64 make build-os

build-windows:
	OS=windows ARCH=amd64 make build-os
	OS=windows ARCH=arm64 make build-os
