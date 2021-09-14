//go:build tool

// Starter template for including tools.go with build oriented tools that aren't necessarily part of the go module dependency tree.
// Might not need this, as the tools versioning will locally provision the tools.

package tools

import (
	_ "github.com/magefile/mage/mg"
	_ "github.com/magefile/mage/sh"
	_ "github.com/pterm/pterm"
	_ "github.com/sheldonhull/magetools/tooling"
	_ "golang.org/x/tools/cmd/stringer"
)
