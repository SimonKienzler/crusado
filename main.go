package main

import (
	"context"
	"log"

	"github.com/simonkienzler/crusado/pkg/config"
	"github.com/simonkienzler/crusado/pkg/workitems"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/work"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
)

func main() {
	crusadoConfig := config.GetConfig(true)

	// create a connection to the organization
	connection := azuredevops.NewPatConnection(crusadoConfig.OrganizationUrl, crusadoConfig.PersonalAccessToken)

	ctx := context.Background()

	// create clients required by the workitems service
	workitemClient, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}

	workClient, err := work.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}

	// configure the workitems service
	workitemsService := workitems.Service{
		WorkitemClient: workitemClient,
		WorkClient:     workClient,
		DryRun:         true,
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
			log.Fatal(err)
		}

		log.Printf("Got current iteration path: %+v", *currentIteration.Path)
		workitemsService.ProjectConfig.IterationPath = *currentIteration.Path
	} else {
		log.Printf("Using configured iteration path: %+v", workitemsService.ProjectConfig.IterationPath)
	}

	// create user story in current iteration

	userStory, err := workitemsService.CreateUserStory(ctx, "A user story", "hello from crusado")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("User Story: %+v", *userStory)

	// create task underneath the user story

	task, err := workitemsService.CreateTaskUnderneathUserStory(ctx, "A task", "", userStory)
	if err != nil {
		log.Fatalf("Error during task creation: %s", err)
	}

	log.Printf("Task: %+v", *task)
}
