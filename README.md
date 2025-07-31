[![GitHub release](https://img.shields.io/github/release/sgaunet/zabbix-cli.svg)](https://github.com/sgaunet/zabbix-cli/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/zabbix-cli)](https://goreportcard.com/report/github.com/sgaunet/zabbix-cli)
![GitHub Downloads](https://img.shields.io/github/downloads/sgaunet/zabbix-cli/total)
![Test Coverage](https://raw.githubusercontent.com/wiki/sgaunet/zabbix-cli/coverage-badge.svg)
![Coverage CI](https://github.com/sgaunet/zabbix-cli/actions/workflows/coverage.yml/badge.svg)
[![GoDoc](https://godoc.org/github.com/sgaunet/zabbix-cli?status.svg)](https://godoc.org/github.com/sgaunet/zabbix-cli)
[![License](https://img.shields.io/github/license/sgaunet/zabbix-cli.svg)](LICENSE)

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

## 🕐 Project Status: Low Priority

This project is not under active development. While the project remains functional and available for use, please be aware of the following:

### What this means:
- **Response times will be longer** - Issues and pull requests may take weeks or months to be reviewed
- **Updates will be infrequent** - New features and non-critical bug fixes will be rare
- **Support is limited** - Questions and discussions may not receive timely responses

### We still welcome:
- 🐛 **Bug reports** - Critical issues will eventually be addressed
- 🔧 **Pull requests** - Well-tested contributions are appreciated
- 💡 **Feature requests** - Ideas will be considered for future development cycles
- 📖 **Documentation improvements** - Always helpful for the community

### Before contributing:
1. **Check existing issues** - Your concern may already be documented
2. **Be patient** - Responses may take considerable time
3. **Be self-sufficient** - Be prepared to fork and maintain your own version if needed
4. **Keep it simple** - Small, focused changes are more likely to be merged

### Alternative options:
If you need active support or rapid development:
- Look for actively maintained alternatives
- Reach out to discuss taking over maintenance

We appreciate your understanding and patience. This project remains important to us, but current priorities limit our ability to provide regular updates and support.
