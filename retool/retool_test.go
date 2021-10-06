package retool_test

import (
	"os"
	"strings"
	"testing"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"
	tooling "github.com/sheldonhull/magetools/retool"
)

func TestGolang_InitAndClean(t *testing.T) {
	if !strings.Contains(strings.ToLower(os.Getenv("GOTESTS")), "slow") {
		t.Skip("GOTESTS should include 'slow' to run this test")
	}
	is := iz.New(t)
	pterm.DisableStyling()

	defer func() {
		err := os.RemoveAll("_tools")
		is.NoErr(err) // Clean should not fail
	}()

	toolList := []string{"github.com/dustinkirkland/golang-petname/cmd/petname@master"}
	err := tooling.InstallTools(toolList)
	is.NoErr(err) // installing a tool should not fail
}
