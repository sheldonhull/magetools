//go:build examples

// Overcomplicated example of installing python 3 with built in commands.
// Probably a much better way to do this with other packages, but I wanted to save this as one example.
// Additionally, I've reverted to using asdf or other tooling instead.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
)

// installPython3 is a setup helper to ensure python3 is setup on new systems for pre-commit tooling.
//
// Yes this should be done in Docker, or in OS, but in this case I'm trying to simplify bootstrap of a project and automatically handle pre-commit setup in case the environment doesn't have that already setup.
func installPython3() (string, error) { //nolint: cyclop // cyclop: acceptable as this is a cross platform installation function that tries to help setup a system with python if doesn't exist.
	switch runtime.GOOS {
	case "linux":
		pterm.Info.Println("attempting to sudo install python3-pip")
		_, err := os.Stat("apt-get")
		if err != nil {
			return "", fmt.Errorf(
				"linux: automatic setup is only supported on debian based systems with apt-get command available: %w",
				err,
			)
		}
		if err := sh.RunV("sudo", "apt-get", "-yqq", "update"); err != nil {
			pterm.Error.Println("running apt-get update")
			return "", fmt.Errorf("linux: failed to apt-get update: %w", err)
		}
		pterm.Warning.Println("This can take 5-10 mins, depending on your system.")
		if err := sh.RunV("sudo", "apt-get", "-yqq", "install", "python3-pip"); err != nil {
			pterm.Error.Println("failed to sudo install python3-pip")
			return "", fmt.Errorf("linux: failed to sudo install python3-pip: %w", err)
		}
	case "darwin":
		pterm.Info.Println("will attempt to install with homebrew")
		_, err := os.Stat("brew")
		if err != nil {
			return "", fmt.Errorf(
				"darwin: automatic setup only works with homebrew installed. Try again after installing homebrew. https://docs.brew.sh/Installation: %w",
				err,
			)
		}
		if err := sh.RunV("brew", "install", "python3@3.10"); err != nil {
			return "", fmt.Errorf("failure on: brew install python3@3.10: %w", err)
		}
	case "windows":
		_, err := os.Stat("scoop")
		if err != nil {
			return "", fmt.Errorf(
				"windows: automatic setup only works with Scoop installed. Try again after installing Scoop. https://github.com/ScoopInstaller/Scoop#installation: %w",
				err,
			)
		}
		if err := sh.RunV("scoop", "install", "python"); err != nil {
			return "", fmt.Errorf("scoop install python3@3.10: %w", err)
		}
	}
	path, err := exec.LookPath("python3")
	return path, err
}
