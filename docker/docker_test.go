// +build testfiles

// fancy uses pterm to provide some nice output that's not really critical but makes the experience nicer with summary
package docker_test

import (
	"testing"

	iz "github.com/matryer/is"

	// mage:import
	"github.com/sheldonhull/magetools/docker"
)

func TestBuild(t *testing.T) {
	is := iz.New(t)
	err := docker.Docker{}.Build("Dockerfile", "magetoolsubuntu", "latest")
	is.NoErr(err) // build Dockerfile without error
}
