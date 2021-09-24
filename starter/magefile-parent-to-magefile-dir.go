//go:build mage
// +build mage

package main

// If you are importing a remote mage config setup, you should put the import in here, as nested target discovery isn't supported.
// For example: "github.com/sheldonhull/magetools/gotools" provides a preset list of tasks that will be automatically discovered by mage on import without any new code.
// NOTE: some of these are sourced as common helper libraries, but at anytime they can just be copied into magefiles directory and internalized on demand.
// Only importing common build specific helpers from public repos (build a generic docker image, run a go linting tool, changelog tooling etc)
//
// - gotools: Provides init command to just setup golangci-lint and formatters, and common formatter running sequence of a gofmt, goimports, goreturns, and gofumpt.
// - licensing: Uses a google project from github to inventory licenses for problematic licenses, as well as internally vendor the license files for adherence to MIT and other tools.

// If you are importing a remote mage config setup, you should put the import in here, as nested target discovery isn't supported.
// For example: "github.com/sheldonhull/magetools/gotools" provides a preset list of tasks that will be automatically discovered by mage on import without any new code.
// Consuming the package as a standard Go library doesn't go here, but in the magefiles/mage/magefile.go file.

import (

	// mage:import
	_ "github.com/sheldonhull/magetools/gotools"

	// mage:import
	_ "github.com/sheldonhull/magetools/gotools"

	// mage:import
	_ "github.com/sheldonhull/magetools/licensing"
	_ "mycurrentrepo/magefiles"
)