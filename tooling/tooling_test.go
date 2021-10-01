package tooling_test

import (
	"os"
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

// This is a test for the minimal output, doesn't need to run unless manully invoked with:
// PRESENTATION_TEST=1 go test ./gotools/ -v -tags=integration -run ^TestGo_SilentInit$
// Try cleaning a specific directory to see output if already ran
// sudo rm -rf $(go env GOPATH)/pkg/mod/github.com/fatih.
func TestGo_SilentInit(t *testing.T) {
	is := iz.New(t)
	if os.Getenv("PRESENTATION_TEST") != "1" {
		t.Skip("PRESENTATION_TEST != 1 so skipping")
	}
	// pterm.DisableColor()
	// pterm.DisableStyling()
	toolList := []string{
		"github.com/goreleaser/goreleaser@v0.174.1",
		"github.com/golangci/golangci-lint/cmd/golangci-lint@master",
		"github.com/dustinkirkland/golang-petname/cmd/petname@master",
		"mvdan.cc/gofumpt@latest",
		"golang.org/x/tools/gopls@latest",
		"github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest",
		"github.com/ramya-rao-a/go-outline@latest",
		"github.com/cweill/gotests/gotests@latest",
		"github.com/fatih/gomodifytags@latest",
	}
	err := tooling.SilentInstallTools(toolList)
	is.NoErr(err) // Init should not fail
}
