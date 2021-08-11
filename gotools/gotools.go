// Provide Go linting, formatting and other basic tooling.
package gotools

import (
	"github.com/magefile/mage/mg"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/tooling"
)

type Golang mg.Namespace

// Lint runs golangci-lint tooling.
func (Golang) Lint() error {
	pterm.Info.Println("Running golangci-lint")
	if err := tooling.RunTool("golangci-lint", "run", "./..."); err != nil {
		return err
	}

	return nil
}
