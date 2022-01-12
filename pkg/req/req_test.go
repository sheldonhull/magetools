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
	want := filepath.Join(req.GetGoPath(), "bin", app)
	got, err := req.QualifyGoBinary(app)
	is.NoErr(err)       // QualifyGoBinary should not error when resolving the path for the file
	is.Equal(want, got) // Filepath should match
}

func Test_ResolveBinaryByInstall(t *testing.T) {
	is := iz.New(t)
	var got, want string
	var err error
	pterm.DisableOutput()
	// pterm.EnableDebugMessages()
	pterm.DisableStyling()
	app := "gofumpt"
	goInstallCmd := "mvdan.cc/gofumpt@latest"
	want = filepath.Join(req.GetGoPath(), "bin", app)
	pterm.Debug.Printfln("want: filepath: %s", want)
	got, err = req.ResolveBinaryByInstall(app, goInstallCmd)
	is.NoErr(err)       // ResolveBinaryByInstall should not error
	is.Equal(want, got) // Binary path should be returned

	_ = sh.Rm(want) // Remove the gofumpt binary to ensure we can test failure case
	got, err = req.ResolveBinaryByInstall(app, goInstallCmd)
	is.NoErr(err)       // ResolveBinaryByInstall should not error after reinstalling
	is.Equal(want, got) // Binary path should be returned after installation
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
