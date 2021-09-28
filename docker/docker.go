// Provide docker build and run commands
package docker

import (
	"fmt"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
)

type Docker mg.Namespace

// Build runs docker build command against the provided dockerfile.
func (Docker) Build(dockerfile, imagename, tag string) error {
	pterm.Info.Println("Building docker image")

	dockerfileDirectory := filepath.Dir(dockerfile)
	// if err != nil {
	// 	pterm.Error.Println("unable to resolve directory of dockerfile: %q", dockerfile)

	// 	return err
	// }
	imagetag := fmt.Sprintf("%s:%s", imagename, tag)
	if err := sh.Run("docker", "build", "--pull", "--rm", "-f", dockerfile, "-t", imagetag, dockerfileDirectory); err != nil {
		return err
	}

	return nil
}
