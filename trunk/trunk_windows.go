// go:build windows

// Package trunk contains tasks to automate setup and usage in Mage based projects for the great trunk.io tooling.
package trunk

import (
	"os"
	"os/exec"

	"github.com/bitfield/script"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/ci"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

// ⚙️ Trunk contains tasks related to Trunk.io tooling.
type Trunk mg.Namespace

// ⚙️ Init installs trunk and ensures the plugins are setup.
func (Trunk) Init() {
	pterm.DefaultSection.Println("(Trunk) Init")
	mg.SerialDeps(
		Trunk{}.Install,
		Trunk{}.InstallPlugins,
	)
}

// ⚙️ InstallTrunk installs trunk.io tooling if it isn't already found.
func (Trunk) Install() error {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultSection.Println("(Trunk) Install")
	prefixPath()
	// if there's a package.json then use npm install with npm install -D @trunkio/launcher, else install as a global tool
	_, err := exec.LookPath("trunk")
	if err != nil && os.IsNotExist(err) {
		if _, err := os.Stat("package.json"); err == nil {
			pterm.Info.Printfln("Found package.json, installing trunk.io as a dev dependency")
			_, err := script.Exec("npm install -D @trunkio/launcher").Stdout()
			if err != nil {
				return err
			}
		} else {
			pterm.Info.Printfln("No package.json found, installing trunk.io globally")
			_, err := script.Exec("npm install -g @trunkio/launcher").Stdout()
			if err != nil {
				return err
			}
		}
	} else {
		pterm.Success.Printfln("trunk.io already installed, skipping")
	}

	return nil
}

// prefixPath checks for trunk installed as a dev dependency, and if so, prepends the node modules bin/trunk path to the PATH so it can be referenced without npm in front.
func prefixPath() {
	magetoolsutils.CheckPtermDebug()
	if _, err := os.Stat("node_modules/.bin/trunk"); err == nil {
		pterm.Debug.Printfln("Found trunk.io installed as a dev dependency, prefixing PATH")
		os.Setenv("PATH", "node_modules/.bin"+string(os.PathListSeparator)+os.Getenv("PATH"))
	}
}

// ⚙️ InstallPlugins ensures the required runtimes are installed.
func (Trunk) InstallPlugins() error {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultSection.Println("(Trunk) InstallPlugins")
	prefixPath()
	trunkArgs := []string{
		"install",
	}
	if ci.IsCI() {
		trunkArgs = append(trunkArgs, "--ci")
	}
	return sh.RunV("trunk", trunkArgs...)
}

// ⚙️ Upgrade upgrades trunk using itself and also the plugins.
func (Trunk) Upgrade() error {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultSection.Println("(Trunk) Upgrade")
	prefixPath()
	trunkArgs := []string{
		"upgrade",
	}
	if ci.IsCI() {
		trunkArgs = append(trunkArgs, "--ci")
	}
	return sh.RunV("trunk", trunkArgs...)
}
