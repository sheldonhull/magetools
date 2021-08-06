// ci helps identify when a task is running in a ci context and not interactively
// Currently this supports checking:
// 1. Azure DevOps
// 2. GitHub Actions
//
// Any calling packages can just run `isci := ci.IsCI()`
package ci

import (
	"os"

	"github.com/pterm/pterm"
)

// IsCI will set the global variable for IsCI based on lookup of the environment variable.
func IsCI() bool {
	_, exists := os.LookupEnv("AGENT_ID")
	if exists {
		pterm.Info.Println("Azure DevOps match based on AGENT_ID. Setting IS_CI = 1")

		return true
	}
	_, exists = os.LookupEnv("CI")
	if exists {
		pterm.Info.Println("GitHub actions match based on [CI] env variable. Setting IS_CI = 1")

		return true
	}

	return false
}
