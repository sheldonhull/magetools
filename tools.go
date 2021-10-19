//go:build tools
// +build tools

// Tooling that Mage or other automation tools use, that is _not_ part of the core code base.
// This signifies to Go that these tools are build tooling and not part of the dependency chain for building the project.
// Additionally, it's ignored for everything like go build.
// To ensure these are downloaded, run go mod tidy

package tools

import (
	_ "github.com/olekukonko/tablewriter"
)
