// Package tooling provides common tooling install setup for go linting, formatting, and other Go tools with nice console output for interactive use.
package tooling

import (
	"bufio"
	"fmt"
	"os/exec"
	"time"

	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
)

// InstallTools installs tooling for the project in a local directory to avoid polluting global modules.
func InstallTools(tools []string) error {
	start := time.Now()

	pterm.DefaultSection.Println("Installing Tools")
	env := map[string]string{}
	args := []string{"install"}

	// spinnerLiveText, _ := pterm.DefaultSpinner.Start("InstallTools")
	defer func() {
		duration := time.Since(start)
		msg := fmt.Sprintf("tooling installed: %v\n", duration)
		pterm.Success.Println(msg)
		// spinnerLiveText.Success(msg) // Resolve spinner with success message.
	}()

	// as of last time I checked 2021-09, go get/install wasn't noted to be safe to run in parallel, so keeping it simple with just loop
	for i, t := range tools {
		msg := fmt.Sprintf("install [%d] %s]", i, t)
		// spinner, _ := pterm.DefaultSpinner.
		// 	WithSequence("|", "/", "-", "|", "/", "-", "\\").
		// 	WithRemoveWhenDone(true).
		// 	WithText(msg).
		// 	WithShowTimer(true).Start()

		err := sh.RunWith(env, "go", append(args, t)...)
		if err != nil {
			pterm.Warning.Printf("Could not install [%s] per [%v]\n", t, err)
			// spinner.Fail(fmt.Sprintf("Could not install [%s] per [%v]\n", t, err))

			continue
		}
		// spinner.Success(msg)
		pterm.Success.Println(msg)
	}

	return nil
}

// SilentInstallTools reads the stdout and then uses a spinner to show progress.
// This is designed to swallow up a lot of the noise with go install commands.
// Originally found from: https://www.yellowduck.be/posts/reading-command-output-line-by-line/ and modified.
func SilentInstallTools(toolList []string) error {
	// delay := time.Second * 1 // help prevent jitter
	spin, _ := pterm.DefaultSpinner. // WithDelay((delay)).WithRemoveWhenDone(true).
						WithShowTimer(true).
						WithSequence("|", "/", "-", "|", "/", "-", "\\").
						WithText("Installing tools").
						Start()

	for _, item := range toolList {
		cmd := exec.Command("go", "install", item)

		// Get a pipe to read from standard out
		r, _ := cmd.StdoutPipe()

		// Use the same pipe for standard error
		cmd.Stderr = cmd.Stdout

		// Make a new channel which will be used to ensure we get all output
		done := make(chan struct{})

		// Create a scanner which scans r in a line-by-line fashion
		scanner := bufio.NewScanner(r)

		// Use the scanner to scan the output line by line and log it
		// It's running in a goroutine so that it doesn't block
		go func(item string) {
			// Read line by line and process it

			for scanner.Scan() {
				line := scanner.Text()
				spin.UpdateText(line)
			}
			// We're all done, unblock the channel
			done <- struct{}{}
		}(item)

		// Start the command and check for errors
		err := cmd.Start()
		if err != nil {
			spin.Fail(err)
			_ = spin.Stop()
			return err
		}

		// Wait for all output to be processed
		<-done

		// Wait for the command to finish
		err = cmd.Wait()
		if err != nil {
			spin.Fail(err)
			_ = spin.Stop()
			return err
		}
		spin.Success(item)
	}

	return nil
}
