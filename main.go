package main

import (
	"context"
	"log"
	"os"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
)

type CrusadoConfig struct {
	OrganizationUrl     string
	PersonalAccessToken string
	ProjectName         string
	CurrentIteration    string
}

func getConfig() CrusadoConfig {
	organizationUrl := os.Getenv("AZURE_ORG_URL")
	personalAccessToken := os.Getenv("AZURE_PAT")
	projectName := os.Getenv("AZURE_PROJECT_NAME")
	currentIteration := os.Getenv("CURRENT_ITERATION")

	return CrusadoConfig{
		OrganizationUrl:     organizationUrl,
		PersonalAccessToken: personalAccessToken,
		ProjectName:         projectName,
		CurrentIteration:    currentIteration,
	}
}

func main() {
	config := getConfig()

	// Create a connection to your organization
	connection := azuredevops.NewPatConnection(config.OrganizationUrl, config.PersonalAccessToken)

	ctx := context.Background()

	workitemClient, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}

	project := config.ProjectName
	workItemType := "User Story"
	validateOnly := true
	document := buildJSONPatchDocument(
		"work item creation test",
		"hello from the crusado CLI",
		project,
		config.CurrentIteration,
	)

	workItem, err := workitemClient.CreateWorkItem(ctx, workitemtracking.CreateWorkItemArgs{
		Document:     &document,
		Project:      &project,
		Type:         &workItemType,
		ValidateOnly: &validateOnly,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", workItem)
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
