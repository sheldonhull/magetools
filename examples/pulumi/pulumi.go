//go:build examples

// Package pulumi contains task to help with running Pulumi tools.
package pulumi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"

	"github.com/sheldonhull/magetools/gotools"
	mtu "github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

const (
	// PulumiOrg is the Pulumi organization the stacks are under.
	PulumiOrg = "myorg"

	// PulumiProjectDir is the directory containing the projects, that will be used to set the working directory for pulumi.
	PulumiProjectDir = "pulumi"
)

// Pulumi namespace contains task to help with running Pulumi tools.
type Pulumi mg.Namespace

// Get returns the fully qualified Pulumi stack name, including the org, project, and stage.
// This looks like `myorg/project/stage`.
func GetPulumiStackName(project, stage string) string {
	mtu.CheckPtermDebug()
	return strings.Join([]string{PulumiOrg, project, stage}, "/")
}

// Tidy runs go mod tidy in the nested module directory.
// This helps seperate the project dependency load download from the actual planning.
func (Pulumi) Tidy() error {
	mtu.CheckPtermDebug()
	if err := os.Chdir(PulumiProjectDir); err != nil {
		return err
	}
	if err := (gotools.Go{}.Tidy()); err != nil {
		return err
	}
	return nil
}

// 🚀 Up 👉 Parameters(project, stack string): Eg: `mage pulumi:up myproject dev`.
// Runs pulumi up/apply to target.
//
// Example: `mage pulumi:up myproject dev`.
func (Pulumi) Up(project, stage string) error {
	mtu.CheckPtermDebug()
	return sh.RunV(
		"pulumi",
		"--cwd",
		filepath.Join(PulumiProjectDir, project),
		"--stack="+GetPulumiStackName(project, stage),
		"up",
		"--yes",
		"--emoji",
	)
}

// 🚀 Destroy 👉 Parameters(project, stack string): Eg: `mage pulumi:destroy myproject dev`.
// Runs pulumi destroy/apply to target.
//
// Example: `mage pulumi:destroy myproject dev`.
func (Pulumi) Destroy(project, stage string) error {
	mtu.CheckPtermDebug()
	return sh.RunV(
		"pulumi",
		"--cwd",
		filepath.Join(PulumiProjectDir, project),
		"--stack="+GetPulumiStackName(project, stage),
		"destroy",
		"--yes",
		"--emoji",
	)
}

// 🚀 Refresh 👉 Parameters(project, stack string): Eg: `mage pulumi:refresh myproject dev`.
// Runs pulumi refresh/apply to target.
//
// Example: `mage pulumi:refresh myproject dev`.
func (Pulumi) Refresh(project, stage string) error {
	mtu.CheckPtermDebug()
	return sh.RunV(
		"pulumi",
		"--cwd",
		filepath.Join(PulumiProjectDir, project),
		"--stack="+GetPulumiStackName(project, stage),
		"refresh",
		"--yes",
		"--emoji",
	)
}

// 🚀 Diff 👉 Parameters(project, stack string): Eg: `mage pulumi:preview myproject dev`.
// Runs pulumi preview/apply to target.
//
// Example: `mage pulumi:preview myproject dev`.
func (Pulumi) Diff(project, stage string) error {
	mtu.CheckPtermDebug()
	return sh.RunV(
		"pulumi",
		"--cwd",
		filepath.Join(PulumiProjectDir, project),
		"--stack="+GetPulumiStackName(project, stage),
		"preview",
		"--diff",
		"--yes",
		"--emoji",
	)
}

// 🚀 Watch 👉 Parameters(project, stack string):  Eg: `mage pulumi:watch myproject dev`. Runs pulumi watch of state with automatic deploys as the project file changes.
//
// Example: `mage pulumi:watch myproject dev`.
func (Pulumi) Watch(project, stage string) error {
	mtu.CheckPtermDebug()
	return sh.RunV(
		"pulumi",
		"--cwd",
		filepath.Join(PulumiProjectDir, project),
		"--stack="+GetPulumiStackName(project, stage),
		"watch",
		"--emoji",
	)
}

// GetVersion is an example function of returning back a semver through whatever means you want.
func GetVersion() (string, error) {
	mtu.CheckPtermDebug()
	return "0.0.1", nil
}

// 🚀 Bump 👉 Parameters(project, stack string):  Eg: `mage pulumi:bump myproject dev`. Updates the container semver tag to use in deployment.
//
// Example: `mage pulumi:bump myproject dev`.
func (Pulumi) Bump(project, stage string) error {
	mtu.CheckPtermDebug()
	ver, err := GetVersion()
	if err != nil {
		pterm.Error.WithShowLineNumber(true).WithLineNumberOffset(1).Println("unable to get semver from cli tool")
		return err
	}
	return sh.RunV(
		"pulumi",
		"--cwd",
		filepath.Join(PulumiProjectDir, project),
		"config",
		"set",
		"--stack="+GetPulumiStackName(project, stage),
		"--path",
		"myproject:data.semvertag",
		ver,
		"--emoji",
	)
}

// ⚙️ SetContext: 👉 Parameters(project, stage string). Eg: `mage pulumi:setcontext myproject dev`
//
// Context is pulled by the command `kubectl config current-context`, so you should set this by running generating the current context via kubectl.
//
// Example: `mage pulumi:setcontext myproject dev`.
func (Pulumi) SetContext(project, stage string) error {
	mtu.CheckPtermDebug()
	kubeconfig, isSet := os.LookupEnv("KUBECONFIG")
	if !isSet {
		kubeconfig = ""
	}

	out, err := sh.Output(
		"kubectl",
		"config",
		"current-context",
		fmt.Sprintf("--kubeconfig=%s", kubeconfig),
	)
	if err != nil {
		return err
	}
	currentContext := strings.TrimSpace(out)
	pterm.Info.Printfln("kubectl config current-context: %q", currentContext)
	currentProjectDir := filepath.Join(PulumiProjectDir, project)
	// This was an example from Pulumi docs I adapted to show where possible to change context based on job.
	if currentContext == "minikube" {
		if err := sh.RunV("pulumi", "--cwd", currentProjectDir, "--stack="+GetPulumiStackName(project, stage), "config", "set", "kubernetes-go-exposed-deployment:isMinikube", "true", "--emoji"); err != nil {
			return err
		}
	}
	if err := sh.RunV("pulumi", "--cwd", currentProjectDir, "stack", "select", GetPulumiStackName(project, stage), "--emoji"); err != nil {
		return err
	}
	if err := sh.RunV("pulumi", "--cwd", currentProjectDir, "--stack="+GetPulumiStackName(project, stage), "config", "set", "kubernetes:kubeconfig", kubeconfig, "--emoji"); err != nil {
		return err
	}
	return nil
}
