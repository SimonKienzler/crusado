package template

import (
	"context"
	"fmt"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/work"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/crusado"
	"github.com/simonkienzler/crusado/pkg/workitems"
	"github.com/spf13/cobra"
	"github.com/thediveo/klo"
)

var (
	RootCmd = &cobra.Command{
		Use:     "template",
		Aliases: []string{"t"},
		Short:   "Work with your crusado templates",
		Long:    `Use the template subcommands to list, show and apply templates.`,
		Args:    cobra.NoArgs,
		Run:     nil,
	}
)

var (
	outputFlag string
)

func init() {
	RootCmd.AddCommand(ListCmd)
	RootCmd.AddCommand(ShowCmd)
	RootCmd.AddCommand(ApplyCmd)
}

func crusadoService() *crusado.Service {
	cfg := config.GetConfigOrDie()

	return &crusado.Service{
		TemplatesDirectory: cfg.TemplatesDirectory,
	}
}

func workitemsService(ctx context.Context, useDryRunMode bool) (*workitems.Service, error) {
	cfg := config.GetConfigOrDie()

	// create a connection to the organization
	connection := azuredevops.NewPatConnection(cfg.OrganizationURL, cfg.PersonalAccessToken)

	workitemClient, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}

	workClient, err := work.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}

	iterationPath, err := workitems.GetIterationPathFromOffset(ctx, workClient, cfg.ProjectName, iterationOffsetFlag)
	if err != nil {
		return nil, err
	}

	// configure the workitems service
	workitemsService := workitems.Service{
		WorkitemClient: workitemClient,
		DryRun:         useDryRunMode,

		ProjectName:   cfg.ProjectName,
		AreaPath:      cfg.ProjectName,
		IterationPath: iterationPath,
	}

	return &workitemsService, nil
}

func getPrinter(outputFormat string) (klo.ValuePrinter, error) {
	return klo.PrinterFromFlag(outputFormat, &crusado.PrinterSpecs)
}

func prettyPrintTemplate(template *crusado.Template) {
	// TODO add color
	fmt.Printf("Name:             %s\n", template.Name)
	fmt.Printf("Type:             %s\n", template.Type)
	fmt.Printf("Title:            %s\n", template.Title)
	fmt.Printf("Number of Tasks:  %d\n", len(template.Tasks))
	fmt.Print("Task Overview:\n")
	for _, task := range template.Tasks {
		fmt.Printf("  - %s\n", task.Title)
	}
}
