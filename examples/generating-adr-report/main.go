//go:build examples

/*
This is an example of using mage to generate a report on `adrgen` managed files in a repo.

This parses the text output of the table and splits into a text slice that gets converted into a markdown formatted table.



*/
package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/olekukonko/tablewriter"
	"github.com/pterm/pterm"

	mtu "github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

// Adr contains mage tasks related to managing ADRs.
type Adr mg.Namespace

const (
	// adrReportFile contains the output of adrgen list parsed into a markdown table and exported to this file.
	adrReportFile = "docs/adr/report.md"
)

// 📈 Report saves adrgen list results into a markdown file
func (Adr) Report() error {
	mtu.CheckPtermDebug()
	b := bytes.Buffer{}
	c := exec.Command("adrgen", "list")
	c.Stdout = &b
	if err := c.Run(); err != nil {
		return err
	}
	pterm.Debug.Printf(b.String())
	scanner := bufio.NewScanner(strings.NewReader(b.String()))
	tableString := &strings.Builder{}
	tableString.WriteString("# ADR Report\n\n")
	table := tablewriter.NewWriter(tableString)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	table.SetTablePadding(" ")

	table.SetHeader([]string{"Title", "Status", "Date", "ID", "FileName"})
	var startParsing bool = false

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, "|")
		if !startParsing {
			pterm.Debug.Printf("startParsing is: %v\n", startParsing)

			if strings.Contains(strings.ToLower(line), "title") {
				startParsing = true
				pterm.Debug.Printf("startParsing has been changed to: %v\n", startParsing)
				continue
			}
			continue
		}
		r, err := regexp.Compile(`\s{2,100}`)
		if err != nil {
			pterm.Error.Printf("regex compile\n")
			return err
		}
		result := r.ReplaceAllString(line, "|")
		pterm.Debug.Printf("%v\n", result)
		splitString := strings.Split(result, "|")
		table.Append(splitString)
	}
	table.Render() // Send output
	ioutil.WriteFile(adrReportFile, []byte(tableString.String()), os.ModeAppend)
	return nil
}
