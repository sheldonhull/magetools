// Provide Go linting, formatting and other basic tooling.
package gotools

import (
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	// "time".

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
	"github.com/sheldonhull/magetools/tooling"
	modfile "golang.org/x/mod/modfile"
)

type (
	Go mg.Namespace
)

// golang tools to ensure are locally vendored.
var toolList = []string{ //nolint:gochecknoglobals // ok to be global for tooling setup
	"github.com/goreleaser/goreleaser@v0.174.1",
	// "golang.org/x/tools/cmd/goimports@master",
	// "github.com/sqs/goreturns@master",
	"github.com/golangci/golangci-lint/cmd/golangci-lint@master",
	"github.com/dustinkirkland/golang-petname/cmd/petname@master",
	"mvdan.cc/gofumpt@latest",
	// "github.com/daixiang0/gci@latest",

	// Additionally to simplify init command adding the commands to install from VSCode go installer
	"golang.org/x/tools/gopls@latest",
	"github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest",
	"github.com/ramya-rao-a/go-outline@latest",
	"github.com/cweill/gotests/gotests@latest",
	"github.com/fatih/gomodifytags@latest",
	"github.com/josharian/impl@latest",
	"github.com/haya14busa/goplay/cmd/goplay@latest",
	"github.com/go-delve/delve/cmd/dlv@latest",
	"github.com/rogpeppe/godef@latest",
	"github.com/mfridman/tparse@latest",   // nice table output after running test
	"github.com/segmentio/golines@latest", // handles nice clean line breaks of long lines
}

// getModuleName returns the name from the module file.
// Original help on this was: https://stackoverflow.com/a/63393712/68698
func (Go) GetModuleName() string {
	magetoolsutils.CheckPtermDebug()
	goModBytes, err := ioutil.ReadFile("go.mod")
	if err != nil {
		pterm.Warning.Println("getModuleName() can't find ./go.mod")
		// Running one more check above the parent directory in case this is running in a test or nested directory for some reason.
		// Only 1 level lookback for now.
		goModBytes, err = ioutil.ReadFile("../go.mod")
		if err != nil {
			pterm.Warning.Println("getModuleName() not able to find ../go.mod")
			return ""
		}
		pterm.Info.Println("found go.mod in ../go.mod")
	}
	modName := modfile.ModulePath(goModBytes)
	return modName
}

// NOTE: this didn't work compared to running with RunV, so I'm commenting out for now.
// golanglint is alias for running golangci-lint.
// var golanglint = sh.RunCmd("golangci-lint") //nolint:gochecknoglobals // ok to be global for tooling setup

// ⚙️  Init runs all required steps to use this package.
func (Go) Init() error {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultHeader.Println("Go Init()")
	if err := tooling.SilentInstallTools(toolList); err != nil {
		return err
	}
	if err := (Go{}.Tidy()); err != nil {
		return err
	}
	pterm.Success.Println("✅  Go Init")
	return nil
}

// 🧪 Run go test. Optional: GOTEST_FLAGS '-tags integration', Or write your own GOTEST env logic.
// Example of checking based on GOTEST style environment variable:
//
// 	if !strings.Contains(strings.ToLower(os.Getenv("GOTESTS")), "slow") {
//		t.Skip("GOTESTS should include 'slow' to run this test")
// }.
func (Go) Test() error {
	magetoolsutils.CheckPtermDebug()
	var vflag string

	if mg.Verbose() {
		vflag = "-v"
	}
	testFlags := os.Getenv("GOTEST_FLAGS")
	if testFlags != "" {
		pterm.Info.Printf("GOTEST_FLAGS provided: %q\n", testFlags)
	}

	pterm.Info.Println("Running go test")
	if err := sh.RunV("go", "test", "./...", "-cover", "-shuffle", "on", "-race", vflag, testFlags); err != nil {
		return err
	}
	pterm.Success.Println("✅ Go Test")
	return nil
}

// 🧪 Run gotestsum.
func (Go) TestSum() error {
	magetoolsutils.CheckPtermDebug()
	var vflag string

	if mg.Verbose() {
		vflag = "-v"
	}
	testFlags := os.Getenv("GOTEST_FLAGS")
	if testFlags != "" {
		pterm.Info.Printf("GOTEST_FLAGS provided: %q\n", testFlags)
	}

	pterm.Info.Println("Running go test")
	if err := sh.RunV("gotestsum", "--format", "pkgname", "--", "-cover", "-shuffle", "on", "-race", vflag, testFlags, "./..."); err != nil {
		return err
	}
	pterm.Success.Println("✅ gotestsum")
	return nil
}

// 🔎  Run golangci-lint without fixing.
func (Go) Lint() error {
	magetoolsutils.CheckPtermDebug()
	// var vflag string

	// // outFormat := "tab"
	// if mg.Verbose() {
	// 	vflag = "-v"
	// }
	pterm.Info.Println("Running golangci-lint")
	if err := sh.RunV("golangci-lint", "run"); err != nil {
		pterm.Error.Println("golangci-lint failure")

		return err
	}
	// pterm.Info.Println("Running golangci-lint")
	// if err := golanglint("run"); err != nil {
	// 	return err
	// }
	pterm.Success.Println("✅ Go Lint")
	return nil
}

// ⚙️ Lint runs golangci-lint tooling using .golangci.yml settings.
// Recommend setting fast: false in your config and allow tool to set.
// Recommend not setting enable-all: true in config to allow cli to call this for linting.
// REMOVED: 20201-09-29 due to issues found in https://github.com/golangci/golangci-lint/issues/1490
// Will manually invoke desired formatters.
// func (Go) Fmt() error {
// 	// var vflag string

// 	// if mg.Verbose() {
// 	// 	vflag = "-v"
// 	// }

// 	pterm.Info.Println("Running golangci-lint formatter")
// 	if err := sh.RunV("golangci-lint", "run", "--fix", "--enable", "gofumpt,gci", "--fast"); err != nil {
// 		// if err := golanglint("run", "--fix", "--presets", "format", "--fast"); err != nil {
// 		pterm.Error.Println("golangci-lint failure")

// 		return err
// 	}
// 	// if err := golanglint("run", "--fix", vflag); err != nil {
// 	// 	return err
// 	// }
// 	return nil
// }

// ✨ Fmt runs gofumpt. Export SKIP_GOLINES=1 to skip golines.
// Important. Make sure golangci-lint config disables gci, goimports, and gofmt.
// This will perform all the sorting and other linters can cause conflicts in import ordering.
func (Go) Fmt() error {
	magetoolsutils.CheckPtermDebug()

	if err := AddGoPkgBinToPath(); err != nil {
		return err
	}
	gfpath, err := QualifyGoBinary("gofumpt")
	if err != nil {
		pterm.Error.Printfln("unable to find gofumpt: %v", err)
		return err
	}
	if err := sh.Run(gfpath, "-l", "-w", "."); err != nil {
		return err
	}

	pterm.Success.Println("✅ Go Fmt")
	return nil
}

// GetGoPath returns the GOPATH value.
func GetGoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

// AddGoPkgBinToPath ensures the go/bin directory is available in path for cli tooling.
func AddGoPkgBinToPath() error {
	gopath := GetGoPath()
	goPkgBinPath := filepath.Join(gopath, "bin")
	if !strings.Contains(os.Getenv("PATH"), goPkgBinPath) {
		pterm.Debug.Printf("Adding %q to PATH\n", goPkgBinPath)
		updatedPath := strings.Join([]string{goPkgBinPath, os.Getenv("PATH")}, string(os.PathListSeparator))
		if err := os.Setenv("PATH", updatedPath); err != nil {
			pterm.Error.Printf("Error setting PATH: %v\n", err)
			return err
		}
		pterm.Info.Printf("Updated PATH: %q\n", updatedPath)
	}
	pterm.Debug.Printf("bypassed PATH update as already contained %q\n", goPkgBinPath)
	return nil
}

// QualifyGoBinary provides a fully qualified path for an installed Go binary to avoid path issues.
func QualifyGoBinary(binary string) (string, error) {
	gopath := GetGoPath()

	qualifiedPath := filepath.Join(gopath, "bin", binary)
	if _, err := os.Stat(qualifiedPath); err != nil {
		pterm.Error.Printfln("%q not found in bin", binary)
		return "", err
	}
	pterm.Debug.Printfln("%q full path: %q", binary, qualifiedPath)
	return qualifiedPath, nil
}

// ✨ Wrap runs golines powered by gofumpt.
func (Go) Wrap() error {
	magetoolsutils.CheckPtermDebug()
	if err := AddGoPkgBinToPath(); err != nil {
		return err
	}
	gopath := GetGoPath()

	gfpath := filepath.Join(gopath, "bin", "gofumpt")
	if _, err := os.Stat(gfpath); err != nil {
		pterm.Error.Printf("gofumpt not found in bin, run mage go:init\n")
		return err
	}
	pterm.Debug.Printf("gofumpt full path: %q\n", gfpath)
	if err := sh.Run(
		"golines",
		".",
		"--base-formatter",
		gfpath,
		"-w",
		"--max-len=120",
		"--reformat-tags"); err != nil {
		return err
	}
	pterm.Success.Println("✅ Go Fmt")
	return nil
}

// 🧹 Tidy tidies.
func (Go) Tidy() error {
	magetoolsutils.CheckPtermDebug()
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return err
	}
	pterm.Success.Println("✅ Go Tidy")
	return nil
}

// 🏥 Doctor will provide config details.
func (Go) Doctor() {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultHeader.Printf("🏥 Doctor Diagnostic Checks\n")
	pterm.DefaultSection.Printf("🏥  Environment Variables\n")

	primary := pterm.NewStyle(pterm.FgLightCyan, pterm.BgGray, pterm.Bold)
	// secondary := pterm.NewStyle(pterm.FgLightGreen, pterm.BgWhite, pterm.Italic)
	if err := pterm.DefaultTable.WithHasHeader().
		WithBoxed(true).
		WithHeaderStyle(primary).
		WithData(pterm.TableData{
			{"Variable", "Value"},
			{"GOVERSION", runtime.Version()},
			{"GOOS", runtime.GOOS},
			{"GOARCH", runtime.GOARCH},
			{"GOROOT", runtime.GOROOT()},
		}).Render(); err != nil {
		pterm.Error.Printf(
			"pterm.DefaultTable.WithHasHeader of variable information failed. Continuing...\n%v",
			err,
		)
	}
	pterm.Success.Println("Doctor Diagnostic Checks")
}

// 🏥  LintConfig will return output of golangci-lint config.
func (Go) LintConfig() error {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultHeader.Println("🏥 LintConfig Diagnostic Checks")
	pterm.DefaultSection.Println("🔍 golangci-lint linters with --fast")
	if err := sh.RunV("golangci-lint", "linters", "--fast"); err != nil {
		pterm.Error.Println("unable to run golangci-lint")
		return err
	}
	pterm.DefaultSection.Println("🔍  golangci-lint linters with plain run")
	if err := sh.RunV("golangci-lint", "linters"); err != nil {
		pterm.Error.Println("unable to run golangci-lint")
		return err
	}

	pterm.Success.Println("LintConfig Diagnostic Checks")
	return nil
}
