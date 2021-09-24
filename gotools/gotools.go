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
}

// ⚙️  Init runs all required steps to use this package.
func (Golang) Init() error {
	if err := tooling.InstallTools(toolList); err != nil {
		return err
	}

	return nil
}

// 🔎  Run golangci-lint and fix by default. Pass [true] to set checkOnly and not auto-fix.
func (Golang) Lint(checkOnly bool) error {
	var vflag string
	fx := "--fix"

	if mg.Verbose() {
		vflag = "-v"
	}

	if checkOnly {
		fx = ""
	}
	pterm.Info.Println("Running golangci-lint")
	if err := sh.Run("golangci-lint", "run", "./...", fx, vflag); err != nil {
		return err
	}

	return nil
}

// ⚙️ Lint runs golangci-lint tooling.
func (Golang) Fmt() error {
	pterm.Info.Println("Running gofmt, gofumpt, goimports, and gci ")
	p, _ := pterm.DefaultProgressbar.WithTotal(4).WithTitle("running formatters").WithRemoveWhenDone(true).Start() //nolint:gomnd
	defer func() {
		p.Title = "formatting completed"
		_, _ = p.Stop()
		pterm.Success.Printf("fmt complete: %s\n", p.GetElapsedTime().String())
	}()

	p.Title = "gofmt"
	if err := sh.Run("gofmt", "-s", "-w", "."); err != nil {
		return err
	}
	p.Increment()
	p.Title = "gofumpt"
	if err := sh.Run("gofumpt", "-l", "-w", "."); err != nil {
		return err
	}
	p.Increment()

	p.Title = "goimports"
	if err := sh.Run("goimports", "-w", "."); err != nil {
		return err
	}
	p.Increment()

	p.Title = "gci"
	if err := sh.Run("gci", "-w", "."); err != nil {
		return err
	}
	p.Increment()

	return nil
}

// 🧹 Tidy tidies.
func (Golang) Tidy() error {
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return err
	}

	return nil
}
