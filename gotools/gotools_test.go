package gotools_test

import (
	"testing"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/gotools"
)

func TestGolang_InitAndClean(t *testing.T) {
	is := iz.New(t)
	pterm.DisableColor()
	pterm.DisableStyling()

	err := gotools.Golang{}.Init()
	is.NoErr(err) // Init should not fail
}

func TestGolang_Tidy(t *testing.T) {
	is := iz.New(t)
	pterm.DisableOutput()
	err := gotools.Golang{}.Tidy()
	is.NoErr(err) // Tidy should not fail
}

func TestGolang_Fmt(t *testing.T) {
	is := iz.New(t)
	pterm.DisableOutput()
	err := gotools.Golang{}.Fmt()
	is.NoErr(err) // Fmt should not fail
}

func TestGolang_Lint(t *testing.T) {
	is := iz.New(t)
	pterm.DisableOutput()
	err := gotools.Golang{}.Lint() // Lint is check-only, not auto-fix
	is.NoErr(err)                  // Lint should not fail
}
