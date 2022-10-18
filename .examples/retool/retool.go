//go:build examples

// retool tooling helps bootstrap a project with locally vendored executables and tools.
// This runs tooling from local `_tools` directory, rather than installing globally.
package retool

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
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
// }

// createDirectories creates the local working directories for build artifacts and tooling.
func createDirectories() error {
	magetoolsutils.CheckPtermDebug()
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
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultHeader.Println("InstallTools")
	start := time.Now()
	if err := createDirectories(); err != nil {
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

	// spinnerLiveText, _ := pterm.DefaultSpinner.Start("InstallTools")
	defer func() {
		duration := time.Since(start)
		msg := fmt.Sprintf("tooling installed: %v\n", duration)
		pterm.Success.Println(msg)
		// spinnerLiveText.Success(msg) // Resolve spinner with success message.
	}()
	for idx, tool := range tools {
		msg := fmt.Sprintf("install [%d] %s]", idx, tool)

		// if windows detected, add the exe to the binary path
		var extension string
		if runtime.GOOS == "windows" {
			extension = ".exe"
		}

		toolPath := filepath.Join("_tools", tool+extension)
		if _, err := os.Stat(toolPath); err == nil {
			pterm.Info.Printf("🔄 [%s] already installed, bypassed.\n", toolPath)

			continue
		}
		_, err := sh.OutputWith(env, "go", append(args, tool)...)
		if err != nil {
			pterm.Warning.Printf("Could not install [%s] per [%v]\n", tool, err)
		}
		pterm.Success.Println(msg)
	}

	pterm.Info.Println("Tools successfully installed (make sure dir is part of .gitignore)")

	return nil
}

// tool runs a command using a cached binary.
func RunTool(cmd string, args ...string) error {
	magetoolsutils.CheckPtermDebug()
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
