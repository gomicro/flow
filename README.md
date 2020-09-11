# Flow
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/gomicro/flow/Build/master)](https://github.com/gomicro/flow/actions?query=workflow%3ABuild+branch%3Amaster)
[![Go Reportcard](https://goreportcard.com/badge/github.com/gomicro/flow)](https://goreportcard.com/report/github.com/gomicro/flow)
[![License](https://img.shields.io/github/license/gomicro/flow.svg)](https://github.com/gomicro/flow/blob/master/LICENSE.md)
[![Release](https://img.shields.io/github/release/gomicro/flow.svg)](https://github.com/gomicro/flow/releases/latest)

Flow is a tool for deploying microservices to AWS.

# Why Flow?

Even when you're looking to do a very discrete set of aws-cli actions for build and deploy, depending on your project, there is a significant portion of time spent prepping your environment (upwards of 1 min per build). This prep, which has nothing to do with your project, can consist of setting up Python, maybe install pip, and then install the aws-cli. Additionally you reguarly will run into version conflicts between your installed version of Python and the required version by the aws-cli breaking your build.

Flow aims to drastically reduce that time by allowing for a single download with no additional dependencies. Downloading a single, precompiled binary into your build is simpler, faster, and less fragile.

To additionally note, Flow is not intended to replace the entire toolkit of the [aws-cli](https://github.com/aws/aws-cli). If you need them, please do look to there.

# Installation

## Precompiled Binary

See the [Latest Release](https://github.com/gomicro/flow/releases/latest) page for a download link to the binary compiled for your system.

## From Source

Requires Golang version 1.14 or higher

```
go get github.com/gomicro/flow
```

# Versioning

The tool will be versioned in accordance with [Semver 2.0.0](http://semver.org).  See the [releases](https://github.com/gomicro/forge/releases) section for the latest version.  Until version 1.0.0 the tool is considered to be unstable.

# License
See [LICENSE.md](./LICENSE.md) for more information.

