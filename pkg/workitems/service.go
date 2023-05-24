package workitems

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/simonkienzler/crusado/pkg/config"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/work"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
)

// addOp is a shortcut variable for the Add operation.
var addOp = webapi.OperationValues.Add

// Service deals with workitems and provides functions
// for creating user stories or tasks.
type Service struct {
	WorkitemClient workitemtracking.Client
	WorkClient     work.Client
	DryRun         bool
	ProjectConfig  config.ProjectConfig
}

func (s *Service) GetCurrentIteration(ctx context.Context) (*work.TeamSettingsIteration, error) {
	iterations, err := s.WorkClient.GetTeamIterations(ctx, work.GetTeamIterationsArgs{
		Project:   &s.ProjectConfig.Name,
		Timeframe: stringPointer("current"),
	})
	if iterations == nil {
		return nil, err
	}

	if len(*iterations) != 1 {
		return nil, fmt.Errorf("Search for current iteration returned %d results", len(*iterations))
	}

	currentIteration := (*iterations)[0]

	return &currentIteration, nil
}

func (s *Service) CreateUserStory(ctx context.Context, title, description string) (*workitemtracking.WorkItem, error) {
	project := s.ProjectConfig.Name
	workItemType := "User Story"
	validateOnly := s.DryRun
	document := s.buildBasicWorkItemJSONPatchDocument(title, description)

	return s.WorkitemClient.CreateWorkItem(ctx, workitemtracking.CreateWorkItemArgs{
		Document:     &document,
		Project:      &project,
		Type:         &workItemType,
		ValidateOnly: &validateOnly,
	})
}

func (s *Service) CreateTaskUnderneathUserStory(ctx context.Context, title, description string, parent *workitemtracking.WorkItem) (*workitemtracking.WorkItem, error) {
	project := s.ProjectConfig.Name
	workItemType := "Task"
	validateOnly := s.DryRun
	document := s.buildBasicWorkItemJSONPatchDocument(title, description)

	if parent == nil {
		return nil, errors.New("cannot create task underneath user story without parent")
	}

	log.Printf("Parent URL: %s", *parent.Url)

	document = append(document, buildJSONPatchOperation(
		addOp, "/relations/-", workitemtracking.WorkItemRelation{
			Url: parent.Url,
			Rel: stringPointer("System.LinkTypes.Hierarchy-Reverse"),
		},
	))

	return s.WorkitemClient.CreateWorkItem(ctx, workitemtracking.CreateWorkItemArgs{
		Document:     &document,
		Project:      &project,
		Type:         &workItemType,
		ValidateOnly: &validateOnly,
	})
}

func (s *Service) buildBasicWorkItemJSONPatchDocument(title, description string) []webapi.JsonPatchOperation {
	return []webapi.JsonPatchOperation{
		buildJSONPatchOperation(addOp, "/fields/System.Title", title),
		buildJSONPatchOperation(addOp, "/fields/System.Description", description),
		buildJSONPatchOperation(addOp, "/fields/System.AreaPath", s.ProjectConfig.AreaPath),
		buildJSONPatchOperation(addOp, "/fields/System.IterationPath", s.ProjectConfig.IterationPath),
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
