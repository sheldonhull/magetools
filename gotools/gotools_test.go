package gotools_test

import (
	"testing"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"

	"github.com/sheldonhull/magetools/gotools"
)

func TestGo_InitAndClean(t *testing.T) {
	is := iz.New(t)
	pterm.DisableColor()
	pterm.DisableStyling()

	err := gotools.Go{}.Init()
	is.NoErr(err) // Init should not fail
}

func TestGo_Tidy(t *testing.T) {
	is := iz.New(t)
	pterm.DisableOutput()
	err := gotools.Go{}.Tidy()
	is.NoErr(err) // Tidy should not fail
}

func TestGo_Fmt(t *testing.T) {
	is := iz.New(t)
	pterm.DisableOutput()
	err := gotools.Go{}.Fmt()
	is.NoErr(err) // Fmt should not fail
}

func TestGo_Lint(t *testing.T) {
	is := iz.New(t)
	pterm.DisableOutput()
	err := gotools.Go{}.Lint() // Lint is check-only, not auto-fix
	is.NoErr(err)              // Lint should not fail
}

func TestGo_Doctor(t *testing.T) {
	pterm.DisableOutput()
	gotools.Go{}.Doctor() // Lint should never fail, as only diagnostic info
}
