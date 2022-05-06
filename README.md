# Hymenaios

[![GitHub](https://img.shields.io/github/license/hymenaios-io/hymenaios)](https://github.com/hymenaios-io/Hymenaios/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/hymenaios-io/Hymenaios)](https://goreportcard.com/report/github.com/hymenaios-io/Hymenaios)
[![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/hymenaios-io/hymenaios?filename=go.mod)](https://go.dev/dl/)
[![GitHub package.json dependency version (subfolder of monorepo)](https://img.shields.io/github/package-json/dependency-version/hymenaios-io/hymenaios/react?filename=web%2Fui%2Freact-app%2Fpackage.json)](https://reactjs.org/)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/hymenaios-io/hymenaios/Test?label=Tests)](https://github.com/hymenaios-io/Hymenaios/actions/workflows/test.yml)


[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/hymenaios-io/hymenaios/Binary%20Build?label=Binary%20Build)](https://github.com/hymenaios-io/Hymenaios/actions/workflows/build-binary.yml)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/hymenaios-io/hymenaios)](https://github.com/hymenaios-io/Hymenaios/releases)
[![GitHub all releases](https://img.shields.io/github/downloads/hymenaios-io/hymenaios/total)](https://github.com/hymenaios-io/Hymenaios/releases)
[![GitHub release (latest by SemVer)](https://img.shields.io/github/downloads/hymenaios-io/hymenaios/latest/total)](https://github.com/hymenaios-io/Hymenaios/releases/latest)

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/hymenaios-io/hymenaios/Docker%20Build?label=Docker%20Build)](https://github.com/hymenaios-io/Hymenaios/actions/workflows/build-docker.yml)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/hymenaios/hymenaios?sort=semver)](https://hub.docker.com/r/hymenaios/hymenaios/tags)
[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/hymenaios/hymenaios?sort=semver)](https://hub.docker.com/r/hymenaios/hymenaios/tags)
[![Docker Pulls](https://img.shields.io/docker/pulls/hymenaios/hymenaios)](https://hub.docker.com/r/hymenaios/hymenaios)

Hymenaios will query websites at a user defined interval for new software releases and then trigger Gotify/Slack notification(s) and/or WebHook(s) when one has been found.
For example, you could set it to monitor the Hymenaios repo ([hymenaios-io/hymenaios](https://github.com/hymenaios-io/hymenaios)). This will query the [GitHub API](https://api.github.com/repos/hymenaios-io/hymenaios/releases) and track the "tag_name" variable. When this variable changes from what it was on a previous query, a GitHub-style WebHook could be sent that triggers  something (like AWX) to update Hymenaios on your server.

##### Table of Contents

- [Hymenaios](#hymenaios)
  - [Demo](#demo)
  - [Command-line arguments](#command-line-arguments)
  - [Building from source](#building-from-source)
    - [Prereqs](#prereqs)
    - [Go changes](#go-changes)
    - [React changes](#react-changes)
  - [Getting started](#config-formatting)
    - [Config formatting](#getting-started)

## Demo

A demo of Hymenaios can be seen on our website [here](https://hymenaios.io/demo).

## Command-line arguments

```bash
$ hymenaios -h
Usage of /usr/local/bin/hymenaios:
  -config.check
        Print the fully-parsed config.
  -config.file string
        Hymenaios configuration file path. (default "config.yml")
  -log.level string
        ERROR, WARN, INFO, VERBOSE or DEBUG (default "INFO")
  -log.timestamps
        Enable timestamps in CLI output.
  -test.gotify string
        Put the name of the Gotify service to send a test message.
  -test.service string
        Put the name of the Service to test the version query.
  -test.slack string
        Put the name of the Slack service to send a test message.
  -web.cert-file string
        HTTPS certificate file path.
  -web.listen-host string
        IP address to listen on for UI, API, and telemetry. (default "0.0.0.0")
  -web.listen-port string
        Port to listen on for UI, API, and telemetry. (default "8080")
  -web.pkey-file string
        HTTPS private key file path.
  -web.route-prefix string
        Prefix for web endpoints (default "/")
```

## Building from source

#### Prereqs

The backend of Hymenaios is built with [Go](https://go.dev/) and the frontend with [React](https://reactjs.org/). The React frontend is built and then [embedded](https://pkg.go.dev/embed) into the Go binary so that those web files can be served.
- [Go 1.18+](https://go.dev/dl/)
- [NodeJS 16](https://nodejs.org/en/download/)

#### Go changes

To see the changes you've made by modifying any of the `.go` files, you must recompile Hymenaios. You could recompile the whole app with a `make build`, but this will also recompile the React components. To save time (and CPU power), you can use the existing React static and recompile just the Go part by running `make go-build`. (Running this in the root dir will produce the `hymenaios` binary)

#### React changes

To see the changes after modifying anything in `web/ui/react-app`, you must recompile both the Go backend as well as the React frontend. This can be done by running `make build`. (Running this in the root dir will produce the `hymenaios` binary)

## Getting started

To get started with Hymenaios, simply download the binary from the [releases page](https://github.com/hymenaios-io/Hymenaios/releases), and setup the config for that binary.

For further help, check out the [Getting Started](https://hymenaios.io/docs/getting-started/) page on our website.

#### Config formatting

The config can be broken down into 6 key areas. ([Further help](https://hymenaios.io/docs/config/))
- [defaults](https://hymenaios.io/docs/config/defaults/) - This is broken down into areas with defaults for [services](https://hymenaios.io/docs/config/defaults/#service-portion), [gotifies](https://hymenaios.io/docs/config/defaults/#gotify-portion), [slacks](https://hymenaios.io/docs/config/defaults/#slack-portion) and [webhooks](https://hymenaios.io/docs/config/defaults/#webhook-portion).
- [settings](https://hymenaios.io/docs/config/settings/) - Settings for the Hymenaios server.
- [service](https://hymenaios.io/docs/config/service/) - A dictionary mapping of all the services to monitor as well as what to notify when a new release is found.
- [gotify](https://hymenaios.io/docs/config/gotify/) - A dictionary mapping of targets for Gotify messages.
- [slack](https://hymenaios.io/docs/config/slack/) - A dictionary mapping of targets for Slack messages.
- [webhook](https://hymenaios.io/docs/config/webhook/) - A dictionary mapping of targets for WebHooks.
