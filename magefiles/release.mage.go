package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

// Changelog contains tasks related to generating release notes with a tool like changie.
type Changelog mg.Namespace

// configCILogger helps add error prefix output for Azure DevOps CI if it detects it's in this context.
func configCILogger() {
	pterm.Error = *pterm.Error.WithShowLineNumber().WithLineNumberOffset(1) //nolint:reassign // changing prefix later, not an issue.
	pterm.Error.Prefix = pterm.Prefix{Text: "##[error]", Style: &pterm.Style{}}
	pterm.Debug.Prefix = pterm.Prefix{Text: "##[debug]", Style: &pterm.Style{}}
	pterm.Warning.Prefix = pterm.Prefix{Text: "##[warning]", Style: &pterm.Style{}}

	pterm.Debug.Printfln("configureLogging() success")
}

// getVersion returns the version and path for the changefile to use for the semver and release notes.
func getVersion() (releaseVersion, cleanPath string, err error) { //nolint:unparam // this can end up being needed for other tasks, so leaving it in for now
	magetoolsutils.CheckPtermDebug()
	configCILogger()
	releaseVersion, err = sh.Output("changie", "latest")
	if err != nil {
		pterm.Error.Printfln("changie pulling latest release note version failure: %v", err)
		return "", "", err
	}
	cleanVersion := strings.TrimSpace(releaseVersion)
	cleanPath = filepath.Join(".changes", cleanVersion+".md")
	if os.Getenv("GITHUB_WORKSPACE") != "" {
		cleanPath = filepath.Join(os.Getenv("GITHUB_WORKSPACE"), ".changes", cleanVersion+".md")
	}
	return cleanVersion, cleanPath, nil
}

// 📦 Bump the application as an interactive command, merging changelog, and running format and git add.
func (Changelog) Bump() error {
	magetoolsutils.CheckPtermDebug()
	configCILogger()
	pterm.DefaultSection.Println("(Changelog) Bump()")
	if err := sh.RunV("changie", "batch", "auto"); err != nil {
		pterm.Warning.Printf("changie batch failure (non-terminating as might be repeating batch command): %v", err)
	}
	if err := sh.RunV("changie", "merge"); err != nil {
		return err
	}
	if err := sh.RunV("trunk", "fmt"); err != nil {
		return err
	}
	if err := sh.RunV("trunk", "check", "--ci"); err != nil {
		pterm.Warning.Printfln(
			"trunk check failure. This is non-terminating for the mage task, but you should check it before merging",
		)
	}
	if err := sh.RunV("git", "add", ".changes/*"); err != nil {
		return err
	}
	if err := sh.RunV("git", "add", "CHANGELOG.md"); err != nil {
		return err
	}

	releaseVersion, _, err := getVersion()
	if err != nil {
		return err
	}
	pterm.Info.Println(" Are you ready to create a commit with these changes?")
	confirm, err := pterm.DefaultInteractiveConfirm.
		WithDefaultValue(false).
		WithRejectText("no").
		WithConfirmText("yes").
		WithDefaultValue(false).Show()
	if err != nil {
		return err
	}
	if !confirm {
		pterm.Warning.Println("someone changed their mind")
		return nil
	}
	response, err := pterm.DefaultInteractiveTextInput.
		WithMultiLine(true).
		WithDefaultText(fmt.Sprintf("feat: 🚀 create release %s", releaseVersion)).Show()
	if err != nil {
		return err
	}
	if err := sh.RunV("git", "commit", "-m", response); err != nil {
		return err
	}
	return nil
}

// 📦 Merge updates the changelog without bumping the version.
// This is useful for when you are picking up after the changie batch has already completed, but need to re-run the changie merge.
func (Changelog) Merge() error {
	magetoolsutils.CheckPtermDebug()
	configCILogger()
	pterm.DefaultSection.Println("(Changelog) Merge()")
	if err := sh.RunV("changie", "merge"); err != nil {
		return err
	}
	if err := sh.RunV("trunk", "fmt"); err != nil {
		return err
	}
	if err := sh.RunV("trunk", "check", "--ci"); err != nil {
		pterm.Warning.Printfln(
			"trunk check failure. This is non-terminating for the mage task, but you should check it before merging",
		)
	}
	if err := sh.RunV("git", "add", ".changes/*"); err != nil {
		return err
	}
	if err := sh.RunV("git", "add", "CHANGELOG.md"); err != nil {
		return err
	}
	releaseVersion, _, err := getVersion()
	if err != nil {
		return err
	}
	pterm.Info.Println(" Are you ready to create a commit with these changes?")
	confirm, err := pterm.DefaultInteractiveConfirm.
		WithDefaultValue(false).
		WithRejectText("no").
		WithConfirmText("yes").
		WithDefaultValue(false).Show()
	if err != nil {
		return err
	}
	if !confirm {
		pterm.Warning.Println("someone changed their mind")
		return nil
	}
	response, err := pterm.DefaultInteractiveTextInput.
		WithMultiLine(true).
		WithDefaultText(fmt.Sprintf("feat: 🚀 create release %s", releaseVersion)).Show()
	if err != nil {
		return err
	}
	if err := sh.RunV("git", "commit", "-m", response); err != nil {
		return err
	}
	return nil
}
