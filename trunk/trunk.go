// Package trunk contains tasks to automate setup and usage in Mage based projects for the great trunk.io tooling.
package trunk

import (
	"errors"
	"os"
	"os/exec"
	"runtime"

	"github.com/bitfield/script"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

// ⚙️ Trunk contains tasks related to Trunk.io tooling.
type Trunk mg.Namespace

// ⚙️ Init installs trunk and ensures the plugins are setup.
func (Trunk) Init() {
	mg.SerialDeps(
		Trunk{}.Install,
		Trunk{}.InstallPlugins,
	)
}

// checkCompatibility checks the current OS and warns if trunk.io is not supported with an error that can be either bypassed or handled.
func checkCompatibility() error {
	magetoolsutils.CheckPtermDebug()
	if runtime.GOOS == "windows" {
		pterm.Warning.Printfln("trunk.io is not supported on: %s", runtime.GOOS)
		return errors.New("trunk.io is not supported on: " + runtime.GOOS)
	}
	return nil
}

// ⚙️ InstallTrunk installs trunk.io tooling if it isn't already found.
func (Trunk) Install() error {
	magetoolsutils.CheckPtermDebug()
	if err := checkCompatibility(); err != nil {
		pterm.Warning.Printfln("[non-terminating] skipping trunk.io installation due to incompatible OS: %v", err)
		return nil
	}
	_, err := exec.LookPath("trunk")
	if err != nil && os.IsNotExist(err) {
		pterm.Warning.Printfln("unable to resolve aqua cli tool, please install for automated project tooling setup: https://aquaproj.github.io/docs/tutorial-basics/quick-start#install-aqua")
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
	if err := checkCompatibility(); err != nil {
		pterm.Warning.Printfln("[non-terminating] skipping trunk.io installation due to incompatible OS: %v", err)
		return nil
	}
	return sh.RunV("trunk", "install")
}

// ⚙️ Upgrade upgrades trunk using itself and also the plugins.
func (Trunk) Upgrade() error {
	if err := checkCompatibility(); err != nil {
		pterm.Warning.Printfln("[non-terminating] skipping trunk.io installation due to incompatible OS: %v", err)
		return nil
	}
	return sh.RunV("trunk", "install")
}
