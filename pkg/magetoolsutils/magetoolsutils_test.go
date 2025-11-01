// loghelper provides simple helper functions for enabling or disabling more logging with Pterm.

package magetoolsutils_test

import (
	"testing"

	u "github.com/sheldonhull/magetools/pkg/magetoolsutils"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"
)

func Test_CheckPtermDebug(t *testing.T) {
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
			is := iz.New(t)
			t.Setenv(tt.envvar, "0")
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, false) // 0 should be false
			pterm.DisableDebugMessages()
		})
		t.Run(tt.envvar+" false should be false", func(t *testing.T) {
			is := iz.New(t)
			t.Setenv(tt.envvar, "false")
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, false) // "false" should be false
			pterm.DisableDebugMessages()
		})
		t.Run(tt.envvar+" False should be false", func(t *testing.T) {
			is := iz.New(t)
			t.Setenv(tt.envvar, "False")
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, false) // "False" should be false
			pterm.DisableDebugMessages()
		})
		t.Run(tt.envvar+" blank should be false", func(t *testing.T) {
			is := iz.New(t)
			t.Setenv(tt.envvar, "")
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, false) // "" should be false
			pterm.DisableDebugMessages()
		})
		t.Run(tt.envvar+" 1 should be true", func(t *testing.T) {
			is := iz.New(t)
			t.Setenv(tt.envvar, "1")
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, true) // "1" should be true
			pterm.DisableDebugMessages()
		})
		t.Run(tt.envvar+" true should be true", func(t *testing.T) {
			is := iz.New(t)
			t.Setenv(tt.envvar, "true")
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, true) // "true" should be true
			pterm.DisableDebugMessages()
		})
		t.Run(tt.envvar+" True should be true", func(t *testing.T) {
			is := iz.New(t)
			t.Setenv(tt.envvar, "True")
			u.CheckPtermDebug()
			is.Equal(pterm.PrintDebugMessages, true) // "True" should be true
			pterm.DisableDebugMessages()
		})
	}
}
