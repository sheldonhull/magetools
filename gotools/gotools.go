// Provide Go linting, formatting and other basic tooling.
package gotools

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"

	"github.com/sheldonhull/magetools/tooling"
)

type Go mg.Namespace

// golang tools to ensure are locally vendored.
var toolList = []string{ //nolint:gochecknoglobals // ok to be global for tooling setup
	"github.com/goreleaser/goreleaser@v0.174.1",
	"golang.org/x/tools/cmd/goimports@master",
	"github.com/sqs/goreturns@master",
	"github.com/golangci/golangci-lint/cmd/golangci-lint@master",
	"github.com/dustinkirkland/golang-petname/cmd/petname@master",
	"mvdan.cc/gofumpt@latest",
	"github.com/daixiang0/gci@latest",

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
	if err := sh.RunV("golangci-lint", "run", "--enable-all"); err != nil {
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
func (Go) Fmt() error {
	// var vflag string

	// if mg.Verbose() {
	// 	vflag = "-v"
	// }

	pterm.Info.Println("Running golangci-lint formatter")
	if err := sh.RunV("golangci-lint", "run", "--fix", "--presets", "format", "--fast"); err != nil {
		// if err := golanglint("run", "--fix", "--presets", "format", "--fast"); err != nil {
		pterm.Error.Println("golangci-lint failure")

		return err
	}
	// if err := golanglint("run", "--fix", vflag); err != nil {
	// 	return err
	// }

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

// Lint run golangci-lint.
// func (Go) Linter() {
// 	sh.RunV("golangci-lint", "run")
// }
