//go:build examples

package main

import (
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"

	mtu "github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

// Build namespace contains tasks specific to running build and CI actions.
type Build mg.Namespace

// ✔️ CICheckout ensures the CI system is able to push updates by using a full checkout instead of detached. Requires `branch` parameter input.
//
// Without this, most build systems checkout in a detached manner, meaning the CI system can't push tags, fixes, or anything else.
//
// If you want the CI system to be the only one that can do tagging for example, you'd want to ensure a full checkout to allow any semver history to be parsed and checkout in full, not in detached state.
func (Build) CICheckout(branch string) error {
	mtu.CheckPtermDebug()
	pterm.Info.Printfln("(Build) CICheckout: branch: %q", branch)

	// Use log to set the CI system commit based on the last person.
	username, err := sh.Output("git", "log", "-1", "--pretty=format:\"%an\"")
	if err != nil {
		return err
	}
	email, err := sh.Output("git", "log", "-1", "--pretty=format:\"%ae\"")
	if err != nil {
		return err
	}

	// In local context only.
	if err := sh.RunV("git", "config", "user.name", strings.TrimSpace(username)); err != nil {
		return err
	}
	if err := sh.RunV("git", "config", "user.email", strings.TrimSpace(email)); err != nil {
		return err
	}
	if err := sh.RunV("git", "fetch", "origin"); err != nil {
		return err
	}
	// In some systems, like Azure Devops you get the branch name prepended with refs/heads.
	// This won't work when doing the checkout, so trim out that prefix if it exists.
	if err := sh.RunV("git", "checkout", strings.ReplaceAll(branch, "refs/heads/", "")); err != nil {
		return err
	}
	if err := sh.RunV("git", "pull", "--progress"); err != nil {
		return err
	}
	pterm.Success.Println("(Build) CICheckout")
	return nil
}
