//go:examples

// This package provides an example of executing what typically would be a bash script to generate a kubeconfig for a service account in Kuberenetes.
// Instead of pipeline redirects, it captures the output, cleans and does base64 decoding with Go packages.
// It's a good example of how you can run even command line driven tools in Go, and not require bash scripting.
// These are the type of things that saving into a mage project can make live a lot better when it comes time to repeat this clunky style of inputs and outputs and have pretty good confidence it will keep working for a long time to come! #mageIsMagic.
// Since this is a "build" style script, mostly errors are just returned as is, instead of special handling being required. Kinda strange at first glace if new to Mage, but build/automation scripts don't require as much custom error output.

// While I had a lot of the initial pieced together logic, I found this fantastic script from Armory docs and adapted this, without using the rewrite context as it didn't make sense for my usecase. Original article here: https://docs.armory.io/armory-enterprise/armory-admin/manual-service-account/
// 📍 Long-term it might be more powerful to extract this to a dedicated tool using the Kubernetes API directly, but for expediency I thought this was a nice interim solution to maintain a single interface for commands like this.

package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	mtu "github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

const (
	// ArtifactDirectory is where we drop built artifacts, like the temporary kubeconfig contents. This is an ephemeral workspace.
	ArtifactDirectory = ".artifacts"
	// CachedDirectory contains files locally built and used, but not to commit to source control, and not a typical file in the artifact directory that can be regularly purged during rebuilds.
	CachedDirectory = "~/.cache/kubeconfigs"
	// Kubectrl is the name of the kubectl binary, in case you want to provide a fully qualified path for some reason.
	Kubectl = "kubectl"
	// PermissionUserReadWriteExecute is the octal permission for read, write, & execute only for owner.
	PermissionUserReadWriteExecute = 0o0700
)

type K8 mg.Namespace

// GenerateSvcKubeConfigs generates the kubeconfig files for the service accounts from an input of Clusters to loop through.
// Example is provided in a struct format, but this could be an api call, a yaml input file, or whatever else you want.
// This input would then be used to generate a kubeconfig with embedded cert and auth to use for CI systems and other purposes.
// This is then used to provide input files to pulumi for embedded kubeconfig connections that are encrypted for this stack.

// Cluster is the input struct for the GenerateSvcKubeConfigs function.
type Cluster struct {
	Name string
	// IncludeInAuthLoop is a flag to indicate if this cluster should be included in kubeconfig generation or bypassed.
	IncludeInAuthLoop bool
	// Namespace is the Kubernetes namespace to use for the context and to get the service account from.
	Namespace string
	// ServiceAccount is the name of the service account that we'll be generating credentials for, like `svc-deploy-account`.
	ServiceAccountName string
}

// Generate a kubeconfig for all requested service accounts.
func (K8) GenerateSvcKubeConfigs() error {
	mtu.CheckPtermDebug()

	Clusters := []Cluster{
		{
			Name:               "mycluster",
			IncludeInAuthLoop:  true,
			Namespace:          "ahoy",
			ServiceAccountName: "captainslog-stardate",
		},
		{
			Name:               "mycluster02",
			IncludeInAuthLoop:  true,
			Namespace:          "ahoy",
			ServiceAccountName: "captainslog-stardate",
		},
	}

	if _, err := os.Stat(ArtifactDirectory); err != nil {
		if err := os.Mkdir(ArtifactDirectory, PermissionUserReadWriteExecute); err != nil {
			return err
		}
	}

	if _, err := os.Stat(CachedDirectory); err != nil {
		if err := os.Mkdir(CachedDirectory, PermissionUserReadWriteExecute); err != nil {
			return err
		}
	}

	for _, cluster := range Clusters {
		pterm.DefaultHeader.Printfln("Processing %s", cluster.Name)
		if !cluster.IncludeInAuthLoop {
			pterm.Info.Printfln("Cluster: %s not included in Auth, so bypassing", cluster.Name)
			continue
		}

		if err := sh.RunV(Kubectl, "config", "set-context", cluster.Name, fmt.Sprintf("--namespace=%s", cluster.Namespace)); err != nil {
			return err
		}
		pterm.Success.Printfln("✔️ set context: %q namespace: %q", cluster.Name, cluster.Namespace)

		temporaryFullFile := filepath.Join(
			ArtifactDirectory,
			fmt.Sprintf("kubeconfig-svc-account-%s-%s.full.tmp", cluster.Namespace, cluster.Name),
		)
		temporaryPartialFile := filepath.Join(
			ArtifactDirectory,
			fmt.Sprintf("kubeconfig-svc-account-%s-%s.partial.tmp", cluster.Namespace, cluster.Name),
		)
		targetKubeconfigFile := filepath.Join(
			CachedDirectory,
			fmt.Sprintf("kubeconfig-svc-account-%s-%s.yaml", cluster.Namespace, cluster.Name),
		)
		pterm.Info.Printfln("Target File: %q", targetKubeconfigFile)
		newContext := cluster.Name

		secretName, err := sh.Output(
			"kubectl",
			"get",
			"serviceaccount",
			cluster.ServiceAccountName,
			"--context",
			cluster.Name,
			"--namespace",
			cluster.Namespace,
			"-o",
			"jsonpath='{.secrets[0].name}'",
		)
		if err != nil {
			return err
		}
		secretName = strings.ReplaceAll(secretName, "'", "")
		pterm.Info.Printfln("secretName: %q", secretName)
		tokenData, err := sh.Output(
			"kubectl",
			"get",
			"secret",
			secretName,
			"--context",
			cluster.Name,
			"--namespace",
			cluster.Namespace,
			"-o",
			"jsonpath='{.data.token}'",
		)
		if err != nil {
			return err
		}
		// Pterm.Debug.Printfln("tokenData: %s", tokenData).

		decodedToken, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(strings.TrimSpace(tokenData), "'", ""))
		if err != nil {
			pterm.Error.Printfln("unable to decode token data: %v", err)
			return err
		}
		// Pterm.Printfln("base64 reader: decoded: %s", decodedToken)

		// Create dedicated kubeconfig
		// Create a full copy.
		fullOut, err := sh.Output("kubectl", "config", "view", "--raw")
		if err != nil {
			return err
		}
		if err := os.WriteFile(temporaryFullFile, []byte(fullOut), PermissionUserReadWriteExecute); err != nil {
			return err
		}

		// Switch working context to correct context.
		if err := sh.Run("kubectl", "--kubeconfig", temporaryFullFile, "config", "use-context", newContext); err != nil {
			return err
		}
		// Minify.
		miniOut, err := sh.Output(
			"kubectl",
			"--kubeconfig",
			temporaryFullFile,
			"config",
			"view",
			"--flatten",
			"--minify",
		)
		if err != nil {
			return err
		}
		if err := os.WriteFile(temporaryPartialFile, []byte(miniOut), PermissionUserReadWriteExecute); err != nil {
			return err
		}

		// Create token user.
		if err := sh.Run("kubectl", "config", "--kubeconfig", temporaryPartialFile, "set-credentials", fmt.Sprintf("%s-%s-token-user", cluster.Name, cluster.Namespace), "--token", string(decodedToken)); err != nil {
			return err
		}
		// Set context to use token user.
		if err := sh.Run("kubectl", "config", "--kubeconfig", temporaryPartialFile, "set-context", newContext, "--user", fmt.Sprintf("%s-%s-token-user", cluster.Name, cluster.Namespace)); err != nil {
			return err
		}
		// Set context to correct namespace.
		if err := sh.Run("kubectl", "config", "--kubeconfig", temporaryPartialFile, "set-context", newContext, "--namespace", cluster.Namespace); err != nil {
			return err
		}
		// Flatten/minify kubeconfig.
		newKubeConfigContents, err := sh.Output(
			"kubectl",
			"config",
			"--kubeconfig",
			temporaryPartialFile,
			"view",
			"--flatten",
			"--minify",
		)
		if err != nil {
			return err
		}
		pterm.Debug.Printfln("newKubeConfigContents: %s", newKubeConfigContents)
		if err := os.WriteFile(targetKubeconfigFile, []byte(newKubeConfigContents), PermissionUserReadWriteExecute); err != nil {
			return err
		}
		pterm.Success.Printfln("Kubeconfig created: %q", targetKubeconfigFile)

		// Remove tmp files. If you don't want to debug, you could wrap these in a defer to ensure they always get cleaned up in CI context.
		sh.Rm(temporaryFullFile)
		sh.Rm(temporaryPartialFile)

		pterm.Warning.Println(
			"Don't forget to change your context back to whatever you need.\n\nIt's defaulting now to the latest in this kubeconfig.",
		)
	}
	return nil
}
