// Package asdf installs asdf packages and makes sure asdf is installed for you.
//
// Requires Linux or MacOS.
package asdf

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/bitfield/script"
	"github.com/dustin/go-humanize"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

// Asdf contains tasks for asdf.
type Asdf mg.Namespace

func checklinux() {
	if runtime.GOOS != "Linux" {
		_ = mg.Fatalf(1, "this command is only supported on Linux and you are on: %s", runtime.GOOS)
	}
}

// relTime returns just a simple relative time humanized, without the "ago" suffix.
func relTime(t time.Time) string {
	return strings.ReplaceAll(humanize.Time(t), " ago", "")
}

// ⚙️ Install will attempt to install Asdf toolchain (not the actual plugins).
//
// NOTE: 👀 Will try to modify with native Git package commands later, just a quick draft.
func (Asdf) Install() error {
	checklinux()
	magetoolsutils.CheckPtermDebug()
	start := time.Now()
	defer func(start time.Time) {
		pterm.Success.Printf("✅ (Install) Asdf() [%s]\n", relTime(start))
	}(start)
	pterm.DefaultHeader.Println("⚙️ Direnv")

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	p := script.Exec(
		"bash -c 'git clone https://github.com/asdf-vm/asdf.git " + userHomeDir + "/.asdf --branch v0.9.0'",
	)
	output, err := p.String()

	if p.ExitStatus() == 128 { //nolint:gomnd // gomnd: exit code is specific to git
		pterm.Warning.Println("likely already cloned this repo since exit code 128")
	} else {
		pterm.Error.PrintOnError(err)
	}
	pterm.Info.Println(output)

	script.Exec("bash -c 'source " + userHomeDir + "/.asdf/.asdf/asdf.sh'")
	pterm.Warning.Printfln(
		"You need to source this as part of your .zshrc file in your profile if you haven't done this already.\n\nsource \"$HOME/.asdf/asdf.sh\"\n\n",
	)
	return nil
}

// 💾 Apps will install Asdf plugins and the apps based on .tool-versions.
func (Asdf) Apps() error {
	checklinux()
	magetoolsutils.CheckPtermDebug()

	asdfBinary, err := exec.LookPath("asdf")
	if err != nil {
		pterm.Error.Printfln("asdf binary not found. have you sourced this? " +
			"You need to source this as part of your .zshrc file in your profile if you haven't done this already.\n\nsource \"$HOME/.asdf/asdf.sh\"\n\n",
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
			pterm.Error.Printfln("[asdf] failure to add plugin: %v: %q err: %v", pluginName, out, err)
			continue
		}
		pterm.Success.Printfln("[asdf] ✅ %q added [version for install: %s]", pluginName, pluginVersion)
	}
	if err := sh.RunV(asdfBinary, "plugin", "update", "--all"); err != nil {
		pterm.Error.Printfln("[asdf] plugin update failure: %v", err)
		return err
	}
	pterm.Success.Println("[asdf] plugin update successful")
	if err := sh.RunV(asdfBinary, "install"); err != nil {
		pterm.Error.Printfln("[asdf] install failure: %v", err)
		return err
	}
	pterm.Success.Println("[asdf] plugin install successful")
	return nil
}
