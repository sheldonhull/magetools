// tooling helps bootstrap a project with locally vendored executables and tools.
package tooling

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
)

// _tools is a directory containing local tooling for the project.
const toolDirectory = "_tools"

// mkdirPermissions creates sets the permission.
const mkdirPermissions = 0o700

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
		if err := os.MkdirAll(dir, mkdirPermissions); err != nil {
			pterm.Error.Printf("failed to create dir: [%s] with error: %v\n", dir, err)

			return err
		}
		pterm.Success.Printf("✅ [%s] dir created\n", dir)
	}

	return nil
}

// InstallTools installs tooling for the project in a local directory to avoid polluting global modules.
func InstallTools(tools []string) error {
	if err := os.MkdirAll("_tools", mkdirPermissions); err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	env := map[string]string{"GOBIN": filepath.Join(wd, "_tools")}
	args := []string{"install"}

	// _, update := os.LookupEnv("UPDATE")
	// if update {
	// 	args = []string{""}
	// }

	pterm.DefaultSection.Println("installing tooling in local project")
	p, _ := pterm.DefaultProgressbar.WithTotal(len(tools)).WithTitle("Installing stuff").WithRemoveWhenDone(true).Start()
	defer func() {
		p.Title = "tooling installed"
		_, _ = p.Stop()
		pterm.Success.Printf("tooling installed: %s\n", p.GetElapsedTime().String())
	}()

	for _, t := range tools {
		// if windows detected, add the exe to the binary path
		var extension string
		if runtime.GOOS == "windows" {
			extension = ".exe"
		}

		toolPath := filepath.Join("_tools", t+extension)
		if _, err := os.Stat(toolPath); err == nil {
			pterm.Info.Printf("🔄 [%s] already installed, bypassed.\n", toolPath)

			continue
		}
		p.Title = "Installing " + t
		_, err := sh.OutputWith(env, "go", append(args, t)...)
		if err != nil {
			pterm.Warning.Printf("Could not install [%s] per [%v]\n", t, err)
		}
		pterm.Success.Printf("install %s\n", t)
		p.Increment()
	}

	p.Title = "Tools successfully installed (make sure dir is part of .gitignore)"
	_, _ = p.Stop()

	return nil
}

// tool runs a command using a cached binary.
func RunTool(cmd string, args ...string) error {
	// if windows detected, add the exe to the binary path
	var extension string
	if runtime.GOOS == "windows" {
		extension = ".exe"
	}
	if mg.Verbose() {
		return sh.RunV(filepath.Join("_tools", cmd+extension), args...)
	}

	return sh.Run(filepath.Join("_tools", cmd+extension), args...)
}
