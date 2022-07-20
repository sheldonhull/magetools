//go:build mage

// This magefile provides a template to get started quickly using the helpers.
// You can drop this into an existing project instead of the default mage --init generated one to quickly get going.
// You'll still have to run: `go mod tidy` to get the dependencies sorted out.
// go get "github.com/sheldonhull/magetools/ci"
// go get "github.com/sheldonhull/magetools/fancy"
// go get "github.com/sheldonhull/magetools/gotools"

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/ci"
	"github.com/sheldonhull/magetools/fancy"

	// mage:import
	"github.com/sheldonhull/magetools/gotools"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build.
const ptermMargin = 10

// artifactDirectory is a directory containing artifacts for the project and shouldn't be committed to source.
const artifactDirectory = "_artifacts"

const permissionUserReadWriteExecute = 0o0700

// tools is a list of Go tools to install to avoid polluting global modules.
// Gotools module already sets up most of the basic go tools.
// var toolList = []string{ //nolint:gochecknoglobals // ok to be global for tooling setup
// 	"github.com/goreleaser/goreleaser@v0.174.1",
// 	"golang.org/x/tools/cmd/goimports@master",
// 	"github.com/sqs/goreturns@master",
// 	"github.com/golangci/golangci-lint/cmd/golangci-lint@master",
// 	"mvdan.cc/gofumpt@latest",
// 	"github.com/daixiang0/gci@latest",
// }

// createDirectories creates the local working directories for build artifacts and tooling.
func createDirectories() error {
	for _, dir := range []string{artifactDirectory} {
		if err := os.MkdirAll(dir, permissionUserReadWriteExecute); err != nil {
			pterm.Error.Printf("failed to create dir: [%s] with error: %v\n", dir, err)

			return err
		}
		pterm.Success.Printf("✅ [%s] dir created\n", dir)
	}

	return nil
}

// Init runs multiple tasks to initialize all the requirements for running a project for a new contributor.
func Init() error {
	fancy.IntroScreen(ci.IsCI())
	pterm.Success.Println("running Init()...")
	mg.Deps(Clean, createDirectories)
	if err := (gotools.Go{}.Init()); err != nil {
		return err
	}

	return nil
}

// Clean up after yourself.
func Clean() {
	pterm.Success.Println("Cleaning...")
	for _, dir := range []string{artifactDirectory} {
		err := os.RemoveAll(dir)
		if err != nil {
			pterm.Error.Printf("failed to removeall: [%s] with error: %v\n", dir, err)
		}
		pterm.Success.Printf("🧹 [%s] dir removed\n", dir)
	}
	mg.Deps(createDirectories)
}
