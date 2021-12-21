// fancy uses pterm to provide some nice output that's not really critical but makes the experience nicer with summary
package docker_test

import (
	"os"
	"strings"
	"testing"

	iz "github.com/matryer/is"
	"github.com/sheldonhull/magetools/docker"
)

func TestBuild(t *testing.T) {
	if !strings.Contains(strings.ToLower(os.Getenv("GOTESTS")), "slow") {
		t.Skip("GOTESTS should include 'slow' to run this test")
	}

	is := iz.New(t)
	err := docker.Docker{}.Build("Dockerfile", "magetoolsubuntu", "latest")
	is.NoErr(err) // build Dockerfile without error
}
