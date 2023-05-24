package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/simonkienzler/crusado/pkg/config"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/work"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
)

var addOp = webapi.OperationValues.Add

func main() {
	config := config.GetConfig(true)

	// Create a connection to your organization
	connection := azuredevops.NewPatConnection(config.OrganizationUrl, config.PersonalAccessToken)

	ctx := context.Background()

	// get current iteration
	currentIterationPath := config.IterationPath
	if !config.UseIterationPathFromEnvVar {
		log.Print("Getting path of current iteration...")

		workClient, err := work.NewClient(ctx, connection)
		if err != nil {
			log.Fatal(err)
		}

		currentIteration, err := getCurrentIteration(ctx, workClient, config.ProjectName)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Got current iteration path: %+v", *currentIteration.Path)
		currentIterationPath = *currentIteration.Path
	} else {
		log.Printf("Using configured iteration path: %+v", currentIterationPath)
	}

	// create user story in current iteration

	workitemClient, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}

	workItem, err := createUserStory(ctx, workitemClient, config.ProjectName, currentIterationPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("User Story: %+v", *workItem)

	// create task underneath the user story
	task, err := createTaskUnderneathUserStory(ctx, workitemClient, config.ProjectName, currentIterationPath, workItem.Url)
	if err != nil {
		log.Fatalf("Error during task creation: %s", err)
	}

	log.Printf("Task: %+v", *task)
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
	document := buildBasicWorkItemJSONPatchDocument(
		"this is a user story",
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

func createTaskUnderneathUserStory(ctx context.Context, workitemClient workitemtracking.Client, projectName, currentIteration string, parentUrl *string) (*workitemtracking.WorkItem, error) {
	project := projectName
	workItemType := "Task"
	validateOnly := true
	document := buildBasicWorkItemJSONPatchDocument(
		"this is a task",
		"hello from the crusado CLI",
		project,
		currentIteration,
	)

	if parentUrl == nil {
		return nil, errors.New("cannot create task underneath user story without parent url")
	}

	log.Printf("Parent URL: %s", *parentUrl)

	document = append(document, buildJSONPatchOperation(
		addOp, "/relations/-", workitemtracking.WorkItemRelation{
			Url: parentUrl,
			Rel: stringPointer("System.LinkTypes.Hierarchy-Reverse"),
		},
	))

	return workitemClient.CreateWorkItem(ctx, workitemtracking.CreateWorkItemArgs{
		Document:     &document,
		Project:      &project,
		Type:         &workItemType,
		ValidateOnly: &validateOnly,
	})
}

func buildBasicWorkItemJSONPatchDocument(title, description, areaPath, iterationPath string) []webapi.JsonPatchOperation {
	return []webapi.JsonPatchOperation{
		buildJSONPatchOperation(addOp, "/fields/System.Title", title),
		buildJSONPatchOperation(addOp, "/fields/System.Description", description),
		buildJSONPatchOperation(addOp, "/fields/System.AreaPath", areaPath),
		buildJSONPatchOperation(addOp, "/fields/System.IterationPath", iterationPath),
	}
}

func buildJSONPatchOperation(op webapi.Operation, path string, value interface{}) webapi.JsonPatchOperation {
	return webapi.JsonPatchOperation{
		Op:    &op,
		Path:  &path,
		Value: value,
	}
}

func stringPointer(s string) *string {
	return &s
}
