package gittools_test

import (
	"testing"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"

	"github.com/sheldonhull/magetools/gittools"
)

func TestGolang_InitAndClean(t *testing.T) {
	is := iz.New(t)
	pterm.DisableColor()
	pterm.DisableStyling()

	err := gittools.Gittools{}.Init()
	is.NoErr(err) // Init should not fail
}
