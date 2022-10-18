//go:build mage

// ⚡ Core Mage Tasks
package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/ci"
	"github.com/sheldonhull/magetools/fancy"

	// mage:import
	_ "github.com/sheldonhull/magetools/docgen"

	// mage:import
	"github.com/sheldonhull/magetools/gotools"

	// mage:import
	_ "github.com/sheldonhull/magetools/precommit"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build.

// artifactDirectory is a directory containing artifacts for the project and shouldn't be committed to source.
const artifactDirectory = ".artifacts"

const permissionUserReadWriteExecute = 0o0777

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
func Init() { //nolint:deadcode // This is not dead code, and I find this insulting golangci-lint.
	fancy.IntroScreen(ci.IsCI())
	pterm.Success.Println("running Init()...")

	mg.SerialDeps(
		Clean,
		createDirectories,
		gotools.Go{}.Tidy,
		// tooling.SilentInstallTools(toolList),
	)
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

type Js mg.Namespace

// Auth setups .yarnrc.yml in your $HOME path.
func (Js) Auth() error {
	pterm.Error.Printfln(
		"`MY_VAR` is a required env var and is missing.\n\nexport MY_VAR=\"mytoken\"\n\n- Add this to your .envrc file\n- ensure `direnv allow` works to load automatically in this project.",
	)
	return nil
}
