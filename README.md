# zabbix-cli

zabbix-cli is a command line tool to interact with Zabbix API.

**It is a work in progress.**

## Features

- Export templates YAML
- Import templates YAML

## zabbix API

The API documentation is available [here](https://www.zabbix.com/documentation/6.0/en/manual/api/reference/configuration/export).

## Getting started

Usage is quite simple :

```bash
A CLI tool to interact with Zabbix API

Usage:
  zabbix-cli [command]

Available Commands:
  config      describes how to configure zabbix-cli
  export      export template
  help        Help about any command
  import      import template
  version     print version of zabbix-cli

Flags:
  -h, --help   help for zabbix-cli

Use "zabbix-cli [command] --help" for more information about a command.
```

## Install

### From binary

Download the binary in the release section.

### From Docker image

Docker registry is: sgaunet/zabbix-cli

- The docker image is multi-arch
- It contains only the binary, use multi stage build to copy zabbix-cli in your image

```Dockerfile
FROM sgaunet/zabbix-cli:latest as zabbix-cli

FROM alpine:latest
COPY --from=zabbix-cli /usr/bin/zabbix-cli /usr/local/bin/zabbix-cli
```

## Development

This project is using:

- golang
- [task for development](https://taskfile.dev/#/)
- docker
- [docker buildx](https://github.com/docker/buildx)
- docker manifest
- [goreleaser](https://goreleaser.com/)
- [venom](https://github.com/ovh/venom) : Tests
- [pre-commit](https://pre-commit.com/)

There are hooks executed in the precommit stage. Once the project cloned on your disk, please install pre-commit:

```bash
brew install pre-commit
```

Install tools:

```bash
task dev:install-prereq
```

And install the hooks:

```bash
task dev:install-pre-commit
```

If you like to launch manually the pre-commmit hook:

```bash
task dev:pre-commit
```

## Tests

Tests are done with [venom](https://github.com/ovh/venom).

```bash
cd tests
docker compose up -d
# Check that the stack is up before running the tests
task tests
```
