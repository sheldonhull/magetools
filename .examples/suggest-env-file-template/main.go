//go:build examples

/*
This is an example of using mage and promptui to improve the development workflow.

Let's say you need to startup a docker compose setup and you expect an .env file to exist.

- You shouldn't commit .env files to source control.
- It's a manual step to create the .env file.
- It's more manual steps to go do this after getting an error in docker compose.

Use mage to solve this problem by doing a quick check for the missing .env file, and prompting for creating from a copy of an existing environment file.

This would improve the workflow for any new contributors and eliminate any wasted time figuring out what to do again in the future.
*/
package main

import (
	"os"
	"path/filepath"

	"github.com/magefile/mage/sh"
	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"

	"github.com/sheldonhull/magetools/ci"
	mtu "github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

const (
	envdir = "env"

	// You can jump start your local .env file by copying the .env.example or using the .env.local that will contain local dev defaults to start with.
	envFile = ".env"
)

// listEnvFiles lists all the matching files from env dir based on glob pattern.
func listEnvFiles(pattern string) ([]string, error) {
	mtu.CheckPtermDebug()
	pterm.Debug.Printf("searching %q for %q\n", envdir, pattern)
	files, err := filepath.Glob(filepath.Join(envdir, pattern))
	if err != nil {
		return nil, err
	}
	return files, nil
}

// promptForSelect runs promptui to have file selected and response returned.
func promptForSelect(items []string) (string, error) {
	mtu.CheckPtermDebug()
	if len(items) == 0 {
		return "", nil
	}

	prompt := promptui.Select{
		Label: "Do you want to copy one of these to start you off? \nCtrl+C to cancel",
		Items: items,
	}
	_, response, err := prompt.Run()
	if err != nil {
		pterm.Error.Printf("Prompt failed %v\n", err)
		return "", err
	}
	pterm.Info.Printf("You choose [%s]\n", response)
	return response, err
}

// offerToSetupEnvFiles checks for expected file and then if not found offers to copy one matching the target pattern.
// This is to help setup dev environment and catch confusing intialization problems due to missing source files that need to be copied from templates.
func offerToSetupEnvFiles(expectedFile string, pattern string) error {
	mtu.CheckPtermDebug()
	if ci.IsCI() {
		pterm.Debug.Println("skipping offerToSetupEnvFiles() per CI system")
		return nil
	}

	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		pterm.Error.Printf("%q doesn't exist yet\n"+
			"This should be initially copied from example file and customized locally.\n"+
			"The changes you make in it don't get committed into source control.\nInstead an empty template is versioned.\n", expectedFile)
	} else {
		pterm.Success.Printf("✅ %q is setup\n", expectedFile)
		return nil
	}

	items, err := listEnvFiles(pattern)
	if err != nil {
		pterm.Error.Println("unable to match pattern in env dir")
		return err
	}
	// append none to the beginning to ensure it's quick to select
	items = append([]string{"none"}, items...)
	file, err := promptForSelect(items)
	if err != nil {
		return err
	}
	if file != "none" {
		err := sh.Copy(expectedFile, file)
		if err != nil {
			return err
		}
	}
	return nil
}
