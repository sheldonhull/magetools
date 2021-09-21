//go:build mage
// +build mage

package main

// If you are importing a remote mage config setup, you should put the import in here, as nested target discovery isn't supported.
// For example: "github.com/sheldonhull/magetools/gotools" provides a preset list of tasks that will be automatically discovered by mage on import without any new code.

import (

	// mage:import
	_ "github.com/sheldonhull/magetools/gotools"

	// mage:import
	_ "github.com/sheldonhull/magetools/licensing"
	_ "mycurrentrepo/magefiles"
)
