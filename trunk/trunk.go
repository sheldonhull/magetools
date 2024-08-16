// Package trunk contains tasks to automate setup and usage in Mage based projects for the great trunk.io tooling.
package trunk

import (
	"os"
	"os/exec"
	"runtime"

	"al.essio.dev/pkg/shellescape"
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
func (Trunk) Install() (err error) {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultSection.Println("(Trunk) Install")

	switch runtime.GOOS {
	case "linux", "darwin":
		err = trunkInstallLinuxDarwin()
	case "windows":
		err = trunkInstallWindows()
	default:
		pterm.Error.Println("Unsupported OS")
	}
	return err
}

// trunkInstallLinuxDarwin installs trunk.io on linux and darwin using the default curl approach.
func trunkInstallLinuxDarwin() (err error) {
	_, err = exec.LookPath("trunk")
	if err != nil && os.IsNotExist(err) {
		pterm.Warning.Printfln(
			"unable to resolve aqua cli tool, please install for automated project tooling setup: https://aquaproj.github.io/docs/tutorial-basics/quick-start#install-aqua",
		)
		_, err := script.Exec("curl https://get.trunk.io -fsSL").Exec("bash -s -- -y").Stdout()
		if err != nil {
			pterm.Error.Printfln(
				"curl %s",
				shellescape.QuoteCommand([]string{"https://get.trunk.io -fsSL | bash -s -- -y"}),
			)
			return err
		}
	} else {
		pterm.Success.Printfln("trunk.io already installed, skipping")
	}
	return err
}

// trunkInstallWindows installs trunk.io on windows, using npm install method.
func trunkInstallWindows() (err error) {
	prefixNodeModulesBin()

	// if there's a package.json then use npm install with npm install -D @trunkio/launcher, else install as a global tool
	_, err = exec.LookPath("trunk")
	if err != nil && os.IsNotExist(err) {
		if _, err = os.Stat("package.json"); err == nil {
			pterm.Info.Printfln("Found package.json, installing trunk.io as a dev dependency")
			_, err = script.Exec("npm install -D @trunkio/launcher").Stdout()
			if err != nil {
				pterm.Error.Printfln(shellescape.QuoteCommand([]string{"npm install -D @trunkio/launcher"}))
				return err
			}
		} else {
			pterm.Info.Printfln("No package.json found, installing trunk.io globally")
			_, err = script.Exec("npm install --global @trunkio/launcher").Stdout()
			if err != nil {
				pterm.Error.Printfln(shellescape.QuoteCommand([]string{"npm install --global @trunkio/launcher"}))
				return err
			}
		}
	} else {
		pterm.Success.Printfln("trunk.io already installed, skipping")
	}
	if err != nil {
		pterm.Warning.Printfln("if any odd errors, try upgrading node/npm using install or upgrade command")
		pterm.Warning.Printfln(shellescape.QuoteCommand([]string{"winget install --id OpenJS.NodeJS --source winget"}))
	}
	return err
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
	if runtime.GOOS == "windows" {
		prefixNodeModulesBin()
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
	if runtime.GOOS == "windows" {
		prefixNodeModulesBin()
	}
	return sh.RunV("trunk", trunkArgs...)
}

// prefixNodeModulesBin checks for trunk installed as a dev dependency, and if so, prepends the node modules bin/trunk path to the PATH so it can be referenced without npm in front.
func prefixNodeModulesBin() {
	magetoolsutils.CheckPtermDebug()
	if _, err := os.Stat("node_modules/.bin/trunk"); err == nil {
		pterm.Debug.Printfln("Found trunk.io installed as a dev dependency, prefixing PATH")
		os.Setenv("PATH", "node_modules/.bin"+string(os.PathListSeparator)+os.Getenv("PATH"))
	}
}
