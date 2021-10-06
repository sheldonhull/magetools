// loghelper provides simple helper functions for enabling or disabling more logging with Pterm.

package magetoolsutils_test

import (
	"os"
	"testing"

	u "github.com/sheldonhull/magetools/pkg/magetoolsutils"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"
)

func TestFuncName(t *testing.T) {
	pterm.DisableStyling()
	pterm.DisableOutput()

	orig := os.Getenv("DEBUG")
	defer os.Setenv("DEBUG", orig)

	t.Run("DEBUG is set to 0", func(t *testing.T) {
		is := iz.New(t)
		os.Setenv("DEBUG", "0")
		u.CheckPtermDebug()
		is.Equal(pterm.PrintDebugMessages, false) // DEBUG = 0 should be false
		os.Unsetenv("DEBUG")
	})
	t.Run("DEBUG is set to 1", func(t *testing.T) {
		is := iz.New(t)
		os.Setenv("DEBUG", "1")
		u.CheckPtermDebug()
		is.Equal(pterm.PrintDebugMessages, true) // DEBUG = 1 should be true
		os.Unsetenv("DEBUG")
	})
	t.Run("DEBUG is set to empty string", func(t *testing.T) {
		is := iz.New(t)
		os.Setenv("DEBUG", "")
		u.CheckPtermDebug()
		is.Equal(pterm.PrintDebugMessages, false) // DEBUG = "" should be false
		os.Unsetenv("DEBUG")
	})
	t.Run("DEBUG is unset", func(t *testing.T) {
		is := iz.New(t)
		os.Unsetenv("DEBUG")
		u.CheckPtermDebug()
		is.Equal(pterm.PrintDebugMessages, false) // DEBUG unset should be false
		os.Unsetenv("DEBUG")
	})
}
