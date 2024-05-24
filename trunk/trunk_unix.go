// go:build linux darwin

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
	magetoolsutils.CheckPtermDebug()
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

	_, err := exec.LookPath("trunk")
	if err != nil && os.IsNotExist(err) {
		pterm.Warning.Printfln(
			"unable to resolve aqua cli tool, please install for automated project tooling setup: https://aquaproj.github.io/docs/tutorial-basics/quick-start#install-aqua",
		)
		_, err := script.Exec("curl https://get.trunk.io -fsSL").Exec("bash -s -- -y").Stdout()
		if err != nil {
			return err
		}
	} else {
		pterm.Success.Printfln("trunk.io already installed, skipping")
	}
	return nil
}

// ⚙️ InstallPlugins ensures the required runtimes are installed.
func (Trunk) InstallPlugins() error {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultSection.Println("(Trunk) InstallPlugins")
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
	trunkArgs := []string{
		"upgrade",
	}
	if ci.IsCI() {
		trunkArgs = append(trunkArgs, "--ci")
	}
	return sh.RunV("trunk", trunkArgs...)
}
