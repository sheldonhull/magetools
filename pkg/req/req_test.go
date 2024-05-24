package req_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitfield/script"
	"github.com/magefile/mage/sh"
	iz "github.com/matryer/is"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/pkg/req"
)

func Test_GetGoPath(t *testing.T) {
	is := iz.New(t)
	// pterm.DisableOutput()
	// pterm.EnableDebugMessages()
	pterm.DisableStyling()
	p := script.Exec("go env GOPATH")
	s, _ := p.String()
	want := strings.TrimSpace(s)
	got := req.GetGoPath()
	is.Equal(want, got) // GOPATH should return the same as go cli.
}

func Test_QualifyGoBinary(t *testing.T) {
	is := iz.New(t)
	// pterm.DisableOutput()
	// pterm.EnableDebugMessages()
	pterm.DisableStyling()
	app := "gofumpt"
	got, err := req.QualifyGoBinary(app)
	is.NoErr(err)      // QualifyGoBinary should not error when resolving the path for the file.
	is.True(got != "") // Binary was found in some manner.
}

// Test_ResolveBinaryByInstall just verifies the tool exists, as this could be in global, or any other tooling location.
// It no longer needs to be in a subdirectory.
func Test_ResolveBinaryByInstall(t *testing.T) {
	is := iz.New(t)
	var got, want string
	var err error
	pterm.DisableOutput()
	// pterm.EnableDebugMessages()
	pterm.DisableStyling()
	app := "gofumpt"
	goInstallCmd := "mvdan.cc/gofumpt@latest"
	want = filepath.Join(
		req.GetGoPath(),
		"bin",
		app,
	) // THIS IS OPTIONAL!!! Might be installed via another method (binaries for example, so no longer requiring path, just cleaning up in case it needs it).
	pterm.Debug.Printfln("want: filepath: %s", want)
	got, err = req.ResolveBinaryByInstall(app, goInstallCmd)
	is.NoErr(err)      // ResolveBinaryByInstall should not error.
	is.True(got != "") // Binary was found in some manner.
	_ = sh.Rm(want)    // Remove the gofumpt binary to ensure we can test failure case.
	got, err = req.ResolveBinaryByInstall(app, goInstallCmd)
	is.NoErr(err)      // ResolveBinaryByInstall should not error after reinstalling.
	is.True(got != "") // Binary was found in some manner.
}

func Test_AddGoPkgBinToPath(t *testing.T) {
	is := iz.New(t)
	// pterm.DisableOutput()
	// pterm.EnableDebugMessages()
	pterm.DisableStyling()
	p := script.Exec("go env GOPATH")
	s, _ := p.String()
	want := strings.TrimSpace(s)
	got := req.GetGoPath()
	is.Equal(want, got) // GOPATH should return the same as go cli.
}
