package main

import (
	"context"
	"log"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/userstorytemplates"
	"github.com/simonkienzler/crusado/pkg/workitems"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/work"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	yaml "gopkg.in/yaml.v3"
)

const userStoryTemplateExample = `name: example-user-story-template
description: This template demonstrates the capabilities of crusado.
storyTitle: Try out crusado
storyDescription: crusado looks like a great tool. We should test it.
tasks:
  - title: Download crusado
    description: step 1
  - title: Test crusado
    description: step 2
  - title: Document test results
    description: step 3
`

func main() {
	// create user story template struct
	userStoryTemplate := config.UserStoryTemplate{}
	yaml.Unmarshal([]byte(userStoryTemplateExample), &userStoryTemplate)
	log.Printf("User Story Template: %+v", userStoryTemplate)

	ctx := context.Background()

	useDryRunMode := true

	workitemsService, err := createWorkitemsService(ctx, useDryRunMode)
	if err != nil {
		log.Fatalf("Error during service creation: %s", err)
	}

	userStoryTemplatesService := &userstorytemplates.Service{
		WorkitemsService: *workitemsService,
	}

	// create user story and tasks from template

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
