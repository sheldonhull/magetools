// loghelper provides simple helper functions for enabling or disabling more logging with Pterm.
package magetoolsutils

import (
	"os"
	"strconv"

	"github.com/magefile/mage/mg"
	"github.com/pterm/pterm"
)

// CheckPtermDebug looks for DEBUG=1 and sets debug level output if this is found to help troubleshoot tasks.
func CheckPtermDebug() {
	switch {
	case pterm.PrintDebugMessages:
		return // Debug messages already enabled. Do nothing.

	case enableDebugMessagesBasedOnEnv("DEBUG"):
		return // General debug env.

	case enableDebugMessagesBasedOnEnv("SYSTEM_DEBUG"):
		return // Azure DevOps debug env.

	case enableDebugMessagesBasedOnEnv("ACTIONS_STEP_DEBUG"):
		return // GitHub Actions debug env.

	case mg.Verbose():
		pterm.EnableDebugMessages()
		pterm.Debug.Println("mg.Verbose() true (-v or MAGEFILE_VERBOSE env var), setting tasks to debug level output")
		return
	}
}

// enableDebugMessagesBasedOnEnv checks if the specified environment variable is set to a truthy value.
// If so, it enables debug messages in pterm and returns true. Otherwise, it returns false.
func enableDebugMessagesBasedOnEnv(name string) bool {
	envValue, isSet := os.LookupEnv(name)
	if !isSet {
		return false
	}

	debug, err := strconv.ParseBool(envValue)
	if err != nil {
		const offset = 2
		pterm.Warning.WithShowLineNumber(true).
			WithLineNumberOffset(offset).
			Printfln("ParseBool(%s): %v\t debug: %v", name, err, debug)
		return false
	}
	if !debug {
		return false
	}

	pterm.EnableDebugMessages()
	pterm.Debug.Printfln("%s env var detected, setting tasks to debug level output", name)

	return true
}
