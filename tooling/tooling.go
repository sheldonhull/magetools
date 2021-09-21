// Package tooling provides common tooling install setup for go linting, formatting, and other Go tools with nice console output for interactive use.
package tooling

import (
	"fmt"
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
