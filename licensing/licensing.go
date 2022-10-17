// Licensing package wraps up go-license to run checks against licenses and validate forbidden license types aren't detected.
//
// Deprecated: This is likely to be removed in an upcoming release as not used and problematic to check output.
// I'll validate usage (which shouldn't be anyone) and update major version release if I remove this so it won't break anything for folks.

package licensing

// This set of tasks helps run Google's licensing tool to check for problem licenses, comply with credit, and also saving source when required.

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
	"github.com/sheldonhull/magetools/tooling"
)

type Licensing mg.Namespace

// licenseDir is the directory where the licenses are stored.
const licenseDir = ".licenses"

// golang tools to ensure are locally vendored.
var toolList = []string{ //nolint:gochecknoglobals // ok to be global for tooling setup
	"github.com/google/go-licenses@master",
}

// mkdirPermissions creates sets the permission.
const mkdirPermissions = 0o700

// createDirectories creates the local working directories for build artifacts and tooling.
func createDirectories() error {
	magetoolsutils.CheckPtermDebug()
	for _, dir := range []string{licenseDir} {
		if err := os.MkdirAll(dir, mkdirPermissions); err != nil {
			pterm.Error.Printf("failed to create dir: [%s] with error: %v\n", dir, err)

			return err
		}
		pterm.Success.Printf("✅ [%s] dir created\n", dir)
	}

	return nil
}

// ⚙️  Init runs all required steps to use this package.
func (Licensing) Init() error {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultHeader.Println("Licensing Init()")
	if err := createDirectories(); err != nil {
		return err
	}
	if err := tooling.InstallTools(toolList); err != nil {
		return err
	}
	pterm.Success.Println("✅ Init")

	return nil
}

// 💾 Save checks the licenses of the files in the given repo and saves to a csv.
func (Licensing) Save() error {
	magetoolsutils.CheckPtermDebug()
	pterm.Info.Println("Checks the licenses and persists to local directory")
	cmdArgs := []string{
		"save", "./...",
		"--save_path",
		licenseDir,
		"--force",
	}

	err := sh.Run("go-licenses", cmdArgs...)
	if err != nil {
		pterm.Error.Println(err)

		return err
	}

	pterm.Success.Println("✅ Licensing Check")

	return nil
}

// 💾 CSV checks the licenses of the files in the given repo persists to a license.csv file.
func (Licensing) CSV() error {
	magetoolsutils.CheckPtermDebug()
	pterm.Info.Println("Checks the licenses and persists to local directory")
	err := sh.Run("go-licenses", "csv", "./...")
	if err != nil {
		pterm.Error.Println(err)

		return err
	}
	pterm.Success.Println("✅ Licenses CSV Check")

	return nil
}

// 🔎 Check look for forbidden licenses.
func (Licensing) Check() error {
	magetoolsutils.CheckPtermDebug()
	pterm.Info.Println("look for forbidden licenses")

	err := sh.Run("go-licenses", "check", "./...")
	if err != nil {
		pterm.Error.Println(err)

		return err
	}
	pterm.Success.Println("Check")

	return nil
}
