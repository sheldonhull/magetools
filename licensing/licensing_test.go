package licensing_test

import (
	"os"
	"testing"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"

	"github.com/sheldonhull/magetools/licensing"
)

func TestInitAndSave(t *testing.T) {
	pterm.DisableStyling()
	is := iz.New(t)
	var err error

	defer func() {
		err = os.RemoveAll("_tools")
		is.NoErr(err) // Clean should not fail

		err = os.RemoveAll("licenses")
		is.NoErr(err) // Clean should not fail
	}()

	err = licensing.Licensing{}.Init()
	is.NoErr(err) // should not error on Init

	err = licensing.Licensing{}.Save()
	is.NoErr(err) // should not error on Save

	err = licensing.Licensing{}.Check()
	is.NoErr(err) // should not error on Checking for forbidden licenses

	_, err = os.Stat("licenses")
	is.NoErr(err) // should not error on finding licenses.csv
}
