package secrets

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
	"github.com/sheldonhull/magetools/pkg/req"
)

const (
	// _reportFormat contains the format this package outputs the report in.
	_reportFormat = "json"

	//  _artifactFilePath is the default path to output the artifact out to.
	//  This can be overridden by setting `GITLEAK_ARTIFACT_PATH`.
	_artifactFilePath = ".artifacts/gitleaks.json"

	// _fullPermissions is the permission for full access on the artifact directory.
	_fullPermissions = 0o755
)

// Secrets contains tasks related to checking for sensitive values using tools like Gitleaks.
type Secrets mg.Namespace

// relTime returns just a simple relative time humanized, without the "ago" suffix.
func relTime(t time.Time) string {
	return strings.ReplaceAll(humanize.Time(t), " ago", "")
}

// setupGitLeaks ensures the tool exists and returns the fully qualified path.
func setupGitleaks() (gitleaks string, err error) {
	appgitleaks := "gitleaks"
	gitleaks, err = req.ResolveBinaryByInstall(appgitleaks, "github.com/zricethezav/gitleaks/v8@latest")
	if err != nil {
		pterm.Error.WithShowLineNumber(true).WithLineNumberOffset(1).Printfln("unable to find %s: %v", gitleaks, err)
		return "", err
	}
	return gitleaks, nil
}

// setArtifactPath sets the artifact output to a default of `.artifacts/gitleaks.json` but allows override with env variable of `GITLEAK_ARTIFACT_PATH`.`.
func setArtifactPath() string {
	artifactFilePath, isSet := os.LookupEnv("GITLEAK_ARTIFACT_PATH")
	if !isSet {
		pterm.Debug.Printfln("GITLEAK_ARTIFACT_PATH env var override artifactOut to: %q", artifactFilePath)
		artifactFilePath = _artifactFilePath
	}

	// check for the directory to actually exist and create if it doesn't to avoid failures.
	basePath := filepath.Dir(artifactFilePath)
	pterm.Debug.Printfln("basePath: %q", basePath)
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath, _fullPermissions); err != nil {
			pterm.Warning.Printfln("Error creating artifact directory: %v", err)
		}
	}

	return artifactFilePath
}

// 🔐 Detect scans for secret violations with gitleaks without git consideration.
//
// This outputs by default to `.artifacts/gitleaks.json` but can be overridden by setting `GITLEAK_ARTIFACT_PATH`.
//
// The defaults for this scan with `--no-git` to focus on file content without history.
func (Secrets) Detect() error {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultSection.Println("🔍 (Secrets) Check()")
	start := time.Now()
	var err error
	defer func(start time.Time) {
		if err == nil {
			pterm.Success.Printf("✅ (Secrets) Check() [%s]\n", relTime(start))
		}
	}(start)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	gitleaks, err := setupGitleaks()
	artifactFilePath := setArtifactPath()
	if err != nil {
		return err
	}
	pterm.Debug.Printfln("artifactOut: %q", artifactFilePath)
	return sh.Run(
		gitleaks,
		"detect",
		"--report-format",
		"json",
		"--source="+wd,
		"--report-path",
		artifactFilePath,
		"--redact",
		"--no-git",
		"--log-level",
		"warn",
		"--verbose",
	)
}

// 🔐 Protect scans the staged artifacts for violations.
//
// This is useful in pre-commit checks.
//
// This outputs by default to `.artifacts/gitleaks.json` but can be overridden by setting `GITLEAK_ARTIFACT_PATH`.
func (Secrets) Protect() error {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultSection.Println("🔍 (Secrets) Check()")
	start := time.Now()
	var err error
	defer func(start time.Time) {
		if err == nil {
			pterm.Success.Printf("✅ (Secrets) Check() [%s]\n", relTime(start))
		}
	}(start)
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}
	gitleaks, err := setupGitleaks()
	artifactFilePath := setArtifactPath()

	if err != nil {
		return err
	}

	pterm.Debug.Printfln("artifactOut: %q", artifactFilePath)
	return sh.Run(
		gitleaks,
		"protect",
		"--report-format",
		_reportFormat,
		"--source="+workingDir,
		"--report-path",
		artifactFilePath,
		"--redact",
		"--staged",
		"--log-level",
		"warn",
		"--verbose",
	)
}
