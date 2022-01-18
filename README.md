# magetools

![Coverage](https://img.shields.io/badge/Coverage-56.0%25-yellow)

General tooling helpers for simplifying cross repository automation using Mage.

<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-34%25-brightgreen.svg?longCache=true&style=flat)</a>

The test coverage shown is informal, as these aren't setup always with full standard tests.
Primarily the tests just import and run to confirm no errors.
When possible, I do test the functionality.

Mage is a good case of a little magic making things faster and easier to adopt, while making it a bit trickier to test functions IMO.

## Other Go Focused Tools for Task Automation

Worth an honorable mention is [Goyek](https://github.com/goyek/goyek) which has a similar goal to Mage, but approaches with less magic and a similar design to writing tests.
I opted not to use after a few months of working with it due to the alpha level changes occuring in project, and more effort required on adding tasks and remote imports[remote-imports-with-goyek].

If you don't have any desire to use remote imports, then Goyek can be a good option if you want to write in a style like Go testing.

Lastly, another worthy mention is [Gosh](https://github.com/mumoshu/gosh).

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

## Getting Started With Mage

I need a way to bootstrap projects quickly with some tasks, so these directories give two different templates.

1. `simple-single-file`: The simpliest for small projects, contains all tasks in single file.
2. `root-imports-with-tasks-in-subdirectory`: Organized for normal projects, this provides a zero-install run (using `go run mage.go taskname`) and also organizes tasks into subdirectory.
This allows the root of the project to remain clean.
Tasks can then be easily split into files in the tasks directory such as: `mage.git.go, mage.js.go, etc`.

## Why [dot] prefix in directory?

Go commands that run like `go test ./...` automatically ignore dot or underscore prefixed files.

I prefer to the dot prefix and underscore isn't very common from what I've observed in the Go echo system.

## Tip

Include `mage_output_file.go` in your gitignore file to avoid it causing consternation in your git diff monitoring tooling.

## Examples

The examples directory contains random mage examples I've used that don't fit yet into a package, or I haven't had the need to add any tests to reuse this in multiple places.

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
This should bypass the cached public version and call directly to quickly update if the cached go pkg version isn't registering an update.

```shell
GOPRIVATE=github.com/sheldonhull/* go get -u
GOPRIVATE=github.com/sheldonhull/* go get -u github.com/sheldonhull/magetools/gotools@latest
```

```powershell
$ENV:GOPRIVATE='github.com/sheldonhull/*'
```

## Allow Zero Install Run

From the [Mage Docs], see [mage.go](starter/root-imports-with-tasks-in-subdirectory/mage.go).

Run this using: `go run main.go` and it should work just like using `mage` directly.

[Mage Docs]: https://magefile.org/zeroinstall

## Future

Possibly best to setup with a dedicated templating tool in the future (something like Cookiecutter) but for now this is just an easy copy and paste in VSCode.

[remote-imports-with-gokey]: https://github.com/goyek/goyek/discussions/114
