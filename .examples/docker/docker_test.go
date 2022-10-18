//go:build examples

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

func Test_DevcontainerBuild(t *testing.T) {
	var err error
	if !strings.Contains(strings.ToLower(os.Getenv("GOTESTS")), "superslow") {
		t.Skip("GOTESTS should include 'superslow' to run this test")
	}
	is := iz.New(t)
	_ = os.Chdir("..")
	err = docker.Devcontainer{}.Build()
	is.NoErr(err) // build Devcontainer without error
}
