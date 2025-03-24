# X

[![Tests Status](https://github.com/koneksahealth/x/workflows/tests/badge.svg)](https://github.com/koneksahealth/x/actions)
[![Release Status](https://github.com/koneksahealth/x/workflows/release/badge.svg)](https://github.com/koneksahealth/x/actions)
[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)
[![Cov](https://github.com/KoneksaHealth/continuous-delivery/raw/gh-pages/badges/x-coverage.svg?raw=true)](https://github.com/koneksahealth/x/actions)

Sub-repository of packages to be leveraged by other Koneksa Health Go-based applications and services.

* [Releases](./docs/releases)

## Getting Started

### Installing
Use go get to retrieve the package to add it to your GOPATH workspace, or project's Go module dependencies.

```
$ go get github.com/koneksahealth/x
```

To update the package use go get -u to retrieve the latest version of the SDK.

```
$ go get -u github.com/koneksahealth/x
```

### Go Modules
If you are using Go modules, your go get will default to the latest tagged release version of the package. To get a 
specific release version of the package use @<tag> in your go get command.

```
$ go get github.com/koneksahealth/x@v1.0.1
```

To get the latest package repository change use @latest.

```
$ go get github.com/koneksahealth/x@latest
```

### Installing from Private Repositories

#### Github Configuration

```
$ git config --global url."git@github.com:".insteadOf "https://github.com/"
```

#### `GOPRIVATE`
Set the `GOPRIVATE` environment variable.

```
$ export GOPRIVATE=github.com/koneksahealth/*
```
