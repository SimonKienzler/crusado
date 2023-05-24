package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/work"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
)

type CrusadoConfig struct {
	OrganizationUrl     string
	PersonalAccessToken string
	ProjectName         string
}

func getConfig() CrusadoConfig {
	organizationUrl := os.Getenv("AZURE_ORG_URL")
	personalAccessToken := os.Getenv("AZURE_PAT")
	projectName := os.Getenv("AZURE_PROJECT_NAME")

	return CrusadoConfig{
		OrganizationUrl:     organizationUrl,
		PersonalAccessToken: personalAccessToken,
		ProjectName:         projectName,
	}
}

func main() {
	config := getConfig()

	// Create a connection to your organization
	connection := azuredevops.NewPatConnection(config.OrganizationUrl, config.PersonalAccessToken)

	ctx := context.Background()

	// get current iteration

	workClient, err := work.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}

	currentIteration, err := getCurrentIteration(ctx, workClient, config.ProjectName)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", *currentIteration.Path)

	// create user story in current iteration

	workitemClient, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}

	workItem, err := createUserStory(ctx, workitemClient, config.ProjectName, *currentIteration.Path)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", *workItem)
}

func getCurrentIteration(ctx context.Context, workClient work.Client, projectName string) (*work.TeamSettingsIteration, error) {
	iterations, err := workClient.GetTeamIterations(ctx, work.GetTeamIterationsArgs{
		Project:   &projectName,
		Timeframe: stringPointer("current"),
	})
	if err != nil || iterations == nil {
		return nil, err
	}

	if len(*iterations) != 1 {
		return nil, fmt.Errorf("Search for current iteration returned %d results", len(*iterations))
	}

	currentIteration := (*iterations)[0]

	return &currentIteration, nil
}

func createUserStory(ctx context.Context, workitemClient workitemtracking.Client, projectName, currentIteration string) (*workitemtracking.WorkItem, error) {
	project := projectName
	workItemType := "User Story"
	validateOnly := true
	document := buildJSONPatchDocument(
		"work item creation test",
		"hello from the crusado CLI",
		project,
		currentIteration,
	)

	return workitemClient.CreateWorkItem(ctx, workitemtracking.CreateWorkItemArgs{
		Document:     &document,
		Project:      &project,
		Type:         &workItemType,
		ValidateOnly: &validateOnly,
	})
}

func buildJSONPatchDocument(title, description, area, iteration string) []webapi.JsonPatchOperation {
	return []webapi.JsonPatchOperation{
		{
			Op:    &webapi.OperationValues.Add,
			Path:  stringPointer("/fields/System.Title"),
			Value: title,
		},
		{
			Op:    &webapi.OperationValues.Add,
			Path:  stringPointer("/fields/System.Description"),
			Value: description,
		},
		{
			Op:    &webapi.OperationValues.Add,
			Path:  stringPointer("/fields/System.AreaPath"),
			Value: area,
		},
		{
			Op:    &webapi.OperationValues.Add,
			Path:  stringPointer("/fields/System.IterationPath"),
			Value: iteration,
		},
	}
}

func stringPointer(s string) *string {
	return &s
}
