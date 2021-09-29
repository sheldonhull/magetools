// Provide Go linting, formatting and other basic tooling.
package gotools

import (
	"io/ioutil"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/tooling"
	modfile "golang.org/x/mod/modfile"
)

type Go mg.Namespace

// golang tools to ensure are locally vendored.
var toolList = []string{ //nolint:gochecknoglobals // ok to be global for tooling setup
	"github.com/goreleaser/goreleaser@v0.174.1",
	// "golang.org/x/tools/cmd/goimports@master",
	"github.com/sqs/goreturns@master",
	"github.com/golangci/golangci-lint/cmd/golangci-lint@master",
	"github.com/dustinkirkland/golang-petname/cmd/petname@master",
	"mvdan.cc/gofumpt@latest",
	// "github.com/daixiang0/gci@latest",

	// Additionally to simplify init command adding the commands to install from VSCode go installer
	"golang.org/x/tools/gopls@latest",
	"github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest",
	"github.com/ramya-rao-a/go-outline@latest",
	"github.com/cweill/gotests/gotests@latest",
	"github.com/fatih/gomodifytags@latest",
	"github.com/josharian/impl@latest",
	"github.com/haya14busa/goplay/cmd/goplay@latest",
	"github.com/go-delve/delve/cmd/dlv@latest",
}

// getModuleName returns the name from the module file.
// Original help on this was: https://stackoverflow.com/a/63393712/68698
func (Go) GetModuleName() string {
	goModBytes, err := ioutil.ReadFile("go.mod")
	if err != nil {
		pterm.Warning.Println("getModuleName() can't find ./go.mod")
		// Running one more check above the parent directory in case this is running in a test or nested directory for some reason.
		// Only 1 level lookback for now.
		goModBytes, err = ioutil.ReadFile("../go.mod")
		if err != nil {
			pterm.Warning.Println("getModuleName() not able to find ../go.mod")
			return ""
		}
		pterm.Info.Println("found go.mod in ../go.mod")
	}
	modName := modfile.ModulePath(goModBytes)
	return modName
}

// NOTE: this didn't work compared to running with RunV, so I'm commenting out for now.
// golanglint is alias for running golangci-lint.
// var golanglint = sh.RunCmd("golangci-lint") //nolint:gochecknoglobals // ok to be global for tooling setup

// ⚙️  Init runs all required steps to use this package.
func (Go) Init() error {
	if err := tooling.InstallTools(toolList); err != nil {
		return err
	}
	if err := (Go{}.Tidy()); err != nil {
		return err
	}
	pterm.Success.Println("✅  Go Init")
	return nil
}

// 🧪 Run go test on project.
func (Go) Test() error {
	var vflag string

	if mg.Verbose() {
		vflag = "-v"
	}

	pterm.Info.Println("Running go test")
	if err := sh.RunV("go", "test", "./...", "-shuffle", "on", "-race", vflag); err != nil {
		return err
	}
	pterm.Success.Println("✅ Go Test")
	return nil
}

// 🔎  Run golangci-lint without fixing.
func (Go) Lint() error {
	// var vflag string

	// // outFormat := "tab"
	// if mg.Verbose() {
	// 	vflag = "-v"
	// }
	pterm.Info.Println("Running golangci-lint")
	if err := sh.RunV("golangci-lint", "run"); err != nil {
		pterm.Error.Println("golangci-lint failure")

		return err
	}
	// pterm.Info.Println("Running golangci-lint")
	// if err := golanglint("run"); err != nil {
	// 	return err
	// }
	pterm.Success.Println("✅ Go Lint")
	return nil
}

// ⚙️ Lint runs golangci-lint tooling using .golangci.yml settings.
// Recommend setting fast: false in your config and allow tool to set.
// Recommend not setting enable-all: true in config to allow cli to call this for linting.
// REMOVED: 20201-09-29 due to issues found in https://github.com/golangci/golangci-lint/issues/1490
// Will manually invoke desired formatters.
// func (Go) Fmt() error {
// 	// var vflag string

// 	// if mg.Verbose() {
// 	// 	vflag = "-v"
// 	// }

// 	pterm.Info.Println("Running golangci-lint formatter")
// 	if err := sh.RunV("golangci-lint", "run", "--fix", "--enable", "gofumpt,gci", "--fast"); err != nil {
// 		// if err := golanglint("run", "--fix", "--presets", "format", "--fast"); err != nil {
// 		pterm.Error.Println("golangci-lint failure")

// 		return err
// 	}
// 	// if err := golanglint("run", "--fix", vflag); err != nil {
// 	// 	return err
// 	// }
// 	return nil
// }

// ✨ Fmt runs gofumpt.
// Important. Make sure golangci-lint config disables gci, goimports, and gofmt.
// This will perform all the sorting and other linters can cause conflicts in import ordering.
func (Go) Fmt() error {
	pterm.Info.Println("Running gofumpt")
	if err := sh.Run("gofumpt", "-l", "-w", "."); err != nil {
		return err
	}
	pterm.Success.Println("✅ gofumpt")
	return nil
}

// 🧹 Tidy tidies.
func (Go) Tidy() error {
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return err
	}
	pterm.Success.Println("✅ Go Tidy")
	return nil
}

// 🏥 Doctor will provide config details.
func (Go) Doctor() {
	pterm.Info.Println("🏥 Doctor Diagnostic Checks")

	pterm.DefaultSection.Println("🔍 golangci-lint linters with --preset format")
	if err := sh.RunV("golangci-lint", "linters", "--enable", "gofumpt,gci"); err != nil {
		pterm.Error.Println("unable to run golangci-lint")
	}
	pterm.DefaultSection.Println("🔍 golangci-lint linters with --fast")
	if err := sh.RunV("golangci-lint", "linters"); err != nil {
		pterm.Error.Println("unable to run golangci-lint")
	}
	pterm.DefaultSection.Println("🔍  golangci-lint linters with plain run")
	if err := sh.RunV("golangci-lint", "linters"); err != nil {
		pterm.Error.Println("unable to run golangci-lint")
	}

	pterm.Success.Println("Doctor Diagnostic Checks")
}
