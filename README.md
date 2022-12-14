# kubeseal-convert
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/EladLeev/kubeseal-convert/Build%20Package)
[![Go Report Card](https://goreportcard.com/badge/github.com/eladleev/kubeseal-convert)](https://goreportcard.com/report/github.com/eladleev/kubeseal-convert)

The missing part of [Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets). :closed_lock_with_key:

## Motivation
`kubeseal-convert` aims to reduce the friction of importing secrets from a pre-existing secret management systems (e.g. Vault, AWS Secrets Manager, etc) into a `SealedSecret`.  
Instaed of:
1. Going into AWS Secret Manager
2. Retrieve the secret who needs to be migrated
3. Create a "normal" k8s secret
4. Fill out the values on the secret
5. Run `kubeseal`

Just run `kubeseal-convert` with the secret path.

Table of Contents
-----------------

- [kubeseal-convert](#kubeseal-convert)
  - [Motivation](#motivation)
  - [Table of Contents](#table-of-contents)
  - [Flags \& Options](#flags--options)
    - [Flags](#flags)
  - [Supported SM Systems](#supported-sm-systems)
    - [AWS Secrets Manager](#aws-secrets-manager)
    - [Hashicorp Vault](#hashicorp-vault)
  - [Build from source](#build-from-source)
    - [Prerequisites](#prerequisites)
    - [Building Steps](#building-steps)
  - [Examples](#examples)
  - [Contributing](#contributing)
  - [License](#license)

## Flags & Options
Same as the `kubeseal` command, `kubeseal-convert` is un-opinionated. It won't commit the secret to Git, apply it to the cluster, or save it on a specific path.  
The `SealedSecret` will be printed to `STDOUT`. You can run it as is, as part of CI, or as part of a Job.

```shell
./kubeseal-convert <SECRETS_STORE> <PATH> --namespace <NS_NAME> --name <SECRET_NAME>
```
### Flags
| Name                  | Description                                                            | Require | Type       |
| --------------------- | ---------------------------------------------------------------------- | ------- | ---------- |
| `-n`, `--name`        | The Sealed Secret name.                                                | `V`     | `string`   |
| `--namespace`         | The Sealed Secret namespace. If not specified, taken from k8s context. |         | `string`   |
| `-a`, `--annotations` | Sets k8s annotations. KV pairs, comma separated.                       |         | `[]string` |
| `-l`, `--labels`      | Sets k8s lables. KV pairs, comma separated.                            |         | `[]string` |
|                       |                                                                        |         |            |
| `-h`, `--help`        | Display help.                                                          |         | `none`     |
| `-v`, `--version`     | Display version.                                                       |         | `none`     |


## Supported SM Systems
:white_check_mark: AWS Secrets Manager  
:white_check_mark: Hashicorp Vault  
:question: Google Secrets Manager  
:question: Azure Key Vault

### AWS Secrets Manager
The AWS client rely on AWS local configuration variables - config file, environment variables, etc.
### Hashicorp Vault
In order to work with the Vault provider, two environment variables needs to be set - `VAULT_TOKEN` and `VAULT_ADDR`.  
Currently, only [`kv-v2`](https://developer.hashicorp.com/vault/docs/secrets/kv/kv-v2) is supported.

## Build from source

### Prerequisites

* Go version 1.19+
* `make` command installed
* `kubeseal` command installed, and a valid communication to the sealed secrets controller.

### Building Steps

1. Clone this repository
```shell
git clone https://github.com/EladLeev/kubeseal-convert && cd kubeseal-convert
```
2. Build using Makefile
```shell
make build
```
This command will generate binaries for Mac (aka Darwin), Linux and Windows platforms:
```shell
./out/
├── darwin
│   ├── amd64
│   │   └── kubeseal-convert
│   └── arm64
│       └── kubeseal-convert
├── linux
│   ├── amd64
│   │   └── kubeseal-convert
│   └── arm64
│       └── kubeseal-convert
└── windows
    ├── amd64
    │   └── kubeseal-convert
    └── arm64
        └── kubeseal-convert
```
In order to build a specific platform run one of these commands:
```shell
make build-darwin
make build-linux
make build-windows
```


3. **[optional]** Set up local env for testing
```shell
make init-dev
```
4. **[optional]** Run the [example](#examples)

## Examples
```shell
./kubeseal-convert sm MyTestSecret --namespace test-ns --name test-secret --annotations converted-by=kubeseal-convert,env=dev --labels test=abc > secret.yaml
```
or
```shell
./kubeseal-convert vlt "mydomain/data/MyTestSecret" --namespace test-ns --name test-secret --annotations converted-by=kubeseal-convert,src=vault --labels test=abc > secret.yaml
```
This will:  
1. Retrieve a secret called `MyTestSecret` from AWS Secrets Manager / Hashicorp Vault
2. Create it on `test-ns` namespace
3. Call it `test-secret`
4. Add few annotations and labels
5. Save it as `secret.yaml` to be push to the repo safely

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details of submitting a pull requests.

## License

This project is licensed under the Apache License - see the [LICENSE](LICENSE) file for details.

