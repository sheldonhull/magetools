// Provide Go linting, formatting and other basic tooling.
package gotools

import (
	"os"
	"path/filepath"
	"runtime"

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

// 🔎 Lint runs golangci-lint tooling.
func (Golang) Lint() error {
	pterm.Info.Println("Running golangci-lint")
	if err := tooling.RunTool("golangci-lint", "run", "./..."); err != nil {
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
	if err := tooling.RunTool("gofumpt", "-l", "-w", "."); err != nil {
		return err
	}
	p.Increment()

	p.Title = "goimports"
	if err := tooling.RunTool("goimports", "-w", "."); err != nil {
		return err
	}
	p.Increment()

	p.Title = "gci"
	if err := tooling.RunTool("gci", "-w", "."); err != nil {
		return err
	}
	p.Increment()

	return nil
}

// 🧹 Clean all Go artifacts.
func (Golang) Clean() error {
	pterm.Success.Println("Cleaning...")

	// if runtime.GOOS == "windows" {
	// 	pterm.Error.Println("You are running on Windows. Request support for this command (haven't programmed in .exe extension handling yet)")

	// 	return errors.New("windows support not implemented yet")
	// }
	for _, item := range []string{
		"goreleaser",
		"goimports",
		"goreturns",
		"golangci-lint",
		"petname",
		"gofumpt",
		"gci",
	} {
		// if windows detected, add the exe to the binary path
		var extension string
		if runtime.GOOS == "windows" {
			extension = ".exe"
		}
		toolPath := filepath.Join("_tools", item+extension)
		if _, err := os.Stat(toolPath); err != nil {
			pterm.Info.Printf("🔄 [%s] item not found, bypassed\n", toolPath)

			break
		}
		err := os.RemoveAll(toolPath)
		if err != nil {
			pterm.Error.Printf("failed to remove: [%s] with error: %v\n", toolPath, err)
		}
		pterm.Success.Printf("🧹 [%s] item removed\n", toolPath)
	}

	return nil
}
