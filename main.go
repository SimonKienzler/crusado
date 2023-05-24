package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
)

func main() {
	organizationUrl := os.Getenv("AZURE_ORG_URL")
	personalAccessToken := os.Getenv("AZURE_PAT")

	// Create a connection to your organization
	connection := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)

	ctx := context.Background()

	// Create a client to interact with the Core area
	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}

	// Get first page of the list of team projects for your organization
	responseValue, err := coreClient.GetProjects(ctx, core.GetProjectsArgs{})
	if err != nil {
		log.Fatal(err)
	}

	index := 0
	for responseValue != nil {
		// Log the page of team project names
		for _, teamProjectReference := range (*responseValue).Value {
			log.Printf("Name[%v] = %v", index, *teamProjectReference.Name)
			index++
		}

		// if continuationToken has a value, then there is at least one more page of projects to get
		if responseValue.ContinuationToken != "" {

			continuationToken, err := strconv.Atoi(responseValue.ContinuationToken)
			if err != nil {
				log.Fatal(err)
			}

			// Get next page of team projects
			projectArgs := core.GetProjectsArgs{
				ContinuationToken: &continuationToken,
			}
			responseValue, err = coreClient.GetProjects(ctx, projectArgs)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			responseValue = nil
		}
	}
}
