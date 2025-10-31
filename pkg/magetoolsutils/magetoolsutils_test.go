// loghelper provides simple helper functions for enabling or disabling more logging with Pterm.

package magetoolsutils_test

import (
	"os"
	"testing"

	u "github.com/sheldonhull/magetools/pkg/magetoolsutils"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"
)

func Test_CheckPtermDebug(t *testing.T) {
	orig := os.Getenv("DEBUG")
	defer os.Setenv("DEBUG", orig)

	origSystem := os.Getenv("SYSTEM_DEBUG")
	defer os.Setenv("SYSTEM_DEBUG", origSystem)

	origActions := os.Getenv("ACTIONS_STEP_DEBUG")
	defer os.Setenv("ACTIONS_STEP_DEBUG", origActions)

	origMage := os.Getenv("MAGEFILE_VERBOSE")
	defer os.Setenv("MAGEFILE_VERBOSE", origMage)

	testCases := []struct {
		desc   string
		envvar string
	}{
		{
			desc:   "General debug env",
			envvar: "DEBUG",
		},
		{
			desc:   "Azure DevOps debug env",
			envvar: "SYSTEM_DEBUG",
		},
		{
			desc:   "GitHub Actions debug env",
			envvar: "ACTIONS_STEP_DEBUG",
		},
		{
			desc:   "Magefile debug env",
			envvar: "MAGEFILE_VERBOSE",
		},
	}
	for _, tt := range testCases {
		tt := tt
		// Cannot be used in parallel tests
		t.Run(tt.envvar+" unset should be false", func(t *testing.T) {
			is := iz.New(t)
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, false) // unset should be false
		})
		t.Run(tt.envvar+" 0 should be false", func(t *testing.T) {
			var err error
			is := iz.New(t)
			err = os.Setenv(tt.envvar, "0")
			is.NoErr(err) // os.SetEnv() should not error
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, false) // 0 should be false
			err = os.Unsetenv(tt.envvar)
			is.NoErr(err) // os.Unsetenv should not error for test cleanup
			pterm.DisableDebugMessages()
		})
		t.Run(tt.envvar+" false should be false", func(t *testing.T) {
			var err error
			is := iz.New(t)
			err = os.Setenv(tt.envvar, "false")
			is.NoErr(err) // os.SetEnv() should not error
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, false) // "false" should be false
			err = os.Unsetenv(tt.envvar)
			is.NoErr(err) // os.Unsetenv should not error for test cleanup
			pterm.DisableDebugMessages()
		})
		t.Run(tt.envvar+" False should be false", func(t *testing.T) {
			var err error
			is := iz.New(t)
			os.Setenv(tt.envvar, orig)
			err = os.Setenv(tt.envvar, "False")
			is.NoErr(err) // os.SetEnv() should not error
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, false) // "False" should be false
			err = os.Unsetenv(tt.envvar)
			is.NoErr(err) // os.Unsetenv should not error for test cleanup
			pterm.DisableDebugMessages()
		})
		t.Run(tt.envvar+" blank should be false", func(t *testing.T) {
			var err error
			is := iz.New(t)
			os.Setenv(tt.envvar, orig)
			err = os.Setenv(tt.envvar, "")
			is.NoErr(err) // os.SetEnv() should not error
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, false) // "" should be false
			err = os.Unsetenv(tt.envvar)
			is.NoErr(err) // os.Unsetenv should not error for test cleanup
			pterm.DisableDebugMessages()
		})
		t.Run(tt.envvar+" 1 should be true", func(t *testing.T) {
			var err error
			is := iz.New(t)
			os.Setenv(tt.envvar, orig)
			err = os.Setenv(tt.envvar, "1")
			is.NoErr(err) // os.SetEnv() should not error
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, true) // "1" should be true
			err = os.Unsetenv(tt.envvar)
			is.NoErr(err) // os.Unsetenv should not error for test cleanup
			pterm.DisableDebugMessages()
		})
		t.Run(tt.envvar+" true should be true", func(t *testing.T) {
			var err error
			is := iz.New(t)
			os.Setenv(tt.envvar, orig)
			err = os.Setenv(tt.envvar, "true")
			is.NoErr(err) // os.SetEnv() should not error
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, true) // "true" should be true
			err = os.Unsetenv(tt.envvar)
			is.NoErr(err) // os.Unsetenv should not error for test cleanup
			pterm.DisableDebugMessages()
		})
		t.Run(tt.envvar+" True should be true", func(t *testing.T) {
			var err error
			is := iz.New(t)
			os.Setenv(tt.envvar, orig)
			err = os.Setenv(tt.envvar, "True")
			is.NoErr(err) // os.SetEnv() should not error
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, true) // "True" should be true
			err = os.Unsetenv(tt.envvar)
			is.NoErr(err) // os.Unsetenv should not error for test cleanup
			pterm.DisableDebugMessages()
		})
	}
}
