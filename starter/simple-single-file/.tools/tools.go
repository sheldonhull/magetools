//go:build tools
// +build tools

// Place this in build directory, tools directory, or anywhere else to avoid conflict with main package in the same directory.
// Tooling that Mage or other automation tools use, that is _not_ part of the core code base.
// This signifies to Go that these tools are build tooling and not part of the dependency chain for building the project.
// Additionally, it's ignored for everything like go build.
// To ensure these are downloaded, run go mod tidy

package tools

import (
	_ "github.com/dustinkirkland/golang-petname"   // Petname is a tool for generating random names for builds, terraform, or other helpers.
	_ "github.com/gosuri/uiprogress"               // Uiprogress provides progress bars for CLI tools that support go routines.
	_ "github.com/magefile/mage/mg"                // Mg contains helper magefile tasks from mage.
	_ "github.com/magefile/mage/sh"                // Sh contains helper shell tasks from mage.
	_ "github.com/manifoldco/promptui"             // Promptui is a tool for generating interactive prompts used in Mage to help developers.
	_ "github.com/pterm/pterm"                     // Pterm is for cross-platform terminal output.
	_ "github.com/sheldonhull/magetools/fancy"     // Fancy is a tool for generating fancy output headers.
	_ "github.com/sheldonhull/magetools/gittools"  // Gittools contains tasks for setting up tooling like git town and bit.
	_ "github.com/sheldonhull/magetools/gotools"   // Gotools contains tasks for setting up tooling with standard gofumpt, testing output, and setup of standard gopls tooling.
	_ "github.com/sheldonhull/magetools/licensing" // Licensing is a tool for generating license audit capture.
	_ "github.com/sheldonhull/magetools/tooling"   // Tooling has helpers to do low noise updates on actions like go install, go mod tidy, using pterm spinners.
)
