//go:build integration
// +build integration

package tooling_test

import (
	"os"
	"strings"
	"testing"

	iz "github.com/matryer/is"
	"github.com/sheldonhull/magetools/tooling"
)

func TestInstallTools(t *testing.T) {
	if !strings.Contains(strings.ToLower(os.Getenv("GOTESTS")), "slow") {
		t.Skip("GOTESTS should include 'slow' to run this test")
	}
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
// PRESENTATION_TEST=1 go test ./tooling/ -v -tags=integration -run ^TestGo_SilentInit$
// go test -v -timeout 30s -tags mage -run ^TestGoSilentInit$ github.com/sheldonhull/magetools/tooling -v -shuffle on -race
// Try cleaning a specific directory to see output if already ran
// sudo rm -rf $(go env GOPATH)/pkg/mod/github.com/fatih.
func TestGo_SilentInit(t *testing.T) {
	if !strings.Contains(strings.ToLower(os.Getenv("GOTESTS")), "slow") {
		t.Skip("GOTESTS should include 'slow' to run this test")
	}
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

// This is a test for the minimal output, doesn't need to run unless manully invoked with:
//
// This isn't race friendly as it's running go install/mod tidy. Ok to bypass race flag on this
//
// PRESENTATION_TEST=1 go test ./tooling/ -v -tags=integration -run ^Test_SpinnerStdOut$
//
// PRESENTATION_TEST=1 go test ./tooling/ -tags integration -run ^TestGoSpinnerStdOut$  -timeout 30s -v -shuffle on
//
// Try cleaning a specific directory to see output if already ran
//
// find  $(go env GOPATH)/pkg/mod/github.com -name fatih -maxdepth 2 -type d -depth
//
// find  $(go env GOPATH)/pkg/mod/github.com -name fatih -maxdepth 2 -depth -type d -exec sudo rm -rf {} \;.
//
// find  $(go env GOPATH)/pkg/mod/github.com -name gosuri -maxdepth 2 -depth -type d -exec sudo rm -rf {} \;
func TestGoSpinnerStdOut(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// t.Run("run with list of inputs", func(t *testing.T) {
	// 	is := iz.New(t)
	// 	if os.Getenv("PRESENTATION_TEST") != "1" {
	// 		t.Skip("PRESENTATION_TEST != 1 so skipping")
	// 	}

	// 	toolList := []string{
	// 		"github.com/goreleaser/goreleaser@latest",
	// 		"github.com/golangci/golangci-lint/cmd/golangci-lint@master",
	// 		"github.com/dustinkirkland/golang-petname/cmd/petname@master",
	// 		"mvdan.cc/gofumpt@latest",
	// 		"golang.org/x/tools/gopls@latest",
	// 		"github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest",
	// 		"github.com/ramya-rao-a/go-outline@latest",
	// 		"github.com/cweill/gotests/gotests@latest",
	// 		"github.com/fatih/gomodifytags@latest",
	// 	}
	// 	err := tooling.SpinnerStdOut("go", []string{"install"}, toolList)
	// 	is.NoErr(err) // SpinnerStdOut should not fail
	// })

	t.Run("run mod tidy with nil slice", func(t *testing.T) {
		is := iz.New(t)
		if os.Getenv("PRESENTATION_TEST") != "1" {
			t.Skip("PRESENTATION_TEST != 1 so skipping")
		}

		err := tooling.SpinnerStdOut("go", []string{"mod", "tidy"}, nil)
		is.NoErr(err) // SpinnerStdOut should not fail running go mod tidy with nil slice
	})

	t.Run("run mod tidy with empty string slice", func(t *testing.T) {
		is := iz.New(t)
		if os.Getenv("PRESENTATION_TEST") != "1" {
			t.Skip("PRESENTATION_TEST != 1 so skipping")
		}

		err := tooling.SpinnerStdOut("go", []string{"mod", "tidy"}, []string{""})
		is.NoErr(err) // SpinnerStdOut should not fail with empty slice
	})
}
