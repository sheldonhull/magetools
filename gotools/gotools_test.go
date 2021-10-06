//go:build integration
// +build integration

package gotools_test

import (
	"os"
	"testing"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/gotools"
)

func TestGo_InitAndClean(t *testing.T) {
	if !strings.Contains(strings.ToLower(os.Getenv("GOTESTS")), "slow") {
		t.Skip("GOTESTS should include 'slow' to run this test")
	}
	is := iz.New(t)
	pterm.DisableColor()
	pterm.DisableStyling()

	err := gotools.Go{}.Init()
	is.NoErr(err) // Init should not fail
}

func TestGo_Tidy(t *testing.T) {
	is := iz.New(t)
	pterm.DisableOutput()
	err := gotools.Go{}.Tidy()
	is.NoErr(err) // Tidy should not fail
}

func TestGo_Fmt(t *testing.T) {
	is := iz.New(t)
	pterm.DisableOutput()
	err := gotools.Go{}.Fmt()
	is.NoErr(err) // Fmt should not fail
}

func TestGo_Lint(t *testing.T) {
	is := iz.New(t)
	pterm.DisableOutput()
	err := gotools.Go{}.Lint() // Lint is check-only, not auto-fix
	is.NoErr(err)              // Lint should not fail
}

func TestGo_GetModuleName(t *testing.T) {
	is := iz.New(t)
	// pterm.DisableOutput()
	got := gotools.Go{}.GetModuleName()               // Lint should never fail, as only diagnostic info
	is.Equal(got, "github.com/sheldonhull/magetools") // GetModuleName should return the module name
}

func TestGo_Doctor(t *testing.T) {
	pterm.DisableOutput()
	gotools.Go{}.Doctor() // Lint should never fail, as only diagnostic info
}

func TestGo_Test(t *testing.T) {
	if !strings.Contains(strings.ToLower(os.Getenv("GOTESTS")), "slow") {
		t.Skip("GOTESTS should include 'slow' to run this test")
	}
	pterm.DisableOutput()

	t.Run("test with no flags", func(t *testing.T) {
		is := iz.New(t)
		err := gotools.Go{}.Test()
		is.NoErr(err) // No error should be returned, and gotools tests shouldn't run by default
		// TODO: mock test invocation to confirm parses, but doesn't rerun each time.
	})
	t.Run("test with GOTEST_FLAGS with -tags=integration", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}
		is := iz.New(t)
		os.Setenv("GOTEST_FLAGS", "-tags=integration")
		err := gotools.Go{}.Test()
		is.NoErr(err) // No error should be returned against --tags=integration
		// TODO: mock test invocation to confirm parses, but doesn't rerun each time.
	})
}
