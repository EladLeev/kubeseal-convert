NAME=kubeseal-convert

build:
	GOARCH=amd64 GOOS=darwin go build -o ${NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${NAME}-linux main.go

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
init-stack:
	localstack start -d
	minikube start
	kubectx minikube

init-secretmanager:
	localstack status services --format json | jq -r .secretsmanager
	awslocal secretsmanager create-secret --name MyTestSecret --description "This is a test" --secret-string "{\"user\":\"Dwight_Schrute\",\"password\":\"beet4life\"}"
	awslocal secretsmanager list-secrets

init-sealedsecrets:
	kubectx minikube
	helm repo add sealed-secrets https://bitnami-labs.github.io/sealed-secrets
	helm install sealed-secrets -n kube-system --set-string fullnameOverride=sealed-secrets-controller sealed-secrets/sealed-secrets --wait

init-vault:
	kubectl exec --stdin=true --tty=true vault-0 -- /bin/sh

init-dev:  init-stack init-secretmanager