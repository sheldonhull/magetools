package licensing_test

import (
// "os"
// "strings"
// "testing"

// "github.com/magefile/mage/sh"
// iz "github.com/matryer/is"
// "github.com/pterm/pterm"
// "github.com/sheldonhull/magetools/licensing"
)

// func TestInitAndSave(t *testing.T) {
// 	if !strings.Contains(strings.ToLower(os.Getenv("GOTESTS")), "slow") {
// 		t.Skip("GOTESTS should include 'slow' to run this test")
// 	}
// 	pterm.DisableStyling()

// 	is := iz.New(t)
// 	var err error

// 	defer func() {
// 		_ = sh.Rm("_tools")
// 		_ = sh.Rm(".licenses")
// 	}()

// 	err = licensing.Licensing{}.Init()
// 	is.NoErr(err) // should not error on Init

// 	_ = licensing.Licensing{}.Save()
// 	// This step might error if there are non-go dependencies so it's ignored for now
// 	// is.NoErr(err) // should not error on Save.

// 	err = licensing.Licensing{}.Check()
// 	is.NoErr(err) // should not error on Checking for forbidden licenses

// 	err = licensing.Licensing{}.CSV()
// 	is.NoErr(err) // should not error on running CSV check

// 	_, err = os.Stat(".licenses")
// 	is.NoErr(err) // should not error on finding licenses.csv
// }
