package cmd

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/userstorytemplates"
	"github.com/simonkienzler/crusado/pkg/workitems"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/work"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

var (
	crusadoCmd = &cobra.Command{
		Use:   "crusado",
		Short: "TODO",
		Long:  `TODO`,
		Args:  args,
		Run:   Crusado,
	}
)

var (
	dryRunFlag bool
)

func Execute() error {
	return crusadoCmd.Execute()
}

func init() {
	// TODO change the default to false at some point, it's the more sensible default for actual usage
	crusadoCmd.PersistentFlags().BoolVar(&dryRunFlag, "dry-run", true, "if set to true, crusado doesn't actually create work items in Azure DevOps")
}

func args(cmd *cobra.Command, args []string) error {
	// TODO validate args
	return nil
}

func Crusado(cmd *cobra.Command, args []string) {
	// TODO implement proper contexts
	ctx := context.Background()

	// create user story template struct
	userStoryTemplateExample, err := ioutil.ReadFile("./example/userStoryTemplate.yaml")
	if err != nil {
		log.Fatalf("Could not read example template file: %s", err)
	}

	userStoryTemplate := config.UserStoryTemplate{}
	yaml.Unmarshal(userStoryTemplateExample, &userStoryTemplate)

	workitemsService, err := createWorkitemsService(ctx, dryRunFlag)
	if err != nil {
		log.Fatalf("Error during service creation: %s", err)
	}

	userStoryTemplatesService := &userstorytemplates.Service{
		WorkitemsService: *workitemsService,
	}

	if err := userStoryTemplatesService.CreateWorkitemsFromUserStoryTemplate(ctx, userStoryTemplate); err != nil {
		log.Fatalf("Error during user story template creation: %s", err)
	}
}

func createWorkitemsService(ctx context.Context, useDryRunMode bool) (*workitems.Service, error) {
	crusadoConfig := config.GetConfig(true)

	// create a connection to the organization
	connection := azuredevops.NewPatConnection(crusadoConfig.OrganizationUrl, crusadoConfig.PersonalAccessToken)

	// create clients required by the workitems service
	workitemClient, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}

	workClient, err := work.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}

	// configure the workitems service
	workitemsService := workitems.Service{
		WorkitemClient: workitemClient,
		WorkClient:     workClient,
		DryRun:         useDryRunMode,
		ProjectConfig: config.ProjectConfig{
			Name:     crusadoConfig.ProjectName,
			AreaPath: crusadoConfig.ProjectName,
		},
	}

	// get current iteration, either from env var or from the API
	workitemsService.ProjectConfig.IterationPath = crusadoConfig.IterationPath
	if !crusadoConfig.UseIterationPathFromEnvVar {
		log.Print("Getting path of current iteration...")

		currentIteration, err := workitemsService.GetCurrentIteration(ctx)
		if err != nil {
			return nil, err
		}

		log.Printf("Got current iteration path: %+v", *currentIteration.Path)
		workitemsService.ProjectConfig.IterationPath = *currentIteration.Path
	} else {
		log.Printf("Using configured iteration path: %+v", workitemsService.ProjectConfig.IterationPath)
	}

	return &workitemsService, nil
}
