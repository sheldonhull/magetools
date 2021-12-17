package gotools_test

import (
	"os"
	"strings"
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
	// pterm.DisableOutput()
	err := gotools.Go{}.Fmt()
	is.NoErr(err) // Fmt should not fail with gofumpt
}

func TestGo_Wrap(t *testing.T) {
	is := iz.New(t)
	// pterm.DisableOutput()
	err := gotools.Go{}.Wrap()
	is.NoErr(err) // Fmt should not fail with golines
}

func TestGo_Lint(t *testing.T) {
	is := iz.New(t)
	// pterm.DisableOutput()
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
	pterm.DisableStyling()
	gotools.Go{}.Doctor() // Lint should never fail, as only diagnostic info
}

func TestLintConfig(t *testing.T) {
	is := iz.New(t)
	pterm.DisableOutput()
	pterm.DisableStyling()
	err := gotools.Go{}.LintConfig()
	is.NoErr(err) // Lint config should run without returning any errors
}

func ExampleGo_Test() {
	pterm.DisableOutput()
	if !strings.Contains(strings.ToLower(os.Getenv("GOTESTS")), "superslow") {
		return
		// t.Skip("GOTESTS should include 'slow' to run this test")
	}
	// Running as mage task
	if err := (gotools.Go{}.Test()); err != nil {
		pterm.Error.Printf("ExampleGo_Test: %v\n", err)
	}

	// Running with GOTEST_FLAGS detection
	os.Setenv("GOTEST_FLAGS", "-tags=integration")
	if err := (gotools.Go{}.Test()); err != nil {
		pterm.Error.Printf("ExampleGo_Test: %v\n", err)
	}

	// Output:
}

func ExampleGo_TestSum() {
	pterm.DisableOutput()
	if !strings.Contains(strings.ToLower(os.Getenv("GOTESTS")), "superslow") {
		return
		// t.Skip("GOTESTS should include 'slow' to run this test")
	}
	// Running as mage task
	if err := (gotools.Go{}.TestSum()); err != nil {
		pterm.Error.Printf("ExampleGo_TestSum: %v\n", err)
	}

	// Running with GOTEST_FLAGS detection
	// os.Setenv("GOTEST_FLAGS", "-tags=integration")
	// if err := (gotools.Go{}.TestSum()); err != nil {
	// 	pterm.Error.Printf("ExampleGo_TestSum: %v", err)
	// }

	// Output:
}
