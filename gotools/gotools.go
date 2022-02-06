// Provide Go linting, formatting and other basic tooling.
//
// Some additional benefits to using this over calling natively are:
//
// - Uses improved gofumpt over gofmt.
//
// - Uses golines with `mage go:wrap` to automatically wrap long expressions.
//
// - If the non-standard tooling isn't installed, it will automatically go install the required tool on calling, reducing the need to run setup processes.
package gotools

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/ci"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
	"github.com/sheldonhull/magetools/pkg/req"
	"github.com/sheldonhull/magetools/tooling"
	"github.com/ztrue/tracerr"
	modfile "golang.org/x/mod/modfile"
)

type (
	Go mg.Namespace
)

const (
	// _maxLength is the maximum length allowed before golines will wrap functional options and similar style calls.
	//
	// For example:
	//
	// log.Str(foo).Str(bar).Str(taco).Msg("foo") if exceeded the length would get transformted into:
	//
	// log.Str(foo).
	//	Str(bar).
	//	Str(taco).
	//	Msg("foo")
	_maxLength = 120
)

// toolList is the list of tools to initially install when running a setup process in a project.
//
// This includes goreleaser, golangci-lint, petname (for random build/titles).
//
// In addition, core tooling from VSCode Install Tool commands are included so using in a Codespace project doesn't require anything other than mage go:init.
var toolList = []string{ //nolint:gochecknoglobals // ok to be global for tooling setup

	// build tools
	"github.com/goreleaser/goreleaser@v0.174.1",
	"github.com/dustinkirkland/golang-petname/cmd/petname@latest",
	"github.com/AlexBeauchemin/gobadge@latest", // create a badge for your markdown from the coverage files.
	// linting tools
	"github.com/golangci/golangci-lint/cmd/golangci-lint@latest",

	// formatting tools
	"github.com/segmentio/golines@latest", // handles nice clean line breaks of long lines
	"mvdan.cc/gofumpt@latest",

	// Testing tools
	"github.com/mfridman/tparse@latest", // nice table output after running test
	"gotest.tools/gotestsum@latest",     // ability to run tests with junit, json output, xml, and more.

	"golang.org/x/tools/gopls@latest",
	"github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest",
	"github.com/ramya-rao-a/go-outline@latest",
	"github.com/cweill/gotests/gotests@latest",
	"github.com/fatih/gomodifytags@latest",
	"github.com/josharian/impl@latest",
	"github.com/haya14busa/goplay/cmd/goplay@latest",
	"github.com/go-delve/delve/cmd/dlv@latest",
	"github.com/rogpeppe/godef@latest",

	// Self setup mage
	"github.com/magefile/mage@latest",
}

// getModuleName returns the name from the module file.
// Original help on this was: https://stackoverflow.com/a/63393712/68698
func (Go) GetModuleName() string {
	magetoolsutils.CheckPtermDebug()
	goModBytes, err := ioutil.ReadFile("go.mod")
	if err != nil {
		pterm.Warning.WithShowLineNumber(true).WithLineNumberOffset(1).Println("getModuleName() can't find ./go.mod")
		// Running one more check above the parent directory in case this is running in a test or nested directory for some reason.
		// Only 1 level lookback for now.
		goModBytes, err = ioutil.ReadFile("../go.mod")
		if err != nil {
			pterm.Warning.WithShowLineNumber(true).
				WithLineNumberOffset(1).
				Println("getModuleName() not able to find ../go.mod")
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
	if err := sh.RunV("go", "test", "./...", "-cover", "-shuffle", "on", "-race", strings.TrimSpace(vflag), strings.TrimSpace(testFlags)); err != nil {
		return err
	}
	pterm.Success.Println("✅ Go Test")
	return nil
}

// 🧪 Run gotestsum.
func (Go) TestSum() error {
	magetoolsutils.CheckPtermDebug()
	appgotestsum := "gotestsum"
	gotestsum, err := req.ResolveBinaryByInstall(appgotestsum, "gotest.tools/gotestsum@latest")
	if err != nil {
		pterm.Error.WithShowLineNumber(true).WithLineNumberOffset(1).Printfln("unable to find %s: %v", gotestsum, err)
		return err
	}

	var vflag string
	if mg.Verbose() {
		vflag = "-v"
	}
	testFlags := os.Getenv("GOTEST_FLAGS")
	if testFlags != "" {
		pterm.Info.Printf("GOTEST_FLAGS provided: %q\n", testFlags)
	}
	// The artifact directory will atttempt to be set to the environment variable: BUILD_ARTIFACTSTAGINGDIRECTORY, but if it isn't set then it will default to .artifacts, which should be excluded in gitignore.
	var artifactDir string
	var ok bool
	artifactDir, ok = os.LookupEnv("BUILD_ARTIFACTSTAGINGDIRECTORY")
	if !ok {
		artifactDir = ".artifacts"
	}
	junitFile := filepath.Join(artifactDir, "junit.xml")
	jsonFile := filepath.Join(artifactDir, "gotest.json")
	coverfile := filepath.Join(artifactDir, "cover.out")
	if err := os.MkdirAll(artifactDir, os.FileMode(0o755)); err != nil { //nolint: gomnd // gomnd, acceptable per permissions
		return err
	}
	format := "dots-v2"
	if ci.IsCI() {
		format = "pkgname"
	}
	pterm.Info.Println("Running go test")
	if err := sh.RunV("gotestsum",
		"--format", format,
		"--junitfile", junitFile,
		"--jsonfile", jsonFile,
		"--",
		"-coverpkg=./...",
		fmt.Sprintf("-coverprofile=%s", coverfile),
		"-covermode", "atomic",
		"-shuffle=on",
		"-race",
		vflag,
		testFlags,
		"./...",
	); err != nil {
		return err
	}
	pterm.Success.Println("✅ gotestsum")
	return nil
}

// 🔎  Run golangci-lint without fixing.
func (Go) Lint() error {
	magetoolsutils.CheckPtermDebug()

	pterm.Info.Println("Running golangci-lint")
	if err := sh.RunV("golangci-lint", "run"); err != nil {
		pterm.Error.WithShowLineNumber(true).WithLineNumberOffset(1).Println("golangci-lint failure")

		return err
	}
	pterm.Success.Println("✅ Go Lint")
	return nil
}

// ✨ Fmt runs gofumpt. Export SKIP_GOLINES=1 to skip golines.
// Important. Make sure golangci-lint config disables gci, goimports, and gofmt.
// This will perform all the sorting and other linters can cause conflicts in import ordering.
func (Go) Fmt() error {
	magetoolsutils.CheckPtermDebug()
	appgofumpt := "gofumpt"
	gofumpt, err := req.ResolveBinaryByInstall(appgofumpt, "mvdan.cc/gofumpt@latest")
	if err != nil {
		pterm.Error.WithShowLineNumber(true).WithLineNumberOffset(1).Printfln("unable to find %s: %v", gofumpt, err)
		return err
	}
	if err := sh.Run(gofumpt, "-l", "-w", "."); err != nil {
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

// ✨ Wrap runs golines powered by gofumpt.
func (Go) Wrap() error {
	magetoolsutils.CheckPtermDebug()
	appgolines := "golines"
	appgofumpt := "gofumpt"
	binary, err := req.ResolveBinaryByInstall(appgolines, "github.com/segmentio/golines@latest")
	if err != nil {
		tracerr.PrintSourceColor(err)
		pterm.Error.WithShowLineNumber(true).WithLineNumberOffset(1).Printfln("unable to find %s: %v", appgolines, err)
		return err
	}
	gofumpt, err := req.ResolveBinaryByInstall(appgofumpt, "mvdan.cc/gofumpt@latest")
	if err != nil {
		tracerr.PrintSourceColor(err)
		pterm.Error.WithShowLineNumber(true).WithLineNumberOffset(1).Printfln("unable to find %s: %v", gofumpt, err)
		return err
	}
	if err := sh.Run(
		binary,
		".",
		"--base-formatter",
		gofumpt,
		"-w",
		fmt.Sprintf("--max-len=%d", _maxLength),
		"--reformat-tags"); err != nil {
		tracerr.PrintSourceColor(err)
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
		tracerr.PrintSourceColor(err)
		pterm.Error.WithShowLineNumber(true).WithLineNumberOffset(1).Printfln(
			"pterm.DefaultTable.WithHasHeader of variable information failed. Continuing...%v",
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
		pterm.Error.WithShowLineNumber(true).WithLineNumberOffset(1).Println("unable to run golangci-lint")
		tracerr.PrintSourceColor(err)
		return err
	}
	pterm.DefaultSection.Println("🔍  golangci-lint linters with plain run")
	if err := sh.RunV("golangci-lint", "linters"); err != nil {
		pterm.Error.WithShowLineNumber(true).WithLineNumberOffset(1).Println("unable to run golangci-lint")
		tracerr.PrintSourceColor(err)
		return err
	}

	pterm.Success.Println("LintConfig Diagnostic Checks")
	return nil
}
