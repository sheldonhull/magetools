package tools

import (
	"os"
	"path/filepath"

	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
)

// _tools is a directory containing local tooling for the project.
const toolDirectory = "_tools"

// mkdirPermissions creates sets the permission
const mkdirPermissions = 0700

// // tools is a list of Go tools to install to avoid polluting global modules.
// var tools = []string{ //nolint:gochecknoglobals // ok to be global for tooling setup
// 	"github.com/goreleaser/goreleaser@v0.174.1",
// 	"golang.org/x/tools/cmd/goimports@master",
// 	"github.com/sqs/goreturns@master",
// 	"github.com/golangci/golangci-lint/cmd/golangci-lint@master",
// 	"github.com/dustinkirkland/golang-petname/cmd/petname@master"
// }

// createDirectories creates the local working directories for build artifacts and tooling.
func createDirectories() error {
	for _, dir := range []string{toolDirectory} {
		if err := os.MkdirAll(dir, 0700); err != nil { //nolint:gomnd // file permissions ok to be literal
			pterm.Error.Printf("failed to create dir: [%s] with error: %v\n", dir, err)

			return err
		}
		pterm.Success.Printf("✅ [%s] dir created\n", dir)
	}

	return nil
}

// Tools installs tooling for the project in a local directory to avoid polluting global modules.
func Tools(tools []string) error {
	if err := os.MkdirAll("_tools", 0700); err != nil { //nolint:gomnd // file permissions ok to be literal
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	env := map[string]string{"GOBIN": filepath.Join(wd, "_tools")}
	args := []string{"get"}

	_, update := os.LookupEnv("UPDATE")
	if update {
		args = []string{"get", "-u"}
	}

	pterm.DefaultSection.Println("Installing go tooling for development")
	p, _ := pterm.DefaultProgressbar.WithTotal(len(tools)).WithTitle("Installing stuff").Start()
	for _, t := range tools {
		p.Title = "Installing " + t
		_, err := sh.OutputWith(env, "go", append(args, t)...)
		if err != nil {
			pterm.Warning.Printf("Could not install [%s] per [%v]\n", t, err)
		}
		pterm.Success.Printf("install %s\n", t)

		p.Increment()
	}

	p.Title = "Tools successfully installed"
	_, _ = p.Stop()

	return nil
}

// tool runs a command using a cached binary.
func tool(cmd string, args ...string) error {
	return sh.Run(filepath.Join("_tools", cmd), args...)
}
