package adr_test

import (
	"os"
	"testing"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/adr"
)

func TestInitAndGenerate(t *testing.T) {
	pterm.DisableStyling()
	is := iz.New(t)
	var err error

	defer func() {
		_ = os.Remove("adrgen.config.yml")
		err = os.RemoveAll("docs")
		is.NoErr(err) // Clean should not fail
	}()

	err = adr.Adr{}.Init()
	is.NoErr(err) // should not error on Init

	_, err = os.Stat("adrgen.config.yml")
	is.NoErr(err) // adrgen.config.yml should be created

	is.NoErr(err) // should not error on finding docs/godocs
	err = adr.Adr{}.Report()
	is.NoErr(err) // should not error on generating report

	_, err = os.Stat("docs/adr/adr-report.md")
	is.NoErr(err) // docs/godocs/adr-report.md should be found

	_, err = os.Stat("docs/adr")
	is.NoErr(err) // should not error on finding docs/godocs
}
