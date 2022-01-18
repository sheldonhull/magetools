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
	// pterm.EnableDebugMessages()

	// pterm.DisableDebugMessages()
	// pterm.DisableStyling()
	// pterm.DisableOutput()
	orig := os.Getenv("DEBUG")
	defer os.Setenv("DEBUG", orig)

	origSystem := os.Getenv("SYSTEM_DEBUG")
	defer os.Setenv("SYSTEM_DEBUG", origSystem)

	origActions := os.Getenv("ACTIONS_STEP_DEBUG")
	defer os.Setenv("ACTIONS_STEP_DEBUG", origActions)
	testCases := []struct {
		desc   string
		envvar string
	}{
		{
			desc:   "DEBUG general flag",
			envvar: "DEBUG",
		},
		{
			desc:   "SYSTEM_DEBUG azure-devops flag",
			envvar: "SYSTEM_DEBUG",
		},
		{
			desc:   "ACTIONS_STEP_DEBUG github-actions flag",
			envvar: "ACTIONS_STEP_DEBUG",
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
