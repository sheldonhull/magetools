// Provide Go linting, formatting and other basic tooling.
package gotools

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"

	"github.com/sheldonhull/magetools/tooling"
)

type Golang mg.Namespace

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

// ⚙️  Init runs all required steps to use this package.
func (Golang) Init() error {
	if err := tooling.InstallTools(toolList); err != nil {
		return err
	}

	return nil
}

// 🔎  Run go test on project.
func (Golang) Test() error {
	var vflag string

	if mg.Verbose() {
		vflag = "-v"
	}

	pterm.Info.Println("Running go test")
	if err := sh.Run("go test", "./...", "-shuffle", "-race", vflag); err != nil {
		return err
	}

	return nil
}

// 🔎  Run golangci-lint without fixing.
func (Golang) Lint() error {
	var vflag string

	if mg.Verbose() {
		vflag = "-v"
	}

	pterm.Info.Println("Running golangci-lint")
	if err := sh.Run("golangci-lint", "run", "./...", vflag); err != nil {
		return err
	}

	return nil
}

// ⚙️ Lint runs golangci-lint tooling using .golangci.yml settings.
func (Golang) Fmt() error {
	var vflag string

	if mg.Verbose() {
		vflag = "-v"
	}

	pterm.Info.Println("Running golangci-lint")
	if err := sh.Run("golangci-lint", "run", "./...", "--fix", vflag); err != nil {
		return err
	}

	return nil
}

// 🧹 Tidy tidies.
func (Golang) Tidy() error {
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return err
	}

	return nil
}
