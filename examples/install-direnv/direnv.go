//go:build examples

// Package direnv provides an example of using Mage tasks to install a tool like Direnv using what normally is a curl | bash style command.
package direnv

import (
	"strings"
	"time"

	"github.com/bitfield/script"
	"github.com/dustin/go-humanize"
	"github.com/pterm/pterm"
)

// relTime returns just a simple relative time humanized, without the "ago" suffix.
func relTime(t time.Time) string {
	return strings.ReplaceAll(humanize.Time(t), " ago", "")
}

// ⚙️ Direnv simplifies project env variable and path accessibility loading for a project.
func Direnv() error {
	start := time.Now()
	defer func(start time.Time) {
		pterm.Success.Printf("✅ (Direnv) [%s]\n", relTime(start))
	}(start)
	pterm.DefaultHeader.Println("⚙️ Direnv")
	p := script.Exec("curl -sfL https://direnv.net/install.sh").Exec("bash")
	output, err := p.String()
	pterm.Info.Println(output)
	pterm.Error.PrintOnError(err)
	if err != nil {
		return err
	}
	return nil
}
