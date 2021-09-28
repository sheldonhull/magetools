package tooling_test

import (
	"testing"

	iz "github.com/matryer/is"

	"github.com/sheldonhull/magetools/tooling"
)

func TestGolang_InstallTools(t *testing.T) {
	is := iz.New(t)
	// pterm.DisableStyling()

	toolList := []string{
		"github.com/goreleaser/goreleaser@v0.174.1",
		"golang.org/x/tools/cmd/goimports@master",
		"github.com/sqs/goreturns@master",
		"github.com/golangci/golangci-lint/cmd/golangci-lint@master",
		"github.com/dustinkirkland/golang-petname/cmd/petname@master",
	}

	err := tooling.InstallTools(toolList)
	is.NoErr(err) // installing a tool should not fail
}
