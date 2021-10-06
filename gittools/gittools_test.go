package gittools_test

import (
	"os"
	"strings"
	"testing"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/gittools"
)

func TestGolang_InitAndClean(t *testing.T) {
	if !strings.Contains(strings.ToLower(os.Getenv("GOTESTS")), "slow") {
		t.Skip("GOTESTS should include 'slow' to run this test")
	}
	is := iz.New(t)
	pterm.DisableColor()
	pterm.DisableStyling()

	err := gittools.Gittools{}.Init()
	is.NoErr(err) // Init should not fail
}
