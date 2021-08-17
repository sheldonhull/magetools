package gotools_test

import (
	"os"
	"testing"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/gotools"
)

func TestGolang_InitAndClean(t *testing.T) {
	is := iz.New(t)
	pterm.DisableColor()
	pterm.DisableStyling()

	defer func() {
		err := os.RemoveAll("_tools")
		is.NoErr(err) // Clean should not fail
	}()

	err := gotools.Golang{}.Init()
	is.NoErr(err) // Init should not fail

	files, err := os.ReadDir("_tools")
	counter := 0
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		counter++
	}
	is.NoErr(err)        // ReadDir for tools
	is.Equal(counter, 7) // should have 7 tools installed

	err = gotools.Golang{}.Clean()
	is.NoErr(err) // Clean should not fail
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