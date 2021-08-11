# magetools

General tooling helpers for simplifying cross repository automation using Mage.

<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-100%25-brightgreen.svg?longCache=true&style=flat)</a>

The test coverage shown is informal, as these aren't setup always with full standard tests.
Primarily the tests just import and run to confirm no errors.
When possible, I do test the functionality.

## Remote Packages

[How to Use Importing](https://magefile.org/importing/)

## Mage

Install with Go 1.16+

```shell
go install github.com/magefile/mage@latest
```

## Mage-Select

Nice little interactive prompt.

Alias: `mage-select` as `mages`

```shell
go install github.com/iwittkau/mage-select@latest
```

## How To Use

```go
// mage:import
_ "github.com/sheldonhull/magetools/ci"
```

Namespaced

```go
// mage:import ci
"github.com/sheldonhull/magetools/ci"
```

## Update

Quickly refresh library by running this on caller.
This should bypass the cached public version and call directly.

```shell
GOPRIVATE=github.com/sheldonhull/* go get -u
GOPRIVATE=github.com/sheldonhull/* go get -u github.com/sheldonhull/magetools/gotools@latest
```
