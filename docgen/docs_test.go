package docgen_test

import (
	"os"
	"testing"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"

	"github.com/sheldonhull/magetools/docgen"
)

func TestInitAndGenerate(t *testing.T) {
	pterm.DisableStyling()
	is := iz.New(t)
	var err error

	defer func() {
		err = os.RemoveAll("./docs")
		is.NoErr(err) // Clean should not fail
	}()

	err = docgen.Docs{}.Init()
	is.NoErr(err) // should not error on Init

	err = docgen.Docs{}.Generate("github")
	is.NoErr(err) // should not error on generate

	_, err = os.Stat("./docs/godocs")
	is.NoErr(err) // should not error on finding docs/godocs
}
