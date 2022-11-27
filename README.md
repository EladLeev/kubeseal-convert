# kubeseal-convert [wip]

The missing part of [Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets). :closed_lock_with_key:

## Motivation
`kubeseal-convert` aims to reduce the friction of importing secrets from a pre-existing secrets management system (e.g. Vault, AWS Secrets Manager, etc) into a `SealedSecret`.  
Instaed of:
1. Going into AWS Secret Manager
2. Retrieve the secret who needs to be migrated
3. Create a "normal" k8s secret
4. Fill out the values on the secret
5. Run `kubeseal`

Just run `kubeseal-convert` with the secret path.

Table of Contents
-----------------

- [kubeseal-convert \[wip\]](#kubeseal-convert-wip)
  - [Motivation](#motivation)
  - [Table of Contents](#table-of-contents)
  - [Flags \& Options](#flags--options)
    - [Flags](#flags)
  - [Supported SM Systems](#supported-sm-systems)
  - [Build from source](#build-from-source)
    - [Prerequisites](#prerequisites)
    - [Building Steps](#building-steps)
  - [Examples](#examples)
  - [Contributing](#contributing)
  - [License](#license)

## Flags & Options
Same as the `kubeseal` command, `kubeseal-convert` is un-opinionated. It won't commit the secret to Git, apply it to the cluster, or save it on a specific path.  
The `SealedSecret` will be printed to `STDOUT`. You can run it as is, as part of CI, or as part of a Job.

```bash
./kubeseal-convert <SECRETS_STORE> <PATH> --namespace <NS_NAME> --name <SECRET_NAME>
```
### Flags
| Name                    | Description                                                            | Require | Type       |
| ----------------------- | ---------------------------------------------------------------------- | ------- | ---------- |
| `-n`, `--name`          | The Sealed Secret name.                                                | `V`     | `string`   |
| `namespace`             | The Sealed Secret namespace. If not specified, taken from k8s context. |         | `string`   |
| `-a`, `--annotations`   | Sets k8s annotations. KV pairs, comma separated.                       |         | `[]string` |
| `-l`, `--labels`        | Sets k8s lables. KV pairs, comma separated.                            |         | `[]string` |
|                         |                                                                        |         |            |
| `-h`, `--help`          | Display help.                                                          |         | `none`     |
| `-v`, `--version`       | Display version.                                                       |         | `none`     |


## Supported SM Systems
:white_check_mark: AWS Secrets Manager  
:hourglass_flowing_sand: Vault (wip)  
:question: Google Secrets Manager  
:question: Azure Key Vault

## Build from source

### Prerequisites

* Go version 1.19+
* `make` command installed
* `kubeseal` command installed, and a valid communication to the sealed secrets controller.

### Building Steps

1. Clone this repository
```bash
git clone https://github.com/EladLeev/kubeseal-convert && cd kubeseal-convert
```
2. Build using Makefile
```bash
make build
```
3. **[optional]** Set up local env for testing
```bash
make init-dev
```
4. **[optional]** Run the [example](#examples)

## Examples
```bash
./kubeseal-convert sm MyTestSecret --namespace test-ns --name test-secret --annotations converted-by=kubeseal-convert,env=dev --labels test=abc
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details of submitting a pull requests.

## License

This project is licensed under the Apache License - see the [LICENSE](LICENSE) file for details.

