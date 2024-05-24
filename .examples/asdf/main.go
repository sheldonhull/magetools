//go:build examples

// Package asdf helps configure your current workspace by reading the asdf tools file and installing required plugins and the associated apps.
package asdf

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	script "github.com/bitfield/script"
	humanize "github.com/dustin/go-humanize"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	mtu "github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

func checklinux() {
	if runtime.GOOS != "Linux" {
		_ = mg.Fatalf(1, "this command is only supported on Linux and you are on: %s", runtime.GOOS)
	}
}

// Install contains tasks normally relegated to bash scripts, mostly focused on linux/macOS compatibility.
type Install mg.Namespace

// relTime returns just a simple relative time humanized, without the "ago" suffix.
func relTime(t time.Time) string {
	return strings.ReplaceAll(humanize.Time(t), " ago", "")
}

// ⚙️ Asdf installs tooling.
//
// NOTE: 👀 Will try to modify with native Git package commands later, just a quick draft.
func (Install) Asdf() error {
	checklinux()
	mtu.CheckPtermDebug()
	start := time.Now()
	defer func(start time.Time) {
		pterm.Success.Printf("✅ (Install) Asdf() [%s]\n", relTime(start))
	}(start)
	pterm.DefaultHeader.Println("⚙️ Direnv")

	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	p := script.Exec("bash -c 'git clone https://github.com/asdf-vm/asdf.git " + dirname + "/.asdf --branch v0.9.0'")
	output, err := p.String()

	if p.ExitStatus() == 128 { //nolint:gomnd // gomnd: exit code is specific to git
		pterm.Warning.Println("likely already cloned this repo since exit code 128")
	} else {
		pterm.Error.PrintOnError(err)
	}
	pterm.Info.Println(output)

	script.Exec("bash -c 'source " + dirname + "/.asdf/.asdf/asdf.sh'")
	pterm.Warning.Printfln(
		"You need to add this to your .zshrc.\n\nsource $HOME/.asdf/asdf.sh\n\n",
	)
	return nil
}

// 💾 Apps will install Asdf plugins and the apps based on .tool-versions. Failures on install are warning level only, and failure only occurs when asdf not found at all.
func (Install) Apps() error { //nolint:funlen // It's long and done for bootstrapping asdf apps, ok to leave as the ugly func it is.
	mtu.CheckPtermDebug()

	asdfBinary, err := exec.LookPath("asdf")
	if err != nil {
		pterm.Error.Printfln(
			"asdf binary not found. have you sourced this? " + "You need to add this to your .zshrc.\n\nsource $HOME/.asdf/asdf.sh\n\n",
		)
	}

	f, err := os.Open(".tool-versions")
	if err != nil {
		return err
	}
	// Remember to close the file at the end of the program.
	defer f.Close()

	// Read the file line by line using scanner.
	scanner := bufio.NewScanner(f)
	var pluginName, pluginVersion string
	var issuecount int

	if err := sh.RunV(asdfBinary, "plugin", "update", "--all"); err != nil {
		pterm.Warning.Printfln("[asdf] plugin update failure: %v", err)
	}

	for scanner.Scan() {
		pterm.Debug.Printfln("[asdf] %q", scanner.Text())
		line := strings.Split(scanner.Text(), " ")
		if len(line) != 2 { //nolint:gomnd // gomnd: array should equal this exactly.
			pterm.Warning.Printfln("[asdf] unstable to split: %v", err)
			continue
		}

		pluginName = line[0]
		pluginVersion = line[1]
		p := script.Exec(fmt.Sprintf("bash -c '%s plugin add %s'", asdfBinary, pluginName))
		out, err := p.String()

		if p.ExitStatus() == 2 { //nolint:gomnd // gomnd: valid exit code for already installed.
			pterm.Success.Printfln("[asdf] ✅ %q is already added", pluginName)
			continue
		} else if p.ExitStatus() != 0 {
			issuecount++
			pterm.Warning.Printfln("[asdf] failure to add plugin: %v: %q err: %v", pluginName, out, err)
			continue
		}

		if err := sh.RunV(asdfBinary, "install", pluginName, pluginVersion); err != nil {
			issuecount++
			pterm.Warning.Printfln("[asdf] install failure: %v", err)
			continue
		}
		pterm.Success.Printfln("[asdf] ✅ %q added [version for install: %s]", pluginName, pluginVersion)
	}

	if issuecount > 0 {
		pterm.Warning.Printfln("[asdf] total install issues: %d", issuecount)
		return nil
	}
	pterm.Success.Println("[asdf] plugin install successful")
	return nil
}

// 💾 AsdfDirenvSetup reruns asdf and direnv setup for proper hooks into loading .envrc automatically in your sessions.
//
// Need to remove the legacy direnvrc file to allow this to resetup correctly, so this does a removal of `$HOME/.config/direnv/direnvrc` to allow this to setup correctly.
func (Install) AsdfDirenvSetup() error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	pterm.Info.Println("🔨 Running direnv setup to ensure hooks with asdf are using latest version")
	removefile := filepath.Join(dirname, ".config", "direnv", "direnvrc")

	if err := sh.Rm(removefile); err != nil {
		return fmt.Errorf(
			"Need to remove file: %q to ensure direnv hooks setup correctly, unable to do this, so might need to remove the file manually and rerun mage install:asdfdirenvsetup: %v\n",
			removefile,
			err,
		)
	}
	if err := sh.RunV("asdf", "direnv", "setup", "--shell", "zsh", "--version", "latest"); err != nil {
		pterm.Error.Printfln(
			"If asdf isn't found, then you might need to source this in your terminal first:\n\nsource \"${HOME}\\.asdf\\asdf.sh",
		)
		return fmt.Errorf("unable to run asdf direnv setup: %w", err)
	}
	return nil
}
