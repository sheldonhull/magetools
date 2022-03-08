//go:build examples

/*
Example of installing a tool like Azure-CLI as required (for Linux/macOS).

This allows the command to be called, and if the cli isn't found it will install it automatically.
*/
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/bitfield/script"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	mtu "github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

const (
	// AcrRegistryName is the registry name in Azure to connect to for images.
	AcrRegistryName = "foobar"
)

// Cluster is the metadata used by the azure cli to authenticate and build the Kubeconfig.
type Cluster struct {
	// Name is the friendly alias to use for the cluster, and not the actual cluster name.
	Name              string
	Stage             string
	TenantID          string
	SubscriptionName  string
	SubscriptionID    string
	ResourceGroupName string
	ClusterName       string
	Location          string
	// IncludeInAuthLoop is used to determine if the cluster should be included in the auth loop and generate a kubeconfig.
	IncludeInAuthLoop bool
}

var Clusters = []Cluster{ //nolint:gochecknoglobals // Example setup that could/should be loaded instead from a cloud vault or local config file.
	{
		Name:              "cluster1",
		Stage:             "dev",
		SubscriptionName:  "example",
		TenantID:          "example",
		SubscriptionID:    "example",
		ResourceGroupName: "example",
		ClusterName:       "example",
		Location:          "eastus",
		IncludeInAuthLoop: true,
	},
	{
		Name:              "",
		Stage:             "dev",
		TenantID:          "",
		SubscriptionName:  "",
		SubscriptionID:    "",
		ResourceGroupName: "",
		ClusterName:       "",
		Location:          "",
		IncludeInAuthLoop: true,
	},
	{
		Name:              "",
		Stage:             "",
		TenantID:          "",
		SubscriptionName:  "",
		SubscriptionID:    "",
		ResourceGroupName: "",
		ClusterName:       "",
		Location:          "etc",
		IncludeInAuthLoop: false,
	},
}

// Azure namespace contains tasks that use Azure-CLI.
type Azure mg.Namespace

// checkLinux exits fatal if not linux, as this command isn't supported.
func checklinux() {
	if runtime.GOOS != "Linux" {
		_ = mg.Fatalf(1, "this command is only supported on Linux and you are on: %s", runtime.GOOS)
	}
}

// ResolveAzureCLI returns the azure cli binary path, but if it fails to find it, attempts to install it.
func ResolveAzureCLI() (binary string, err error) {
	mtu.CheckPtermDebug()
	pterm.Debug.Println("resolveAzureCLI")
	binary, err = exec.LookPath("az")
	if err != nil { //nolint:nestif // nestif: System install checks, ok with nested for now - Sheldon, should be low usage
		pterm.Debug.Println("exec.LookPath(\"az\") failed")
		if os.IsNotExist(err) {
			checklinux() // For exit (with panic) as this automation only works for Linux right now.
			pterm.Debug.Println("os.IsNotExist(err) was true, so attempting to start install")
			pterm.Warning.Println("Attempting to install missing azure-cli for you")
			p := script.Exec("curl -sL https://aka.ms/InstallAzureCLIDeb").Exec("sudo bash")
			output, err := p.String()
			if err != nil {
				return "", fmt.Errorf(
					"attempt to install azure-cli for you seemed to have failed.\n you must install azure-cli before you can run this command.\n\n| sudo bash\n\nThis command installs azure-cli if not already installed: error: %w",
					err,
				)
			}
			pterm.Info.Println(output)
			binary, err = exec.LookPath("az")
			if err != nil {
				pterm.Debug.Println("os.IsNotExist(err) was true, so attempting to start install")
				return "", fmt.Errorf("error finding azure-cli after install: %w", err)
			}
		}
	}
	return binary, nil
}

// ➕ InstallCLI runs azure cli installer. Only for linux or macOS.
func (Azure) InstallCLI() error {
	var err error
	mtu.CheckPtermDebug()
	pterm.DefaultSection.Println("(Az) InstallCLI()")
	binary, err := ResolveAzureCLI()
	if err != nil {
		return err
	}
	pterm.Success.Printfln("(Az) InstallCLI(): Binary located at: %q", binary)
	return nil
}

// 🔑 AksAuth loops through the clusters and generates a kubeconfig for each one in your `$HOME/.kube/config`` with overwrite set true.
func (Azure) AksAuth() error {
	mtu.CheckPtermDebug()
	pterm.DefaultHeader.Println("Generating Kubeconfig")
	var binary string
	var err error
	binary, err = ResolveAzureCLI()
	if err != nil {
		return err
	}
	for _, cluster := range Clusters {
		if !cluster.IncludeInAuthLoop {
			pterm.Info.Printfln("skipping login for: %s since IncludeInAuthLoop is false", cluster.Name)
			continue
		}
		pterm.DefaultSection.Printfln("(Az) AksAuth(): %s", cluster.Name)
		if err := sh.Run(binary, []string{"login", "--tenant", cluster.TenantID, "--use-device-code"}...); err != nil {
			return err
		}
		pterm.Success.Printfln("✅ az login to tenant: %s", cluster.TenantID)

		if err := sh.Run(binary, []string{"account", "set", "--subscription", cluster.SubscriptionID}...); err != nil {
			return err
		}
		pterm.Success.Printfln("✅ az set subscription to: %s", cluster.SubscriptionName)
		if err := sh.Run(binary, []string{
			"aks",
			"get-credentials",
			"--resource-group",
			cluster.ResourceGroupName,
			"--name",
			cluster.ClusterName,
			"--subscription",
			cluster.SubscriptionID,
		}...); err != nil {
			return err
		}
		pterm.Success.Printfln("✅ az get-credentials for: %s", cluster.Name)
	}

	return nil
}

// 🔑 ACRLogin runs the azure cli command to login to the registry.
//
// This will allow you to see the registry images contained in ACR in VSCode Docker explorer.
//
// Running manually would entail: `az acr login --name exampleregistry`.
func (Azure) AcrLogin() error {
	var err error
	mtu.CheckPtermDebug()
	pterm.DefaultHeader.Println("(Az) AcrLogin()")
	binary, err := ResolveAzureCLI()
	if err != nil {
		return err
	}
	if err := sh.RunV(binary, "acr", "login", "--name", AcrRegistryName); err != nil {
		return err
	}
	pterm.Success.Println("(Az) AcrLogin()")
	return nil
}
