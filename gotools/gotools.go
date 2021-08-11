// Provide Go linting, formatting and other basic tooling.
package gotools

import (
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/tooling"
)

// Lint runs golangci-lint tooling.
func Lint() error {
	pterm.Info.Println("Running golangci-lint")
	if err := tooling.RunTool("golangci-lint", "--enable-all"); err != nil {
		return err
	}

	return nil
}
