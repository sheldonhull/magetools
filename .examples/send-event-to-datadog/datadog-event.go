//go:build examples

// This shows an example of providing an output to Datadog events.
// This could be used to pipe a message to datadog events and then datadog could alert in Slack, keeping datadog as the source of truth.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/DataDog/go-datadog-api"
	"github.com/magefile/mage/mg"
	"github.com/pterm/pterm"
)

type Datadog mg.Namespace

var tags = []string{"service:myapp", "app:myapp", os.Getenv("ENVIRONMENT")}

const (
	// AzureDevopsOrg is the org name, without the url.
	AzureDevopsOrg       = "myorg"
	AzureDevopsProject   = "myproject"
	SlackPipelineChannel = "mypipelinechannel"
)

// 📧 Event posts a datadog event to the timeline, which can be integrated into Slack in DD.
func (Datadog) Event() error {
	var err error
	var currentTag string
	apiKey, isSetApi := os.LookupEnv("DD_API_KEY")
	appKey, isSetApp := os.LookupEnv("DD_APP_KEY")

	if !isSetApi || !isSetApp {
		pterm.Warning.Println("[non terminating] this requires DD_API_KEY and DD_APP_KEY to be set.")
		return nil
	}

	client := datadog.NewClient(apiKey, appKey)

	currentTag, err = version.Sbot{}.GetVersion()
	if err != nil {
		pterm.Warning.Println("unable to find tag, proceeding with empty tag value")
	}
	buildPipeline := fmt.Sprintf(
		"https://dev.azure.com/%s/%s/_build/results?buildId=%s&view=results", AzureDevopsOrg, AzureDevopsProject, os.Getenv("BUILD_BUILDID"),
	) // TODO: Remove hard coded values and generate based on build variables.
	var alertType string
	switch strings.ToLower(os.Getenv("AGENT_JOBSTATUS")) {
	case "canceled":
		alertType = "info"
	case "failed":
		alertType = "error"
	case "succeeded":
		alertType = "success"
	case "succeededwithissues":
		alertType = "warning"
	default:
		alertType = "info"
	}
	header := "## Azure Pipelines"
	currentMessage := // Space gap is part of spec.
		header + "\n\n" +
			"- Author" + os.Getenv("BUILD_REQUESTEDFOR") + "\n" +
			"- CommitID" + "\t" + fmt.Sprintf("*%-10s*", os.Getenv("BUILD_SOURCEVERSION")) + "\n" +
			"- Tag" + "\t" + fmt.Sprintf("*%-10s*", currentTag) + "\n" +
			"- Stage" + "\t" + fmt.Sprintf("*%-10s*", os.Getenv("SYSTEM_STAGENAME")) + "\n" +
			"- Status" + "\t" + fmt.Sprintf("*%-10s*", os.Getenv("AGENT_JOBSTATUS")) + "\n" +
			"- Notify" + "\t" + fmt.Sprintf("@%-10s", SlackPipelineChannel)
	// "\n " + `%%%%%%` + "\n" // Prefixed space is part of spec.
	myEvent := &datadog.Event{
		Title:       os.Getenv("BUILD_BUILDNUMBER"),
		Text:        currentMessage, // `"%%%%%%"` + " \n " + currentMessage + "\n " + `%%%%%%` + "\n",.
		Priority:    "Low",
		AlertType:   alertType,
		Tags:        tags,
		Url:         buildPipeline,
		SourceType:  "azure-pipelines",
		Aggregation: currentTag, // Group all the events by the Tag.
	}

	dash, err := client.PostEvent(myEvent)
	if err != nil {
		pterm.Warning.Printfln("[non-terminating] unable to post datadog event: %v", err)
	}
	pterm.Success.Printfln("message sent: %+v", dash)
	return nil
}
