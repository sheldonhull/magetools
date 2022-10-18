//go:build examples

package secrets_test

import (
	"os"
	"testing"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/secrets"
)

func TestGo_SecretsDetect(t *testing.T) {
	is := iz.New(t)
	pterm.DisableStyling()
	err := secrets.Secrets{}.Detect()
	is.NoErr(err) // Secret check should not fail

	// cleanup
	defer func() {
		err := os.RemoveAll(".artifacts")
		is.NoErr(err) // Clean should not fail
	}()
}

func TestGo_SecretsProtect(t *testing.T) {
	is := iz.New(t)
	pterm.DisableStyling()
	err := secrets.Secrets{}.Protect()
	is.NoErr(err) // Secret check should not fail

	// cleanup
	defer func() {
		err := os.RemoveAll(".artifacts")
		is.NoErr(err) // Clean should not fail
	}()
}
