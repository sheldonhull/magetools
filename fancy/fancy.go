// fancy uses pterm to provide some nice output that's not really critical but makes the experience nicer with summary
package fancy

import (
	"time"

	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/ci"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

// ptermMargin defaults to 10. This provides a buffer on items like spinners to avoid whiplash on the console as it refreshes.
const ptermMargin = 10

// IntroScreen provides a nice pterm based output header for console execution with a summary of the runner and time.
func IntroScreen(disableStyling bool) {
	magetoolsutils.CheckPtermDebug()
	if disableStyling || ci.IsCI() {
		pterm.DisableStyling()
	}
	ptermLogo, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("Mage", pterm.NewStyle(pterm.FgLightCyan)),
		pterm.NewLettersFromStringWithStyle("Magic", pterm.NewStyle(pterm.FgLightMagenta))).
		Srender()

	pterm.DefaultCenter.Print(ptermLogo)

	pterm.DefaultCenter.Print(pterm.DefaultHeader.WithFullWidth().
		WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).
		WithMargin(ptermMargin).Sprint("Task Automation With Go"))

	pterm.Info.Println("Taskflow" +
		"\nInitializing dependencies" +
		"\nHelper Libraries: Sheldon Hull " + pterm.LightMagenta("github.com/sheldonhull/magetools") +
		"\n" +
		"\nRunning on: " + pterm.Green(time.Now().Format("02 Jan 2006 - 15:04:05 MST")))

	pterm.Println()
}
