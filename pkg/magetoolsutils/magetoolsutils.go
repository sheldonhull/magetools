// loghelper provides simple helper functions for enabling or disabling more logging with Pterm.
package magetoolsutils

import (
	"os"

	"github.com/pterm/pterm"
)

// checkPtermDebug looks for DEBUG=1 and sets debug level output if this is found to help troubleshoot tasks.
func CheckPtermDebug() {
	if os.Getenv("DEBUG") == "1" {
		pterm.EnableDebugMessages()
		pterm.Debug.Println("DEBUG enabled per env variable DEBUG = 1")
	} else {
		pterm.DisableDebugMessages()
	}
}
